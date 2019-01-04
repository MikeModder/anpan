package anpan

/* structs.go:
 * Contains the various structs used in anpan
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler contains all the data needed for the handler to function
type CommandHandler struct {
	Prefix           string
	Owners           []string
	StatusHandler    StatusHandler
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
	Debug            bool
	//UseDefaultHelp   bool
}

// StatusHandler contains status entries and the change interval
type StatusHandler struct {
	Entries        []string
	SwitchInterval string
}

// Command is the command object
type Command struct {
	Name        string
	Description string
	OwnerOnly   bool
	Hidden      bool
	Permissions int
	Run         func(context Context, args []string)
}

// Context holds the data required for command execution
type Context struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
	Message *discordgo.Message
	User    *discordgo.User
	Guild   *discordgo.Guild
	Member  *discordgo.Member
}
