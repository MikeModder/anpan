package anpan

/* handler.go:
 * Contains the main code of the command handler.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AddPrefix adds a single prefix to the prefixes.
func (c *CommandHandler) AddPrefix(prefix string) {
	c.Prefixes = append(c.Prefixes, prefix)
}

// RemovePrefix removes a single prefix from the prefixes.
func (c *CommandHandler) RemovePrefix(prefix string) {
	for i, v := range c.Prefixes {
		if v == prefix {
			c.Prefixes = append(c.Prefixes[:i], c.Prefixes[i+1:]...)
			break
		}
	}
}

// SetPrefixes changes all prefixes.
func (c *CommandHandler) SetPrefixes(prefixes []string) {
	c.Prefixes = prefixes
}

// GetPrefixes gets the current prefixes.
func (c *CommandHandler) GetPrefixes() []string {
	return c.Prefixes
}

// SetPrerunFunc sets the function to run before the command handler's OnMessage.
func (c *CommandHandler) SetPrerunFunc(prf PrerunFunc) {
	c.PrerunFunc = prf
}

// ClearPrerunFunc removes the prerun function.
func (c *CommandHandler) ClearPrerunFunc() {
	c.PrerunFunc = nil
}

// SetOnErrorFunc sets the function to run when a command returns an error
func (c *CommandHandler) SetOnErrorFunc(oef OnErrorFunc) {
	c.OnErrorFunc = oef
}

// ClearOnErrorFunc clears the onerror function
func (c *CommandHandler) ClearOnErrorFunc() {
	c.OnErrorFunc = nil
}

// AddCommand adds a command to the Commands map.
func (c *CommandHandler) AddCommand(name, desc string, aliases []string, owneronly, hidden bool, perms int, cmdtype CommandType, run CommandRunFunc) {
	c.Commands = append(c.Commands, &Command{
		Aliases:     aliases,
		Description: desc,
		Hidden:      hidden,
		Name:        name,
		OwnerOnly:   owneronly,
		Permissions: perms,
		Run:         run,
		Type:        cmdtype,
	})
}

// SetHelpCommand sets the help command.
func (c *CommandHandler) SetHelpCommand(name string, aliases []string, perms int, help HelpRunFunc) {
	c.HelpCommand = &HelpCommand{
		Aliases:     aliases,
		Name:        name,
		Permissions: perms,
		Run:         help,
	}
}

// ClearHelpCommand clears the current help command.
func (c *CommandHandler) ClearHelpCommand() {
	c.HelpCommand = nil
}

// RemoveCommand removes a command from the Commands map.
// Note that this only searches for the name, aliases don'tt count.
func (c *CommandHandler) RemoveCommand(name string) {
	for i, v := range c.Commands {
		if v.Name == name {
			c.Commands[i] = nil
			return
		}
	}
}

// IsOwner checks if the given user ID is one of the owners.
func (c *CommandHandler) IsOwner(id string) bool {
	for _, o := range c.Owners {
		if id == o {
			return true
		}
	}

	return false
}

func (c *CommandHandler) debugLog(out string) {
	if c.DebugFunc != nil {
		c.DebugFunc(out)
	}
}

func (c *CommandHandler) errorFunc(context Context, name string, err error) {
	if c.OnErrorFunc != nil {
		c.OnErrorFunc(context, name, err)
	}
}

// OnMessage - You don't need to call this! Pass this to AddHandler().
func (c *CommandHandler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse all content.
	content := m.Content
	c.debugLog("Received: \"" + content + "\"")

	var (
		command *Command
		has     bool
		found   bool
		prefix  string
	)

	// Check for one of the prefixes. If the content doesn't start with one of the prefixes, return.
	for i := 0; i < len(c.Prefixes); i++ {
		prefix = c.Prefixes[i]
		if strings.HasPrefix(content, prefix) {
			has = true
			break
		}
	}

	// If none of the prefixes were found, return.
	if !has {
		return
	}

	// Since we just checked for a prefix, now we need to trim off the prefix
	cmd := strings.Split(strings.TrimPrefix(content, prefix), " ")
	c.debugLog(cmd[0])

	// Before continuing, we need the actual channel itself.
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		c.debugLog("Failed to get the channel.")
		c.errorFunc(Context{
			Session: s,
			Message: m.Message,
			User:    m.Author,
		}, cmd[0], fmt.Errorf(ErrDataUnavailable))
	}

	// Check if the author is a bot, and deny entry if IgnoreBots is true.
	if m.Author.Bot && c.IgnoreBots {
		c.debugLog("Author is bot and IgnoreBots is true.")
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, "", fmt.Errorf(ErrBotBlocked))
		return
	}

	// Check if the command is somehow the given help command.
	if cmd[0] != c.HelpCommand.Name {
		for _, v := range c.HelpCommand.Aliases {
			if cmd[0] == v {
				c.HelpCommand.Run(Context{
					Session: s,
					Channel: channel,
					Message: m.Message,
					User:    m.Author,
				}, cmd[1:], c.Commands, c.Prefixes)
			}
		}
	} else {
		c.HelpCommand.Run(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, cmd[1:], c.Commands, c.Prefixes)
	}

	// Execution will continue on if the help command wasn't found.
	// Let's search for the command.
	for !found {
		for _, v := range c.Commands {
			if cmd[0] == v.Name {
				command = v
				found = true
			}

			for _, v2 := range v.Aliases {
				if cmd[0] == v2 {
					command = v
					found = true
				}
			}
		}
	}

	// If the command is still nil, return an error to the OnErrorFunc and stop.
	if command == nil {
		c.debugLog("Invalid command.")
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, cmd[0], fmt.Errorf(ErrCommandNotFound))

		return
	}

	if channel.Type == discordgo.ChannelTypeDM && command.Type == CommandTypeGuild {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, fmt.Errorf(ErrDMOnly))
		c.debugLog("Tried to run DM-only command on a guild.")
		return
	} else if channel.Type == discordgo.ChannelTypeGuildText && command.Type == CommandTypePrivate {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, fmt.Errorf(ErrGuildOnly))
		c.debugLog("Tried to run Guild-Only command inside DM.")
		return
	}

	// Check if it's an owner-only command, and if it is make sure the author is an owner.
	if command.OwnerOnly && !c.IsOwner(m.Author.ID) {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, fmt.Errorf(ErrOwnerOnly))
		c.debugLog("Owner only command.")
		return
	}

	c.debugLog(fmt.Sprintf("Executing %s, firing Pre-Run Function.", command.Name))
	if c.PrerunFunc != nil {
		// Run the user's pre-run function.
		if !c.PrerunFunc(s, m, command.Name, cmd[1:]) {
			return
		}
	}

	var (
		guild  *discordgo.Guild
		member *discordgo.Member
	)

	if guild, err = s.Guild(m.GuildID); err != nil {
		c.debugLog(fmt.Sprintf("Couldn't retrieve Guild (%s), continuing...", err.Error()))
	}

	if member, err = s.State.Member(m.GuildID, m.Author.ID); err != nil {
		c.debugLog(fmt.Sprintf("Couldn't retrieve Member (%s), continuing...", err.Error()))
	}

	if c.CheckPermissions && guild != nil && member != nil && command.Permissions != 0 {
		has = false

		for _, roleID := range member.Roles {
			if has {
				break
			}

			role, err := s.State.Role(guild.ID, roleID)
			if err != nil {
				c.debugLog("Fetching role \"" + roleID + "\" failed: \"" + err.Error() + "\"")
				continue
			}

			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&command.Permissions != 0 {
				has = true
			}
		}
	} else {
		has = true
	}

	if !has {
		return
	}

	context := Context{
		Session: s,
		Message: m.Message,
		User:    m.Author,
		Channel: channel,
		Guild:   guild,
		Member:  member,
	}

	if err = command.Run(context, cmd[1:]); err != nil && c.OnErrorFunc != nil {
		// Run the user's OnErrorFunc if an error occurs.
		c.OnErrorFunc(context, cmd[0], err)
	}
}
