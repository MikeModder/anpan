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

/* context.go:
 * Contains some utility functions for anpan.Context.
 *
 * anpan (c) 2019-2021 MikeModder/MikeModder007, Apfel
 */

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// Reply directly replies with a message.
//
// Parameters:
// message	- The message content.
func (c *Context) Reply(message string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, message)
}

// ReplyComplex combines Reply, ReplyEmbed and ReplyFile as a way to send a message with, for example, Text and an Embed together.
//
// Parameters:
// message	- The message content.
// tts		- Whether the client should read the message out or not.
// embed	- The embed for this message. Refer to discordgo.MessageEmbed for more info.
// files	- The files to send across. These (collectively) cannot pass more than 8 Megabytes. Refer to discordgo.File for information.
func (c *Context) ReplyComplex(message string, tts bool, embed *discordgo.MessageEmbed, files []*discordgo.File) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendComplex(c.Channel.ID, &discordgo.MessageSend{
		Content: message,
		Embed:   embed,
		TTS:     tts,
		Files:   files,
	})
}

// ReplyEmbed directly replies with a embed, but not with a message.
//
// Parameters:
// embed	- The embed for this message. Refer to discordgo.MessageEmbed for more info.
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Channel.ID, embed)
}

// ReplyFile directly replies with a file, but not with a message.
//
// Parameters:
// files	- The files to send across. These (collectively) cannot pass more than 8 Megabytes.
func (c *Context) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, filename, file)
}
