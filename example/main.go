package main

import (
	"fmt"
	"time"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func exampleOnErrorFunc(context anpan.Context, command *anpan.Command, content []string, err error) {
	if err == anpan.ErrCommandNotFound {
		return
	}

	fmt.Printf("An error occurred for command \"%s\": \"%s\".\n", command.Name, err.Error())
}

func examplePrerunFunc(context anpan.Context, command *anpan.Command, content []string) bool {
	fmt.Printf("Command \"%s\" is being run by \"%s#%s\" (ID: %s).\n", command.Name, context.User.Username, context.User.Discriminator, context.User.ID)
	return true
}

func main() {
	fmt.Println("Example bot for anpan.\nVersion 1.2.0.\nInitializing...")

	// Here we create an appropriate client.
	client, err := discordgo.New("Bot <your token here>")
	if err != nil {
		fmt.Printf("Creating a session failed: \"%s\".\n", err.Error())
		return
	}

	// In here we create a handler with the supplied data...
	handler := anpan.New([]string{"e!"}, []string{"your id", "another one"}, client.StateEnabled, true, true, true, examplePrerunFunc, exampleOnErrorFunc, nil)
	client.AddHandler(handler.MessageHandler)

	// ...then we register a command...
	handler.AddCommand("ping", "Check the bot's ping.", []string{"pong"}, false, false, discordgo.PermissionSendMessages, discordgo.PermissionSendMessages, anpan.CommandTypeEverywhere, pingCommand)

	// ...and a help command.
	handler.SetHelpCommand("help", []string{}, discordgo.PermissionSendMessages, discordgo.PermissionSendMessages, helpCommand)

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
	// We need to know what time it is now.
	timestamp := time.Now()

	msg, err := ctx.Reply("Pong!")
	if err != nil {
		return err
	}

	// Now we can compare it to the current time to see how much time went away during the process of sending a message.
	_, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong! Ping took **%dms**!", time.Since(timestamp).Milliseconds()))
	return err
}
