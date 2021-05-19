package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func helpCommand(context anpan.Context, args []string, commands []*anpan.Command, prefixes []string) error {
	// This is a check for a command to make sure only appropriate commands are shown.
	// TODO: Add a permission check
	typeCheck := func(chn discordgo.ChannelType, cmd anpan.CommandType) bool {
		switch cmd {
		case anpan.CommandTypeEverywhere:
			return true

		case anpan.CommandTypePrivate:
			if chn == discordgo.ChannelTypeDM {
				return true
			}

		case anpan.CommandTypeGuild:
			if chn == discordgo.ChannelTypeGuildText {
				return true
			}
		}

		return false
	}

	// Here we found out that the user inquires for one specific command.
	if len(args) >= 1 {
		for _, command := range commands {
			// If this command isn't it, continue searching.
			if args[0] != command.Name {
				continue
			}

			// If this command is supposed to stay hidden, or if this isn't the correct place, stop the command.
			if command.Hidden || !typeCheck(context.Channel.Type, command.Type) {
				return nil
			}

			// Some useful declarations.
			var (
				owneronlystring = "No"
				typestring      = "Anywhere"
			)

			// Is this command only accessible by the owners?
			if command.OwnerOnly {
				owneronlystring = "Yes"
			}

			// What type do we actually have?
			switch command.Type {
			case anpan.CommandTypePrivate:
				typestring = "Private"

			case anpan.CommandTypeGuild:
				typestring = "Guild-only"
			}

			// Time to tell the user about the prefixes.
			prefixesBuilder := strings.Builder{}
			if len(prefixes) == 1 {
				prefixesBuilder.WriteString(fmt.Sprintf("The prefix is %s", prefixes[0]))
			} else {
				prefixesBuilder.WriteString("The prefixes are ")

				for i, prefix := range prefixes {
					if i+1 == len(prefixes) {
						prefixesBuilder.WriteString(fmt.Sprintf("and %s", prefix))
					} else {
						prefixesBuilder.WriteString(fmt.Sprintf("%s, ", prefix))
					}
				}
			}

			// Maybe our command has a few nicknames...
			aliases := "**None.**"
			if len(command.Aliases) > 0 {
				aliases = strings.Join(command.Aliases, "`, `")
				aliases = "`" + aliases + "`"
			}

			// ..but anyway! Time to return the message.
			_, err := context.ReplyEmbed(&discordgo.MessageEmbed{
				Title:       "Help",
				Color:       0x08a4ff,
				Description: fmt.Sprintf("**%s**\nAliases: %s\nOwner only: **%s**\nUsable: **%s**", command.Description, aliases, owneronlystring, typestring),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf(" %s.", prefixesBuilder.String()),
				},
			})

			// We're done.
			return err
		}

		// We've not found anything :(
		_, err := context.Reply("Command `" + args[0] + "` doesn't exist.")
		return err
	}

	// Well, we now know the user wants to know what commands we actually have.
	var (
		count          int
		commandsSorted = make([]*anpan.Command, len(commands))
		embed          = &discordgo.MessageEmbed{
			Title: "Commands",
			Color: 0x08a4ff,
		}
		names = make([]string, len(commands))
	)

	// Get all names...
	for i, cmd := range commands {
		names[i] = cmd.Name
	}

	// ...sort them alphabetically...
	sort.Strings(names)

	// ...and arrange the commands accordingly.
	for i, v := range names {
		for _, v2 := range commands {
			if v2.Name == v {
				commandsSorted[i] = v2
				break
			}
		}

		if commandsSorted[i] == nil {
			return fmt.Errorf("sort failure")
		}
	}

	// Now that we've sorted the commands, we can show them to the user.
	for _, cmd := range commandsSorted {
		if !cmd.Hidden && typeCheck(context.Channel.Type, cmd.Type) {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   cmd.Name,
				Value:  cmd.Description,
				Inline: count%2 == 0,
			})

			count++
		}
	}

	// We want a footer for additional information.
	var footer strings.Builder

	// How many commands do we have?
	if count == 1 {
		footer.WriteString("There is 1 command.")
	} else {
		footer.WriteString(fmt.Sprintf("There are %d commands.", count))
	}

	footer.WriteString(" | ")

	if len(prefixes) == 1 {
		footer.WriteString(fmt.Sprintf("The prefix is %s.", prefixes[0]))
	} else {
		prefixesBuilder := strings.Builder{}

		for i, prefix := range prefixes {
			if i+1 == len(prefixes) {
				prefixesBuilder.WriteString(fmt.Sprintf("and %s", prefix))
			} else {
				prefixesBuilder.WriteString(fmt.Sprintf("%s, ", prefix))
			}
		}

		footer.WriteString(fmt.Sprintf("The prefixes are %s.", prefixesBuilder.String()))
	}

	// Let them know about the prefixes.
	embed.Footer = &discordgo.MessageEmbedFooter{Text: footer.String()}

	// Time to give them help.
	_, err := context.ReplyEmbed(embed)
	return err
}
