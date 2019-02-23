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
	fmt.Println("anpan example bot")

	client, err := discordgo.New("Bot <your_token_here>")
	if err != nil {
		// Make sure you properly handle errors in your code ;)
		panic(err)
	}

	// Arguments are:
	// prefixes    - []string
	// owner ids   - []string
	// ignore bots - bool
	// check perms - bool

	prefixes := []string{
		"e!",
	}
	owners := []string{
		"your id",
		"another one",
	}
	handler := anpan.NewCommandHandler(prefixes, owners, true, true)

	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - int
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	handler.AddCommand("ping", "Check the bot's ping", false, false, 0, anpan.CommandTypeEverywhere, pingCommand)

	// Tell the dicordgo client to use handler's OnMessage function
	client.AddHandler(handler.OnMessage)
	// Optional: we can also have anpan control the playing status
	client.AddHandler(handler.StatusHandler.OnReady)

	err = client.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down")
	client.Close()
}

func pingCommand(ctx anpan.Context, _ []string) error {
	msg, _ := ctx.Reply("Pong!")

	sent, _ := msg.Timestamp.Parse()
	took := time.Now().Sub(sent)

	ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong! Ping took *%s*!", took.String()))
	return nil
}
