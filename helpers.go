package anpan

import (
	"github.com/bwmarrin/discordgo"
)

func checkPermissions(s *discordgo.Session, guildid, memberid string, required int) bool {
	member, err := s.State.Member(guildid, memberid)
	if err != nil {
		return false
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildid, roleID)
		if err != nil {
			return false
		}

		// If they have admin, return true
		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true
		}

		// If Permissions AND required isn't 0, return true
		if role.Permissions&required != 0 {
			return true
		}
	}

	// We didn't catch anything in the above loop
	// so we simply return false
	return false
}
