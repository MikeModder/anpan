package anpan

/* types.go:
 * Declares all types used for and within anpan.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

import "github.com/bwmarrin/discordgo"

// CommandRunFunc is a command's run function.
type CommandRunFunc func(Context, []string) error

// DebugFunc is used for debugging output.
type DebugFunc func(string)

// HelpRunFunc is used for definitions of a help command.
type HelpRunFunc func(context Context, args []string, commands []*Command, prefixes []string) error

// PrerunFunc is the type for the function that can be run before command execution.
// If all goes well, return true. otherwise, false.
type PrerunFunc func(*discordgo.Session, *discordgo.MessageCreate, string, []string) bool

// OnErrorFunc is the type for the function that can be run.
type OnErrorFunc func(Context, string, error)

// CommandType defines where commands can be used.
type CommandType int

// HelpCommand defines a help command.
type HelpCommand struct {
	Aliases []string
	Name    string
	Show    bool
	Run     HelpRunFunc
}

// CommandHandler contains all the data needed for the handler to function.
type CommandHandler struct {
	CheckPermissions bool
	Commands         []*Command
	DebugFunc        DebugFunc
	HelpCommand      HelpCommand
	IgnoreBots       bool
	OnErrorFunc      OnErrorFunc
	Owners           []string
	Prefixes         []string
	PrerunFunc       PrerunFunc
}

// Command represents a command.
type Command struct {
	Aliases     []string
	Hidden      bool
	Description string
	Name        string
	OwnerOnly   bool
	Permissions int
	Run         CommandRunFunc
	Type        CommandType
}

// Context holds the data required for command execution.
type Context struct {
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Member  *discordgo.Member
	Message *discordgo.Message
	Session *discordgo.Session
	User    *discordgo.User
}

// Command types; Either DM-Only, Guild-Only or both.
const (
	CommandTypePrivate CommandType = iota
	CommandTypeGuild
	CommandTypeEverywhere
)
