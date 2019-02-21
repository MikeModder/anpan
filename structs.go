package anpan

/* structs.go:
 * Contains the various structs used in anpan.
 *
 * Anpan (c) 2019 MikeModder/MikeModder007
 */

import "github.com/bwmarrin/discordgo"

// CommandHandler contains all the data needed for the handler to function.
type CommandHandler struct {
	Commands         map[string]*Command
	CheckPermissions bool
	Debug            bool
	ErrorFunction    func(Context, *Command, error)
	IgnoreBots       bool
	Owners           []string
	Prefixes         []string
	StatusHandler    StatusHandler
	SuccessFunction  func(Context, *Command)
	//UseDefaultHelp bool
}

// StatusHandler contains status entries and the change interval.
type StatusHandler struct {
	Entries        []string
	SwitchInterval string
}

// Command is the command object.
type Command struct {
	Description string
	Hidden      bool
	Permissions int
	OwnerOnly   bool
	Name        string
	Run         func(context Context, args []string) error
	Type        CommandType
}

// CommandType defines where commands can be used.
type CommandType int

// Command types; Either DM-Only, Guild-Only or both.
const (
	CommandTypePrivate = iota
	CommandTypeGuild
	CommandTypeEverywhere
)

// Context holds the data required for command execution.
type Context struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
	Message *discordgo.Message
	User    *discordgo.User
	Guild   *discordgo.Guild
	Member  *discordgo.Member
}
