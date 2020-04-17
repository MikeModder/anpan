package anpan

/* Package anpan:
 * Command handler for discordgo. (https://github.com/bwmarrin/discordgo)
 * Inspired by:
 *  - Clinet (https://github.com/JoshuaDoes/clinet)
 *  - Harmony (https://github.com/superwhiskers/harmony)
 *
 * Main differences from harmony:
 * 	- Built in help command
 * 	- Handling automatic setting of playing status.
 *	- Permission checking built into the handler.
 * 	- Owner only commands.
 *	- Hidden commands.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

// New creates a new command handler.
func New(prefixes []string, owners []string, ignoreBots, checkPerms bool) CommandHandler {
	return CommandHandler{
		prefixes:         prefixes,
		owners:           owners,
		ignoreBots:       ignoreBots,
		checkPermissions: checkPerms,
	}
}
