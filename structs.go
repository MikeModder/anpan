package anpan

/* structs.go:
 * Contains the various structs used in anpan.
 *
 * Anpan (c) 2019 MikeModder/MikeModder007
 */

import "github.com/bwmarrin/discordgo"

// CommandHandler contains all the data needed for the handler to function.
type CommandHandler struct {
	Prefixes         []string
	Owners           []string
	StatusHandler    StatusHandler
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
	Debug            bool
	PrerunFunc       PrerunFunc
	OnErrorFunc      OnErrorFunc
	//UseDefaultHelp   bool
}

// StatusHandler contains status entries and the change interval.
type StatusHandler struct {
	Entries        []string
	SwitchInterval string
}

// Command is the command object.
type Command struct {
	Name        string
	Description string
	OwnerOnly   bool
	Hidden      bool
	Permissions int
	Type        CommandType
	Run         CommandRunFunc
}

// Context holds the data required for command execution.
type Context struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
	Message *discordgo.Message
	User    *discordgo.User
	Guild   *discordgo.Guild
	Member  *discordgo.Member
}

/* These are types, not structs, but this is the best place to put them */

// PrerunFunc is the type for the function that can be run before command execution.
type PrerunFunc func(*discordgo.Session, *discordgo.MessageCreate, string, []string) bool

// OnErrorFunc is the type for the function that can be run.
type OnErrorFunc func(Context, string, error)

// CommandRunFunc is a command's run function
type CommandRunFunc func(Context, []string) error

// CommandType defines where commands can be used.
type CommandType int

// Command types; Either DM-Only, Guild-Only or both.
const (
	CommandTypePrivate CommandType = iota
	CommandTypeGuild
	CommandTypeEverywhere
)
