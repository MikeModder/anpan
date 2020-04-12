package anpan

/* help.go:
 * Contains the default help command for anpan.
 *
 * anpan (c) 2020 MikeModder/MikeModder007
 */

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func helpCommand(context Context, args []string, handler CommandHandler) error {
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
		if commannd, has := handler.Commands[args[0]]; has {
			// Check if the command is hidden.
			if commannd.Hidden {
				return nil
			}

			// basic channel check
			if context.Channel.Type == discordgo.ChannelTypeDM && commannd.Type == CommandTypeGuild {
				return nil
			} else if context.Channel.Type == discordgo.ChannelTypeGuildText && commannd.Type == CommandTypePrivate {
				return nil
			}

			// Proper English is always fun.
			var (
				owneronlystring string
				typestring      string
			)

			if commannd.OwnerOnly {
				owneronlystring = "Yes"
			} else {
				owneronlystring = "No"
			}

			if commannd.Type == CommandTypePrivate {
				typestring = "Private"
			} else if commannd.Type == CommandTypeGuild {
				typestring = "Guild-only"
			} else {
				typestring = "No restrictions"
			}

			embed := &discordgo.MessageEmbed{
				Title: "Help!",
				Description: fmt.Sprintf("Help for command `%s`\n Description: `%s`\nOwner only: **%s**\nType: **%s**", commannd.Name, commannd.Description,
					owneronlystring, typestring),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("The bot's prefixes are %s.", handler.Prefixes),
				},
			}

			context.ReplyEmbed(embed)
			return nil
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

	for name := range handler.Commands {
		cmd := handler.Commands[name]
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

	if len(handler.Prefixes) == 1 {
		footer.WriteString(fmt.Sprintf("The bot's prefix is %s.", handler.Prefixes[0]))
	} else {
		prefixes := strings.Builder{}

		for i, prefix := range handler.Prefixes {
			if i+1 == len(handler.Prefixes) {
				prefixes.WriteString(fmt.Sprintf("and %s", prefix))
			} else {
				prefixes.WriteString(fmt.Sprintf("%s, ", prefix))
			}
		}

		footer.WriteString(fmt.Sprintf("The bot's prefixes are %s.", prefixes.String()))
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
