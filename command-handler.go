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

// IsOwner Checks if the given user ID is one of the owners
func (c *CommandHandler) IsOwner(id string) bool {
	for _, o := range c.Owners {
		if id == o {
			return true
		}
	}

	return false
}

// SetDebug Sets debug on or off
func (c *CommandHandler) SetDebug(enabled bool) {
	c.Debug = enabled
}

func (c *CommandHandler) debugLog(out string) {
	if c.Debug {
		fmt.Println(out)
	}
}

// OnMessage You don't need to call this! Pass this to AddHandler()
func (c *CommandHandler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the author is a bot, and deny entry if IgnoreBots is true
	if m.Author.Bot && c.IgnoreBots {
		c.debugLog("author not bot")
		return
	}

	content := m.Content
	c.debugLog(content)

	// Check for the prefix. If the content doesn't start with the prefix, return
	if !strings.HasPrefix(content, c.Prefix) {
		c.debugLog("no prefix")
		return
	}

	cmd := strings.Split(strings.TrimPrefix(content, c.Prefix), " ")
	c.debugLog(cmd[0])

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
			c.debugLog("user bad perms")
			return
		}

		if !checkPermissions(s, m.GuildID, m.Author.ID, command.Permissions) {
			embed := &discordgo.MessageEmbed{
				Title:       "I don't have the correct permissions!",
				Description: "Uh oh, I don't have the right permissions to execute that command for you! Make sure I have the right permissions (ex. Kick Members) then try running the command again!",
				Color:       0xff0000,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			c.debugLog("bot bad perms")
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
			c.debugLog("owner only")
			return
		}

		c.debugLog(command.Name)
		command.Run(s, m, cmd[1:])
	} else {
		// We don't :(
		c.debugLog("not a command")
		return
	}
}
