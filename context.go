package anpan

/* context.go:
 * Contains some utility functions for anpan.Context
 *
 * Anpan (c) 2018 MikeModder/MikeModder007
 */

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Context) Reply(message string) {
	c.Session.ChannelMessageSend(c.Message.ChannelID, message)
}

func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) {
	c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}
