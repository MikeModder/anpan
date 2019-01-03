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
}

type StatusHandler struct {
	Entries        []string
	SwitchInterval int
}

type Command struct {
	Name        string
	Description string
	OwnerOnly   bool
	Permissions int
	Run         func(*discordgo.Session, *discordgo.MessageCreate, []string)
}
