package main

import (
	"fmt"
	"strings"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func helpCommand(context anpan.Context, args []string, commands []*anpan.Command, prefixes []string) error {
	typeCheck := func(chn discordgo.ChannelType, cmd anpan.CommandType) bool {
		switch cmd {
		case anpan.CommandTypeEverywhere:
			return true

		case anpan.CommandTypePrivate:
			if chn == discordgo.ChannelTypeDM {
				return true
			}

			break

		case anpan.CommandTypeGuild:
			if chn == discordgo.ChannelTypeGuildText {
				return true
			}

			break
		}

		return false
	}

	if len(args) >= 1 {
		for _, commannd := range commands {
			if args[0] != commannd.Name {
				continue
			}

			// Basic checks.
			if commannd.Hidden || (context.Channel.Type == discordgo.ChannelTypeDM && commannd.Type == anpan.CommandTypeGuild) || (context.Channel.Type == discordgo.ChannelTypeGuildText && commannd.Type == anpan.CommandTypePrivate) {
				return nil
			}

			// Proper English is always fun.
			var (
				owneronlystring = "No"
				typestring      = "No restrictions"
			)

			if commannd.OwnerOnly {
				owneronlystring = "Yes"
			}

			switch commannd.Type {
			case anpan.CommandTypePrivate:
				typestring = "Private"
				break
			case anpan.CommandTypeGuild:
				typestring = "Guild-only"
				break
			}

			prefixesBuilder := strings.Builder{}
			if len(prefixes) == 1 {
				prefixesBuilder.WriteString(fmt.Sprintf("The bot's prefix is %s", prefixes[0]))
			} else {
				prefixesBuilder.WriteString("The bot's prefixes are ")
				for i, prefix := range prefixes {
					if i+1 == len(prefixes) {
						prefixesBuilder.WriteString(fmt.Sprintf("and %s", prefix))
					} else {
						prefixesBuilder.WriteString(fmt.Sprintf("%s, ", prefix))
					}
				}
			}

			_, err := context.ReplyEmbed(&discordgo.MessageEmbed{
				Title: "Help!",
				Description: fmt.Sprintf("Help for command `%s`\n Description: `%s`\nOwner only: **%s**\nType: **%s**", commannd.Name, commannd.Description,
					owneronlystring, typestring),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf(" %s.", prefixesBuilder.String()),
				},
			})

			return err
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: fmt.Sprintf("`%s` is not a valid command!", args[0]),
			Color:       0xff0000,
		}

		_, err := context.ReplyEmbed(embed)
		return err
	}

	count := len(commands)
	var list string

	for _, cmd := range commands {
		if !cmd.Hidden && typeCheck(context.Channel.Type, cmd.Type) {
			list += fmt.Sprintf("**%s** - `%s`\n", cmd.Name, cmd.Description)
		}
	}

	var footer strings.Builder

	// Grammar is always fun.
	if count == 1 {
		footer.WriteString("There is 1 command.")
	} else {
		footer.WriteString(fmt.Sprintf("There are %d commands.", count))
	}

	footer.WriteString(" | ")

	if len(prefixes) == 1 {
		footer.WriteString(fmt.Sprintf("The bot's prefix is %s.", prefixes[0]))
	} else {
		prefixesBuilder := strings.Builder{}

		for i, prefix := range prefixes {
			if i+1 == len(prefixes) {
				prefixesBuilder.WriteString(fmt.Sprintf("and %s", prefix))
			} else {
				prefixesBuilder.WriteString(fmt.Sprintf("%s, ", prefix))
			}
		}

		footer.WriteString(fmt.Sprintf("The bot's prefixes are %s.", prefixesBuilder.String()))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Commands:",
		Color:       0x08a4ff,
		Description: list,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer.String(),
		},
	}

	context.ReplyEmbed(embed)
	return nil
}
