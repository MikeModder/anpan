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
type HelpRunFunc func(Context, []string, []*Command, []string) error

// PrerunFunc is the type for the function that can be run before command execution.
// If all goes well, return true. otherwise, false.
type PrerunFunc func(*discordgo.Session, *discordgo.MessageCreate, string, []string) bool

// OnErrorFunc is the type for the function that can be run.
type OnErrorFunc func(Context, string, error)

// CommandType defines where commands can be used.
type CommandType int

// HelpCommand defines a help command.
type HelpCommand struct {
	Aliases     []string
	Name        string
	Permissions int
	Run         HelpRunFunc
}

// CommandHandler contains all the data needed for the handler to function.
// Anything inside here must be controlled with Get/Set/Remove function.
type CommandHandler struct {
	checkPermissions bool
	commands         []*Command
	debugFunc        DebugFunc
	helpCommand      *HelpCommand
	ignoreBots       bool
	onErrorFunc      OnErrorFunc
	owners           []string
	prefixes         []string
	prerunFunc       PrerunFunc
	useRoutines      bool
}

// Command represents a command.
type Command struct {
	// Aliases contains all aliases for the command.
	// In most cases, these'll be less favored than the name, or ignored.
	Aliases []string

	// Hidden defines a "hidden" command - it shouldn't be shown in a help message.
	Hidden bool

	// Description defines what the command does.
	Description string

	// Name defines the name of the command.
	Name string

	// OwnerOnly marks a command as an owner-only command. If this is true, permission checks will be ignored.
	OwnerOnly bool

	// SelfPermissions defines what permissions the current bot must meet to execute the command.
	SelfPermissions int

	// UserPermissions defines what permissions a user must meet to execute the command.
	UserPermissions int

	// Run defines the command's function.
	Run CommandRunFunc

	// Type defines when a command will be executed - inside direct messages, a guild or anywhere.
	Type CommandType
}

// Context holds the data required for command execution.
type Context struct {
	// Channel defines the channel in which the command has been executed.
	Channel *discordgo.Channel

	// Guild defines the guild in which the command has been executed.
	// Note that this may be nil.
	Guild *discordgo.Guild

	// Member defines the member in the guild in which the command has been executed.
	// Note that, if guild is nil, this will be nil too.
	Member *discordgo.Member

	// Message defines the message that has executed this command.
	Message *discordgo.Message

	// Session defines the current session that was passed to the OnMessage handler.
	Session *discordgo.Session

	// User defines the user that has executed the command.
	User *discordgo.User
}

const (
	// CommandTypePrivate defines a command that cannot be executed in a guild.
	CommandTypePrivate CommandType = iota

	// CommandTypeGuild defines a command that cannot be executed in direct messages.
	CommandTypeGuild

	// CommandTypeEverywhere defines a command that can be executed anywhere.
	CommandTypeEverywhere
)
