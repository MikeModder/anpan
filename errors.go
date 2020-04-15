package anpan

import "errors"

/* errors.go:
 * Error definitions.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

var (
	// ErrBotBlocked is thrown when the message handler encounters a bot, but ignoring bots was set to true.
	ErrBotBlocked = errors.New("Author is Bot, but bots are ignored")

	// ErrCommandNotFound is thrown when a message tries to invoke an unknown command.
	ErrCommandNotFound = errors.New("Command not found")

	// ErrDataUnavailable is thrown when data is unavailable, like channels, users or something else.
	ErrDataUnavailable = errors.New("Data couldn't be fetched")

	// ErrDMOnly is thrown when a DM-only command is executed on a guild.
	ErrDMOnly = errors.New("DM-Only on Guild")

	// ErrGuildOnly is thrown when a guild-only command is executed in direct messages.
	ErrGuildOnly = errors.New("Guild-Only in DMs")

	// ErrOwnerOnly is thrown when an owner-only command is executed.
	ErrOwnerOnly = errors.New("Owner-Only")

	// ErrSelfInsufficientPermissions is thrown when the bot itself does not have enough permissions.
	ErrSelfInsufficientPermissions = errors.New("Insufficient permissions (Bot)")

	// ErrUserInsufficientPermissions is thrown when the user doesn't meet the required permissions.
	ErrUserInsufficientPermissions = errors.New("Insufficient permissions (User)")
)
