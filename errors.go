// Copyright (c) 2019-2020 MikeModder/MikeModder007, Apfel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software.
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package anpan

import "errors"

/* errors.go:
 * Error definitions.
 *
 * anpan (c) 2019-2020 MikeModder/MikeModder007, Apfel
 */

var (
	// ErrBotBlocked is thrown when the message handler encounters a bot, but ignoring bots was set to true.
	ErrBotBlocked = errors.New("anpan: The given author was a bot and the IgnoreBots setting is true")

	// ErrCommandAlreadyRegistered is thrown when a command by the same name was registered previously.
	ErrCommandAlreadyRegistered = errors.New("anpan: Another command was already registered by this name")

	// ErrCommandNotFound is thrown when a message tries to invoke an unknown command, or when an attempt at removing an unregistered command was made.
	ErrCommandNotFound = errors.New("anpan: Command not found")

	// ErrDataUnavailable is thrown when data is unavailable, like channels, users or something else.
	ErrDataUnavailable = errors.New("anpan: Necessary data couldn't be fetched")

	// ErrDMOnly is thrown when a DM-only command is executed on a guild.
	ErrDMOnly = errors.New("anpan: DM-Only command on guild")

	// ErrGuildOnly is thrown when a guild-only command is executed in direct messages.
	ErrGuildOnly = errors.New("anpan: Guild-Only command in DMs")

	// ErrOwnerOnly is thrown when an owner-only command is executed.
	ErrOwnerOnly = errors.New("anpan: Owner-Only command")

	// ErrSelfInsufficientPermissions is thrown when the bot itself does not have enough permissions.
	ErrSelfInsufficientPermissions = errors.New("anpan: Insufficient permissions for the bot")

	// ErrUserInsufficientPermissions is thrown when the user doesn't meet the required permissions.
	ErrUserInsufficientPermissions = errors.New("anpan: Insufficient permissions for the user")
)
