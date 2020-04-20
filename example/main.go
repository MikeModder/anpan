package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MikeModder/anpan"
	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Example bot for anpan. Version 1.1.0.")

	client, err := discordgo.New("Bot <your token here>")
	if err != nil {
		fmt.Printf("Creating a session failed: \"%s\".\n", err.Error())
		return
	}

	// Create a new handler and add a command.
	handler := anpan.New([]string{"e!"}, []string{"your id", "another one"}, true, true, true)
	handler.AddCommand("ping", "Check the bot's ping.", []string{"pong"}, false, false, discordgo.PermissionSendMessages, discordgo.PermissionSendMessages, anpan.CommandTypeEverywhere, pingCommand)
	handler.SetHelpCommand("help", []string{}, discordgo.PermissionSendMessages, helpCommand)

	handler.SetOnErrorFunc(func(context anpan.Context, name string, err error) {
		fmt.Printf("An error occurred for command \"%s\": \"%s\".\n", name, err.Error())
	})

	// Tell the dicordgo client to use handler's OnMessage function
	client.AddHandler(handler.OnMessage)

	if err = client.Open(); err != nil {
		fmt.Printf("Opening the session failed: \"%s\".\n", err.Error())
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down.")
	if err := client.Close(); err != nil {
		fmt.Printf("Closing the session failed. \"%s\". Ignoring.\n", err.Error())
	}
}

func pingCommand(ctx anpan.Context, args []string) error {
	msg, err := ctx.Reply("Pong!")
	if err != nil {
		return err
	}

	sent, err := msg.Timestamp.Parse()
	if err != nil {
		return err
	}

	_, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong! Ping took *%s*!", time.Now().Sub(sent).String()))
	return err
}
