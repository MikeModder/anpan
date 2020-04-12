package anpan

import (
	"fmt"
	"strings"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func helpCommand(context Context, args []string, commands []*anpan.Command, prefixes []string) error {
	typeCheck := func(chn discordgo.ChannelType, cmd CommandType) bool {
		switch cmd {
		case CommandTypeEverywhere:
			return true

		case CommandTypePrivate:
			if chn == discordgo.ChannelTypeDM {
				return true
			}

			break

		case CommandTypeGuild:
			if chn == discordgo.ChannelTypeGuildText {
				return true
			}

			break
		}

		return false
	}

	if len(args) >= 1 {
		for _, v := range commands {
			if args[0] != v.Name {
				continue
			}

			// Basic checks.
			if commannd.Hidden || (context.Channel.Type == discordgo.ChannelTypeDM && commannd.Type == CommandTypeGuild) || (context.Channel.Type == discordgo.ChannelTypeGuildText && commannd.Type == CommandTypePrivate) {
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
			case CommandTypePrivate:
				typestring = "Private"
				break
			case CommandTypeGuild:
				typestring = "Guild-only"
				break
			}

			_, err := context.ReplyEmbed(&discordgo.MessageEmbed{
				Title: "Help!",
				Description: fmt.Sprintf("Help for command `%s`\n Description: `%s`\nOwner only: **%s**\nType: **%s**", commannd.Name, commannd.Description,
					owneronlystring, typestring),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("The bot's prefixes are %s.", handler.Prefixes),
				},
			})
			return err
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: fmt.Sprintf("`%s` is not a valid command!", args[0]),
			Color:       0xff0000,
		}

		context.ReplyEmbed(embed)
		return nil
	}

	count := len(handler.Commands)
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
