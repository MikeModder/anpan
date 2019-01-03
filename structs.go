package anpan

/* structs.go:
 * Contains the various structs used in anpan
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler Contains all the data needed for the handler to function
/* Prefix: Prefix used for command
 *
 */
type CommandHandler struct {
	Prefix           string
	Owners           []string
	PlayingStatus    StatusHandler
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
	//UseDefaultHelp   bool
}

// StatusHandler Struct for status handler; contains entries and the change interval
type StatusHandler struct {
	Entries        []string
	SwitchInterval int
}

// Command it's literally just a command
type Command struct {
	Name        string
	Description string
	OwnerOnly   bool
	Permissions int
	Run         func(*discordgo.Session, *discordgo.MessageCreate, []string)
}
