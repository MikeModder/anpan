// Copyright (c) 2019 MikeModder/MikeModder007, Apfel
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

/* handler.go:
 * Contains the main code of the command handler.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// GetCheckPermissions returns whether the message handler checks the permissions or not.
func (c *CommandHandler) GetCheckPermissions(enable bool) bool {
	return c.checkPermissions
}

// SetCheckPermissions sets whether to check for permissions or not.
func (c *CommandHandler) SetCheckPermissions(enable bool) {
	c.checkPermissions = enable
}

// GetEnable returns whether the message handler is actually enabled or not.
func (c *CommandHandler) GetEnable() bool {
	return c.enabled
}

// SetEnable sets whether the message handler shall doing its job.
func (c *CommandHandler) SetEnable(enable bool) {
	c.enabled = enable
}

// GetIgnoreBots returns whether the bot ignores other users marked as bots or not.
func (c *CommandHandler) GetIgnoreBots() bool {
	return c.ignoreBots
}

// SetIgnoreBots sets whether to ignore other bots or not.
func (c *CommandHandler) SetIgnoreBots(enable bool) {
	c.ignoreBots = enable
}

// GetUseRoutines returns whether this command handler runs commands in separate goroutines or not.
func (c *CommandHandler) GetUseRoutines() bool {
	return c.useRoutines
}

// SetUseRoutines sets whether this command handler should run commands in separate goroutines or not.
func (c *CommandHandler) SetUseRoutines(enable bool) {
	c.useRoutines = enable
}

// AddPrefix adds a prefix to the handler.
func (c *CommandHandler) AddPrefix(prefix string) {
	c.prefixes = append(c.prefixes, prefix)
}

// RemovePrefix removes the prefix from the handler, if it exists.
func (c *CommandHandler) RemovePrefix(prefix string) {
	for i, v := range c.prefixes {
		if v == prefix {
			copy(c.prefixes[i:], c.prefixes[i+1:])
			c.prefixes[len(c.prefixes)-1] = ""
			c.prefixes = c.prefixes[:len(c.prefixes)-1]

			break
		}
	}
}

// GetAllPrefixes returns the current prefixes slice.
func (c *CommandHandler) GetAllPrefixes() []string {
	return c.prefixes
}

// SetAllPrefixes overwrites all prefixes within the prefixes slice.
func (c *CommandHandler) SetAllPrefixes(prefixes []string) {
	c.prefixes = prefixes
}

// ClearDebugFunc clears the Debug function
// Refer to DebugFunc for more information.
func (c *CommandHandler) ClearDebugFunc() {
	c.debugFunc = nil
}

// GetDebugFunc returns the current debugging function.
// Refer to DebugFunc for more information.
func (c *CommandHandler) GetDebugFunc() DebugFunc {
	return c.debugFunc
}

// SetDebugFunc sets the given debug function as the debugging function for the command handler.
func (c *CommandHandler) SetDebugFunc(df DebugFunc) {
	c.debugFunc = df
}

// ClearOnErrorFunc removes the current OnError function.
// Refer to OnErrorFunc for more details.
func (c *CommandHandler) ClearOnErrorFunc() {
	c.onErrorFunc = nil
}

// GetOnErrorFunc returns the current OnError function.
// Refer to OnErrorFunc for more details.
func (c *CommandHandler) GetOnErrorFunc() OnErrorFunc {
	return c.onErrorFunc
}

// SetOnErrorFunc sets the supplied OnErrorFunc as the one to use.
// Refer to OnErrorFunc for more details.
func (c *CommandHandler) SetOnErrorFunc(oef OnErrorFunc) {
	c.onErrorFunc = oef
}

// ClearPrerunFunc removes the current PrerunFunc.
// Refer to PrerunFunc for more info.
func (c *CommandHandler) ClearPrerunFunc() {
	c.prerunFunc = func(_ *discordgo.Session, _ *discordgo.MessageCreate, _ string, _ []string) bool {
		return true
	}
}

// GetPrerunFunc returns the current PrerunFunc.
// Refer to PrerunFunc for more info.
func (c *CommandHandler) GetPrerunFunc() PrerunFunc {
	return c.prerunFunc
}

// SetPrerunFunc sets the supplied PrerunFunc as the one to use.
// Refer to PrerunFunc for more info.
func (c *CommandHandler) SetPrerunFunc(prf PrerunFunc) {
	c.prerunFunc = prf
}

// AddCommand adds a command to the Commands map.
//
// Parameters:
// name         - The name of the this command.
// description  - The description for this command.
// aliases      - Additional aliases used for this command.
// owneronly    - Whether only owners can access this command or not.
// hidden       - Whether a help command should hide this command or not.
// selfperms    - The necessary permissions for this command. Set this to "0" if any level is fine.
// userperms    - The necessary permissions for the user to meet to use this command. Set this to "0" if any level is fine.
// cmdtype      - The appropriate command type for this command. Use this to limit commands to direct messages or guilds. Refer to CommandType for help.
// function     - The command itself. Refer to CommandFunc for help.
//
// Errors:
// ErrCommandAlreadyRegistered -> There's already a (help) command with this name.
func (c *CommandHandler) AddCommand(name, desc string, aliases []string, owneronly, hidden bool, selfperms, userperms int, cmdtype CommandType, run CommandFunc) error {
	for _, v := range c.commands {
		if v.Name == name {
			return ErrCommandAlreadyRegistered
		}
	}

	if c.helpCommand != nil && c.helpCommand.Name == name {
		return ErrCommandAlreadyRegistered
	}

	c.commands = append(c.commands, &Command{
		Aliases:         aliases,
		Description:     desc,
		Hidden:          hidden,
		Name:            name,
		OwnerOnly:       owneronly,
		SelfPermissions: selfperms,
		UserPermissions: userperms,
		Run:             run,
		Type:            cmdtype,
	})

	return nil
}

// RemoveCommand removes the supplied command from the command array by using its name.
//
// Errors:
// ErrCommandNotFound -> The given name doesn't belong to any command.
func (c *CommandHandler) RemoveCommand(name string) error {
	for i, v := range c.commands {
		if v.Name == name {
			copy(c.commands[i:], c.commands[i+1:])
			c.commands[len(c.commands)-1] = nil
			c.commands = c.commands[:len(c.commands)-1]

			return nil
		}
	}

	return ErrCommandNotFound
}

// GetHelpCommand returns the current set help command.
// Refer to HelpCommandFunc for help.
func (c *CommandHandler) GetHelpCommand() *HelpCommand {
	return c.helpCommand
}

// SetHelpCommand sets the help command.
//
// Parameters:
// name         - The name of the help command; this should be "help" under normal circumstances.
// aliases      - Additional aliases used for the help command.
// selfperms    - The necessary permissions for this help command. Set this to "0" if any level is fine.
// userperms    - The necessary permissions for the user to meet to use this help command. Set this to "0" if any level is fine.
// function     - The help command itself. Refer to HelpCommandFunc for help.
//
// Notes:
// The command handler always checks for the help command first.
//
// Errors:
// ErrCommandAlreadyRegistered  -> There's already another command that has been registered with the same name.
func (c *CommandHandler) SetHelpCommand(name string, aliases []string, selfperms, userperms int, function HelpCommandFunc) error {
	for _, v := range c.commands {
		if v.Name == name {
			return ErrCommandAlreadyRegistered
		}
	}

	c.helpCommand = &HelpCommand{
		Aliases:         aliases,
		Name:            name,
		SelfPermissions: selfperms,
		UserPermissions: userperms,
		Run:             function,
	}

	return nil
}

// ClearHelpCommand clears the current help command.
func (c *CommandHandler) ClearHelpCommand() {
	c.helpCommand = nil
}

// AddOwner adds a user ID as an owner.
func (c *CommandHandler) AddOwner(id string) {
	c.owners = append(c.owners, id)
}

// RemoveOwner removes a user ID from the owner list.
func (c *CommandHandler) RemoveOwner(id string) {
	for i, v := range c.owners {
		if v == id {
			c.owners = append(c.owners[:i], c.owners[i+1:]...)
			break
		}
	}
}

// SetOwners overwrites the current owner list with the given one.
func (c *CommandHandler) SetOwners(ids []string) {
	c.owners = ids
}

// GetOwners returns the current owner list.
func (c *CommandHandler) GetOwners() []string {
	return c.owners
}

// IsOwner checks whether the given ID is set as an owner.
func (c *CommandHandler) IsOwner(id string) bool {
	for _, o := range c.owners {
		if id == o {
			return true
		}
	}

	return false
}

func (c *CommandHandler) debugLog(format string, a ...interface{}) {
	if c.debugFunc != nil {
		c.debugFunc(fmt.Sprintf(format, a...))
	}
}

func (c *CommandHandler) throwError(context Context, command *Command, args []string, err error) {
	if c.onErrorFunc != nil {
		c.onErrorFunc(context, command, args, err)
	}
}

func permissionCheck(session *discordgo.Session, member *discordgo.Member, guild *discordgo.Guild, channel *discordgo.Channel, necessaryPermissions int) error {
	var permissions int

	if member.User.ID == guild.OwnerID {
		return nil
	}

	permissions |= guild.Roles[0].Permissions // everyone role

	for _, roleID := range member.Roles {
		role, err := session.State.Role(guild.ID, roleID)
		if err != nil {
			return err
		}

		if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			return nil
		}

		permissions |= role.Permissions
	}

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.ID == member.User.ID {
			permissions = permissions &^ overwrite.Deny
			permissions |= overwrite.Allow
		}

		for _, roleID := range member.Roles {
			if overwrite.ID == roleID {
				permissions = permissions &^ overwrite.Deny
				permissions |= overwrite.Allow
			}
		}
	}

	if permissions&necessaryPermissions != necessaryPermissions {
		return errors.New("insufficient perms")
	}

	return nil
}

func (c *CommandHandler) onMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	if !c.enabled || event.Author.ID == s.State.User.ID {
		return
	}

	c.debugLog("Received message (%s) by user \"%s\" (%s): \"%s\"", event.Message.ID, event.Author.String(), event.Author.ID, event.Message.Content)

	var (
		has    bool
		prefix string

		command *Command

		channel    *discordgo.Channel
		guild      *discordgo.Guild
		member     *discordgo.Member
		selfMember *discordgo.Member

		err error
	)

	context := Context{
		Channel: channel,
		Guild:   guild,
		Member:  member,
		Message: event.Message,
		User:    event.Author,
		Session: s,
	}

	for i := 0; i < len(c.prefixes); i++ {
		prefix = c.prefixes[i]

		if strings.HasPrefix(event.Message.Content, prefix) {
			has = true
			break
		}
	}

	if !has {
		c.debugLog("Message %s doesn't contain any of the prefixes", event.Message.ID)
		return
	}

	content := strings.Split(strings.TrimPrefix(event.Message.Content, prefix), " ")
	c.debugLog("Parsed message %s: \"%s\"", event.Message.ID, content)

	for _, v := range c.commands {
		if content[0] == v.Name {
			command = v
			break
		}

		for _, alias := range v.Aliases {
			if content[0] == alias {
				command = v
				break
			}
		}
	}

	if command == nil {
		c.throwError(context, nil, content[1:], ErrCommandNotFound)
		return
	}

	if channel, err = s.Channel(event.ChannelID); err != nil {
		c.debugLog("Channel fetching failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrDataUnavailable)
		return
	}

	if channel.Type == discordgo.ChannelTypeDM {
		if c.useRoutines {
			go func() { err = command.Run(context, content[1:]) }()
		} else {
			err = command.Run(context, content[1:])
		}

		if err != nil {
			c.throwError(context, command, content[1:], err)
		}

		return
	}

	if guild, err = s.Guild(event.GuildID); err != nil {
		c.debugLog("Guild fetching failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrDataUnavailable)
		return
	}

	if member, err = s.GuildMember(event.GuildID, event.Author.ID); err != nil {
		c.debugLog("Member fetching failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrDataUnavailable)
		return
	}

	if selfMember, err = s.GuildMember(event.GuildID, s.State.User.ID); err != nil {
		c.debugLog("Bot Member fetching failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrDataUnavailable)
		return
	}

	if !c.checkPermissions {
		if c.useRoutines {
			go func() { err = command.Run(context, content[1:]) }()
		} else {
			err = command.Run(context, content[1:])
		}

		if err != nil {
			c.throwError(context, command, content[1:], err)
		}

		return
	}

	if !(command.SelfPermissions > 0) || !(command.UserPermissions > 0) {
		return
	}

	if err = permissionCheck(s, member, guild, channel, command.UserPermissions); err != nil {
		c.debugLog("Permission check for member failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrUserInsufficientPermissions)
		return
	}

	if err = permissionCheck(s, selfMember, guild, channel, command.SelfPermissions); err != nil {
		c.debugLog("Permission check for bot failed: \"%s\"", err.Error())
		c.throwError(context, command, content[1:], ErrUserInsufficientPermissions)
		return
	}

	if c.useRoutines {
		go func() { err = command.Run(context, content[1:]) }()
	} else {
		err = command.Run(context, content[1:])
	}

	if err != nil {
		c.throwError(context, command, content[1:], err)
	}
}
