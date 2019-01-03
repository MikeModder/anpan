package anpan

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// SetPrefix Change the prefix
func (c *CommandHandler) SetPrefix(prefix string) {
	c.Prefix = prefix
}

// GetPrefix Get the current prefix
func (c *CommandHandler) GetPrefix() string {
	return c.Prefix
}

// AddCommand Add a command to the Commands map
func (c *CommandHandler) AddCommand(name, desc string, owneronly bool, perms int, run func(*discordgo.Session, *discordgo.MessageCreate, []string)) {
	c.Commands[name] = &Command{
		Name:        name,
		Description: desc,
		OwnerOnly:   owneronly,
		Permissions: perms,
		Run:         run,
	}
}

// RemoveCommand Remove a command from the Commands map
func (c *CommandHandler) RemoveCommand(name string) {
	delete(c.Commands, name)
}

// OnMessage You don't need to call this! Pass this to AddHandler()
func (c *CommandHandler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the author is a bot, and deny entry if IgnoreBots is true
	if m.Author.Bot && c.IgnoreBots {
		return
	}

	content := m.Content

	// Check for the prefix. If the content doesn't start with the prefix, return
	if !strings.HasPrefix(content, c.Prefix) {
		return
	}

	cmd := strings.Split(content, " ")

	// Check and see if we have a command by that name
	if command, exist := c.Commands[cmd[0]]; exist {
		// We do, so check permissions
		if !checkPermissions(s, m.GuildID, m.Author.ID, command.Permissions) {
			embed := &discordgo.MessageEmbed{
				Title:       "You don't have the correct permissions!",
				Description: "You don't have the correct permissions to run this command!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}

		if !checkPermissions(s, m.GuildID, m.Author.ID, command.Permissions) {
			embed := &discordgo.MessageEmbed{
				Title:       "I don't have the correct permissions!",
				Description: "Uh oh, I don't have the right permissions to execute that command for you! Make sure I have the right permissions (ex. Kick Members) then try running the command again!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}

		command.Run(s, m, cmd[1:])
	} else {
		// We don't :(
		return
	}
}
