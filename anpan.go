package anpan

import "github.com/bwmarrin/discordgo"

/* Package anpan:
 * Command handler for discordgo. (https://github.com/bwmarrin/discordgo)
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

// New creates a new command handler.
// Note #1: session.StateEnabled must be true.
// Note #2: If you want to use a mention/ping as a prefix, just add it as a prefix in the format of "<@bot id>", replacing "bot id" with the bot's User ID.
// Note #3: This automatically adds the MessageCreate handler to your session.
func New(session *discordgo.Session, prefixes []string, owners []string, ignoreBots, checkPerms, useRoutines bool) CommandHandler {
	handler := CommandHandler{
		enabled:          true,
		prefixes:         prefixes,
		owners:           owners,
		ignoreBots:       ignoreBots,
		checkPermissions: checkPerms,
		useRoutines:      useRoutines,
	}

	session.AddHandler(handler.onMessage)

	return handler
}
