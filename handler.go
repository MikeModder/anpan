package anpan

/* handler.go:
 * Contains the main code of the command handler.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AddPrefix adds a single prefix to the prefixes slice.
func (c *CommandHandler) AddPrefix(prefix string) {
	c.prefixes = append(c.prefixes, prefix)
}

// RemovePrefix removes a single prefix from the prefixes slice.
func (c *CommandHandler) RemovePrefix(prefix string) {
	for i, v := range c.prefixes {
		if v == prefix {
			c.prefixes = append(c.prefixes[:i], c.prefixes[i+1:]...)
			break
		}
	}
}

// SetPrefixes overwrites all prefixes within the prefixes slice.
func (c *CommandHandler) SetPrefixes(prefixes []string) {
	c.prefixes = prefixes
}

// GetPrefixes returns the current prefixes slice.
func (c *CommandHandler) GetPrefixes() []string {
	return c.prefixes
}

// SetPrerunFunc sets the function to run before the command handler's OnMessage.
func (c *CommandHandler) SetPrerunFunc(prf PrerunFunc) {
	c.prerunFunc = prf
}

// GetPrerunFunc returns the current Prerun function.
func (c *CommandHandler) GetPrerunFunc() PrerunFunc {
	return c.prerunFunc
}

// ClearPrerunFunc removes the current Prerun function.
func (c *CommandHandler) ClearPrerunFunc() {
	c.prerunFunc = func(session *discordgo.Session, event *discordgo.MessageCreate, name string, args []string) bool {
		return true
	}
}

// SetOnErrorFunc sets the function to run when a command returns an error
func (c *CommandHandler) SetOnErrorFunc(oef OnErrorFunc) {
	c.onErrorFunc = oef
}

// GetOnErrorFunc returns the current OnError function.
func (c *CommandHandler) GetOnErrorFunc() OnErrorFunc {
	return c.onErrorFunc
}

// ClearOnErrorFunc clears the onerror function
func (c *CommandHandler) ClearOnErrorFunc() {
	c.onErrorFunc = nil
}

// SetDebugFunc sets the given debug function as the debugging function for the command handler.
func (c *CommandHandler) SetDebugFunc(df DebugFunc) {
	c.debugFunc = df
}

// GetDebugFunc returns the currently set debug function.
func (c *CommandHandler) GetDebugFunc() DebugFunc {
	return c.debugFunc
}

// ClearDebugFunc clears the Debug function
func (c *CommandHandler) ClearDebugFunc() {
	c.debugFunc = nil
}

// AddCommand adds a command to the Commands map.
func (c *CommandHandler) AddCommand(name, desc string, aliases []string, owneronly, hidden bool, selfperms, userperms int, cmdtype CommandType, run CommandRunFunc) {
	c.commands = append(c.commands, &Command{
		Aliases:         aliases,
		Description:     desc,
		Hidden:          hidden,
		Name:            name,
		OwnerOnly:       owneronly,
		SelfPermissions: selfperms,
		UserPermissions: userperms,
		Run:             run,
		Type:            cmdtype,
	})
}

// SetHelpCommand sets the help command.
func (c *CommandHandler) SetHelpCommand(name string, aliases []string, perms int, help HelpRunFunc) {
	c.helpCommand = &HelpCommand{
		Aliases:     aliases,
		Name:        name,
		Permissions: perms,
		Run:         help,
	}
}

// ClearHelpCommand clears the current help command.
func (c *CommandHandler) ClearHelpCommand() {
	c.helpCommand = nil
}

// RemoveCommand removes a command from the Commands map.
// Note that this only searches for the name, aliases don'tt count.
func (c *CommandHandler) RemoveCommand(name string) {
	for i, v := range c.commands {
		if v.Name == name {
			c.commands[i] = nil
			return
		}
	}
}

// AddOwner adds a user ID as an owner.
func (c *CommandHandler) AddOwner(id string) {
	c.owners = append(c.owners, id)
}

// RemoveOwner removes a user ID from the owner list.
func (c *CommandHandler) RemoveOwner(id string) {
	for i, v := range c.owners {
		if v == id {
			c.owners = append(c.owners[:i], c.owners[i+1:]...)
			break
		}
	}
}

// SetOwners overwrites the current owner list with the given one.
func (c *CommandHandler) SetOwners(ids []string) {
	c.owners = ids
}

// GetOwners returns the current owner list.
func (c *CommandHandler) GetOwners() []string {
	return c.owners
}

// IsOwner checks whether the given ID is set as an owner.
func (c *CommandHandler) IsOwner(id string) bool {
	for _, o := range c.owners {
		if id == o {
			return true
		}
	}

	return false
}

// IgnoreBots sets whether to ignore bots or not.
func (c *CommandHandler) IgnoreBots(enable bool) {
	c.ignoreBots = enable
}

// CheckPermissions sets whether to check for permissions or not.
func (c *CommandHandler) CheckPermissions(enable bool) {
	c.checkPermissions = enable
}

func (c *CommandHandler) debugLog(out string) {
	if c.debugFunc != nil {
		c.debugFunc(out)
	}
}

func (c *CommandHandler) errorFunc(context Context, name string, err error) {
	if c.onErrorFunc != nil {
		c.onErrorFunc(context, name, err)
	}
}

func permissionCheck(session *discordgo.Session, member *discordgo.Member, guild *discordgo.Guild, channel *discordgo.Channel, necessaryPermissions int) error {
	var permissions int

	if member.User.ID == guild.OwnerID {
		return nil
	}

	for _, roleID := range member.Roles {

		role, err := session.State.Role(guild.ID, roleID)
		if err != nil {
			return err
		}

		if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			return nil
		}

		permissions |= role.Permissions
	}

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.ID == member.User.ID {
			permissions = permissions &^ overwrite.Deny
			permissions |= overwrite.Allow
		}

		for _, roleID := range member.Roles {
			if overwrite.ID == roleID {
				permissions = permissions &^ overwrite.Deny
				permissions |= overwrite.Allow
			}
		}
	}

	if permissions&necessaryPermissions == 0 {
		return errors.New("insufficient perms")
	}

	return nil
}

// OnMessage - You don't need to call this! Pass this to AddHandler().
func (c *CommandHandler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

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
	for i := 0; i < len(c.prefixes); i++ {
		prefix = c.prefixes[i]
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
		}, cmd[0], ErrDataUnavailable)
	}

	// Check if the author is a bot, and deny entry if IgnoreBots is true.
	if m.Author.Bot && c.ignoreBots {
		c.debugLog("Author is bot and IgnoreBots is true.")
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, "", ErrBotBlocked)
		return
	}

	// Check if the command is somehow the given help command.
	if cmd[0] != c.helpCommand.Name {
		for _, v := range c.helpCommand.Aliases {
			if cmd[0] == v {
				c.helpCommand.Run(Context{
					Session: s,
					Channel: channel,
					Message: m.Message,
					User:    m.Author,
				}, cmd[1:], c.commands, c.prefixes)
			}
		}
	} else {
		c.helpCommand.Run(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, cmd[1:], c.commands, c.prefixes)
	}

	// Execution will continue on if the help command wasn't found.
	// Let's search for the command.
	for !found {
		for _, v := range c.commands {
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
		}, cmd[0], ErrCommandNotFound)

		return
	}

	if channel.Type == discordgo.ChannelTypeDM && command.Type == CommandTypeGuild {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, ErrDMOnly)
		c.debugLog("Tried to run a DM-only command on a guild.")
		return
	} else if channel.Type == discordgo.ChannelTypeGuildText && command.Type == CommandTypePrivate {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, ErrGuildOnly)
		c.debugLog("Tried to run a guild-only command inside DM.")
		return
	}

	// Check if it's an owner-only command, and if it is make sure the author is an owner.
	if command.OwnerOnly && !c.IsOwner(m.Author.ID) {
		c.errorFunc(Context{
			Session: s,
			Channel: channel,
			Message: m.Message,
			User:    m.Author,
		}, command.Name, ErrOwnerOnly)
		c.debugLog("Owner only command.")
		return
	}

	c.debugLog(fmt.Sprintf("Executing %s, firing Pre-Run Function.", command.Name))
	if c.prerunFunc != nil {
		// Run the user's pre-run function.
		if !c.prerunFunc(s, m, command.Name, cmd[1:]) {
			return
		}
	}

	var (
		guild      *discordgo.Guild
		member     *discordgo.Member
		selfMember *discordgo.Member
	)

	if guild, err = s.Guild(m.GuildID); err != nil {
		c.debugLog(fmt.Sprintf("Couldn't retrieve guild (%s), continuing...", err.Error()))
	}

	if member, err = s.GuildMember(m.GuildID, m.Author.ID); err != nil {
		c.debugLog(fmt.Sprintf("Couldn't retrieve member (%s), continuing...", err.Error()))
	}

	if selfMember, err = s.GuildMember(m.GuildID, s.State.User.ID); err != nil {
		c.debugLog(fmt.Sprintf("Couldn't retrieve bot member (%s), continuing...", err.Error()))
	}

	var selfHas bool

	if c.checkPermissions && guild != nil && member != nil && selfMember != nil && (command.SelfPermissions != 0 || command.UserPermissions != 0) {
		has = false

		if err := permissionCheck(s, member, guild, channel, command.UserPermissions); err != nil {
			c.debugLog(fmt.Sprintf("User permission check encountered an error: %s", err.Error()))
		} else {
			has = true
		}

		if err := permissionCheck(s, selfMember, guild, channel, command.UserPermissions); err != nil {
			c.debugLog(fmt.Sprintf("Self permission check encountered an error: %s", err.Error()))
		} else {
			selfHas = true
		}
	} else {
		has = true
		selfHas = true
	}

	context := Context{
		Session: s,
		Message: m.Message,
		User:    m.Author,
		Channel: channel,
		Guild:   guild,
		Member:  member,
	}

	if !has {
		c.errorFunc(context, command.Name, ErrUserInsufficientPermissions)
		c.debugLog("User doesn't have sufficient permissions.")
	}

	if !selfHas {
		c.errorFunc(context, command.Name, ErrSelfInsufficientPermissions)
		c.debugLog("Bot doesn't have sufficient permissions.")
	}

	if err = command.Run(context, cmd[1:]); err != nil && c.onErrorFunc != nil {
		// Run the set OnErrorFunc if an error occurs.
		c.onErrorFunc(context, cmd[0], err)
	}
}
