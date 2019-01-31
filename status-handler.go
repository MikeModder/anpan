package anpan

/* status-handler.go:
 * Contains the code for the built-in status handler.
 *
 * Anpan (c) 2019 MikeModder/MikeModder007
 */

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// AddEntry Adds an entry to the status handler's list.
func (s *StatusHandler) AddEntry(entry string) {
	s.Entries = append(s.Entries, entry)
}

// SetEntries Used to set entries in bulk, overrides the current list.
func (s *StatusHandler) SetEntries(entries []string) {
	s.Entries = entries
}

// SetSwitchInterval Sets the time each entry is displayed.
func (s *StatusHandler) SetSwitchInterval(interval string) {
	s.SwitchInterval = interval
}

// OnReady You don't need to call this manually! Pass it to AddHandler to use the built-in status handler.
func (s *StatusHandler) OnReady(session *discordgo.Session, _ *discordgo.Ready) {
	interval, err := time.ParseDuration(s.SwitchInterval)
	if err != nil {
		return
	}

	s.setStatus(session)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.setStatus(session)
	}
}

func (s *StatusHandler) setStatus(session *discordgo.Session) {
	item := rand.Intn(len(s.Entries))
	session.UpdateStatus(0, s.Entries[item])
}
