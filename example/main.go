package main

import (
	"fmt"
	"time"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Example bot for anpan.\nVersion 1.2.0.\nInitializing...")

	// Here we create an appropriate client.
	client, err := discordgo.New("Bot <your token here>")
	if err != nil {
		fmt.Printf("Creating a session failed: \"%s\".\n", err.Error())
		return
	}

	// In here we create a handler with the supplied data...
	handler := anpan.New([]string{"e!"}, []string{"your id", "another one"}, true, true, true)
	client.AddHandler(handler.MessageHandler)

	// ...then we register a command...
	handler.AddCommand("ping", "Check the bot's ping.", []string{"pong"}, false, false, discordgo.PermissionSendMessages, discordgo.PermissionSendMessages, anpan.CommandTypeEverywhere, pingCommand)

	// ...and a help command.
	handler.SetHelpCommand("help", []string{}, discordgo.PermissionSendMessages, discordgo.PermissionSendMessages, helpCommand)

	// And, of course, a function to let us know if something went wrong.
	handler.SetOnErrorFunc(func(context anpan.Context, command *anpan.Command, _ []string, err error) {
		fmt.Printf("An error occurred for command \"%s\": \"%s\".\n", command.Name, err.Error())
	})

	// Now, time to connect...
	if err = client.Open(); err != nil {
		fmt.Printf("Opening the session failed: \"%s\".\n", err.Error())
		return
	}

	// ...and wait until we need to exit.
	anpan.WaitForInterrupt()

	// Now we close the client, assuming it's still open.
	fmt.Println("Shutting down.")
	if err := client.Close(); err != nil {
		fmt.Printf("Closing the session failed. \"%s\". Ignoring.\n", err.Error())
	}

	// And we're done!
}

func pingCommand(ctx anpan.Context, _ []string) error {
	// First, we need a message...
	msg, err := ctx.Reply("Pong!")
	if err != nil {
		return err
	}

	// ...for a timestamp...
	sent, err := msg.Timestamp.Parse()
	if err != nil {
		return err
	}

	// ...to use some math for the final value.
	_, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong! Ping took *%s*!", time.Now().Sub(sent).String()))
	return err
}
