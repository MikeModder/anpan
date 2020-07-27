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

import (
	"os"
	"os/signal"
	"syscall"
)

/* Package anpan:
 * A custom command handler for discordgo. (https://github.com/bwmarrin/discordgo)
 *
 * anpan (c) 2019-2020 MikeModder/MikeModder007, Apfel
 */

// New creates a new command handler.
//
// Parameters:
// prefixes 		- The prefixes to use for the command handler.
// owners			- The owners of this application; these are used for Owner-Only commands.
// ignoreBots		- Whether to ignore users marked as bots or not.
// checkPermissions	- Whether to check permissions or not.
// useRoutines		- Whether to execute commands outside the event's routine.
//
// Notes:
// Refer to MessageHandler to properly activate the command handler.
// State caching must be enabled. StateEnabled shall not be false then.
// If you want to use a mention/ping as a prefix, just add it as a prefix in the format of "<@bot id>", replacing "bot id" with the bot's User ID.
func New(prefixes []string, owners []string, ignoreBots, checkPermssions, useRoutines bool) CommandHandler {
	return CommandHandler{
		enabled:          true,
		prefixes:         prefixes,
		owners:           owners,
		ignoreBots:       ignoreBots,
		checkPermissions: checkPermssions,
		useRoutines:      useRoutines,
	}
}

// WaitForInterrupt makes your application wait for an interrupt.
// A SIGINT, SIGTERM or a console interrupt will make this function stop.
// Note that the Exit function in the os package will make this function stop, too.
func WaitForInterrupt() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
