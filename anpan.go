package anpan

/* Package anpan:
 * Command handler for DiscordGo.
 * Inspired by:
 *  - Clinet (https://github.com/JoshuaDoes/clinet)
 *  - Harmony (https://github.com/superwhiskers/harmony)
 *
 * Main differences from harmony:
 * 	- Built in help command - (maybe? I need to look more into this).
 * 	- Handling automatic setting of playing status.
 *	- Permission checking built into the handler.
 * 	- Owner only commands.
 *  - Hidden commands.
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

// NewCommandHandler Creates a new command handler.
func NewCommandHandler(prefixes []string, owners []string, ignoreBots, checkPerms bool) CommandHandler {
	return CommandHandler{
		Prefixes:         prefixes,
		Owners:           owners,
		StatusHandler:    NewDefaultStatusHandler(),
		Commands:         make(map[string]*Command),
		IgnoreBots:       ignoreBots,
		CheckPermissions: checkPerms,
	}
}

// NewStatusHandler Creates a new status handler.
func NewStatusHandler(entries []string, interval string) StatusHandler {
	return StatusHandler{
		Entries:        entries,
		SwitchInterval: interval,
	}
}

// NewDefaultStatusHandler Creates a new status handler with some "default" settings.
func NewDefaultStatusHandler() StatusHandler {
	return NewStatusHandler([]string{
		"Powered by Golang!",
		"Powered by anpan!",
	}, "60s")
}
