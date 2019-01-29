package anpan

/* context.go:
 * Contains some utility functions for anpan.Context.
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// Reply directly replies with a message.
func (c *Context) Reply(message string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, message)
}

// ReplyEmbed directly replies with a embed, but not with a message.
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Channel.ID, embed)
}

// ReplyFile directly replies with a file, but not with a message.
func (c *Context) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, filename, file)
}
