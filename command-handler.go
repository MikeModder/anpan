package anpan

/* command-handler.go:
 * Contains the main code of the command handler
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// SetPrefix changes the prefix
func (c *CommandHandler) SetPrefix(prefix string) {
	c.Prefix = prefix
}

// GetPrefix gets the current prefix
func (c *CommandHandler) GetPrefix() string {
	return c.Prefix
}

// SetPrerunFunc sets the function to run before the command handler's OnMessage
func (c *CommandHandler) SetPrerunFunc(prf func(*discordgo.Session, *discordgo.MessageCreate)) {
	c.PrerunFunc = prf
}

// ClearPrerunFunc removes the prerun function
func (c *CommandHandler) ClearPrerunFunc() {
	c.PrerunFunc = nil
}

// AddCommand adds a command to the Commands map
func (c *CommandHandler) AddCommand(name, desc string, owneronly bool, hidden bool, perms int, run func(Context, []string)) {
	c.Commands[name] = &Command{
		Name:        name,
		Description: desc,
		OwnerOnly:   owneronly,
		Hidden:      hidden,
		Permissions: perms,
		Run:         run,
	}
}

// RemoveCommand removes a command from the Commands map
func (c *CommandHandler) RemoveCommand(name string) {
	if _, has := c.Commands[name]; has {
		delete(c.Commands, name)
	}
	return
}

// IsOwner checks if the given user ID is one of the owners
func (c *CommandHandler) IsOwner(id string) bool {
	for _, o := range c.Owners {
		if id == o {
			return true
		}
	}

	return false
}

// SetDebug sets debug on or off
func (c *CommandHandler) SetDebug(enabled bool) {
	c.Debug = enabled
}

func (c *CommandHandler) debugLog(out string) {
	if c.Debug {
		fmt.Println(out)
	}
}

// AddDefaultHelpCommand adds the default (library provided) help command to the list of commands
// TODO: users have to manually call this to add the help command, maybe find a way to add it automatially if no help command is detected?
func (c *CommandHandler) AddDefaultHelpCommand() {
	c.AddCommand("help", "Get some help about using the bot", false, false, 0, c.defaultHelpCmd)
}

// OnMessage - You don't need to call this! Pass this to AddHandler()
func (c *CommandHandler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the author is a bot, and deny entry if IgnoreBots is true
	if m.Author.Bot && c.IgnoreBots {
		c.debugLog("Author is bot")
		return
	}

	if c.PrerunFunc != nil {
		c.PrerunFunc(s, m)
	}

	content := m.Content
	c.debugLog(content)

	// Check for the prefix. If the content doesn't start with the prefix, return
	if !strings.HasPrefix(content, c.Prefix) {
		c.debugLog("No prefix in message")
		return
	}

	cmd := strings.Split(strings.TrimPrefix(content, c.Prefix), " ")
	c.debugLog(cmd[0])

	// Check and see if we have a command by that name
	if command, exist := c.Commands[cmd[0]]; exist {
		// We do, so check permissions
		if !checkPermissions(s, m.GuildID, m.Author.ID, command.Permissions) {
			embed := &discordgo.MessageEmbed{
				Title:       "Insufficient Permissions!",
				Description: "You don't have the correct permissions to run this command!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			c.debugLog("Insufficient permissions for User")
			return
		}

		if !checkPermissions(s, m.GuildID, s.State.User.ID, command.Permissions) {
			embed := &discordgo.MessageEmbed{
				Title:       "Insufficient Permissions!",
				Description: "Uh oh, I don't have the right permissions to execute that command for you! Make sure I have the right permissions (ex. Kick Members) then try running the command again!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			c.debugLog("Insufficient permissions for Bot")
			return
		}

		// Check if it's an owner-only command, and if it is make sure the author is an owner
		if command.OwnerOnly && !c.IsOwner(m.Author.ID) {
			// It's an owner-only command, and the user wasn't on the owners list
			embed := &discordgo.MessageEmbed{
				Title:       "You can't run that command!",
				Description: "Sorry, on the bot owner(s) can run that command!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			c.debugLog("Owner only command")
			return
		}

		c.debugLog("Executing " + command.Name)

		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			c.debugLog("Couldn't retrieve Channel, continuing...")
		}

		guild, err := s.Guild(m.GuildID)
		if err != nil {
			c.debugLog("Couldn't retrieve Guild, continuing...")
		}

		member, err := s.State.Member(m.GuildID, m.Author.ID)
		if err != nil {
			c.debugLog("Couldn't retrieve Member, continuing...")
		}

		context := Context{
			Session: s,
			Message: m.Message,
			User:    m.Author,
			Channel: channel,
			Guild:   guild,
			Member:  member,
		}

		command.Run(context, cmd[1:])
	} else {
		// We don't :(
		c.debugLog("Unknown command / not even one")
		return
	}
}

func (c *CommandHandler) defaultHelpCmd(context Context, args []string) {
	if len(args) >= 1 {
		if commannd, has := c.Commands[args[0]]; has {
			// check if the command is hidden
			if commannd.Hidden {
				return
			}

			embed := &discordgo.MessageEmbed{
				Title:       "Help!",
				Description: fmt.Sprintf("Help for command `%s`\n Description: `%s`\nOwner only: `%v`", commannd.Name, commannd.Description, commannd.OwnerOnly),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("The bot's prefix is %s", c.Prefix),
				},
			}

			context.ReplyEmbed(embed)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: fmt.Sprintf("`%s` is not a valid command!", args[0]),
			Color:       0xff0000,
		}

		context.ReplyEmbed(embed)
		return

	}

	count := len(c.Commands)
	var list string

	for name := range c.Commands {
		cmd := c.Commands[name]
		if !cmd.Hidden {
			list += fmt.Sprintf("`%s - %s`\n", cmd.Name, cmd.Description)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Commands:",
		Description: list,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("The bot's prefix is %s | There are %d commands.", c.Prefix, count),
		},
	}

	context.ReplyEmbed(embed)
}
