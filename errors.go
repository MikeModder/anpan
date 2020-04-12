package anpan

/* errors.go:
 * Error definitions.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

// Different error messages for the OnMessage handler.
// Note: These will be passed to OnErrorFunc, though parts of the context will not be available.
const (
	ErrBotBlocked                  = "Author is Bot, but IgnoreBots is true"
	ErrCommandNotFound             = "Command not found"
	ErrDataUnavailable             = "Data couldn't be fetched"
	ErrDMOnly                      = "DM-Only on Guild"
	ErrGuildOnly                   = "Guild-Only in DMs"
	ErrOwnerOnly                   = "Owner-Only"
	ErrSelfInsufficientPermissions = "Insufficient permissions (Bot)"
	ErrUserInsufficientPermissions = "Insufficient permissions (User)"
)
