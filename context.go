package anpan

/* context.go:
 * Contains some utility functions for anpan.Context.
 *
 * anpan (c) 2020 MikeModder/MikeModder007, Apfel
 */

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// Reply directly replies with a message.
func (c *Context) Reply(message string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, message)
}

// ReplyComplex combines Reply, ReplyEmbed and ReplyFile as a way to send a message with, for example, Text and an Embed together.
func (c *Context) ReplyComplex(message string, tts bool, embed *discordgo.MessageEmbed, files []*discordgo.File) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendComplex(c.Channel.ID, &discordgo.MessageSend{
		Content: message,
		Embed:   embed,
		Tts:     tts,
		Files:   files,
	})
}

// ReplyEmbed directly replies with a embed, but not with a message.
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Channel.ID, embed)
}

// ReplyFile directly replies with a file, but not with a message.
func (c *Context) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, filename, file)
}
