package anpan

/* Package anpan:
 * Command handler for discordgo. (https://github.com/bwmarrin/discordgo)
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

// New creates a new command handler.
// Prefixes are your command prefixes. If you want to use your
func New(prefixes []string, owners []string, ignoreBots, checkPerms, useRoutines bool) CommandHandler {
	return CommandHandler{
		prefixes:         prefixes,
		owners:           owners,
		ignoreBots:       ignoreBots,
		checkPermissions: checkPerms,
		useRoutines:      useRoutines,
	}
}
