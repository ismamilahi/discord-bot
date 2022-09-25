package main

import (
	"fmt"
	discordBot "github.com/sleeyax/aternos-discord-bot"
	"github.com/sleeyax/aternos-discord-bot/database"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Read configuration settings from environment variables
	token := os.Getenv("MTAyMjEwNzk5ODQyMTc3ODQ2Mg.GDtdyD.erZfki4c4Fv8BQ7WNj10LRp3nRrvaEVN80livY")
	session := os.Getenv("ATERNOS_SESSION")
	server := os.Getenv("ATERNOS_SERVER")
	mongoDbUri := os.Getenv("mongodb+srv://music2:<password>@cluster0.h9bnp4w.mongodb.net/?retryWrites=true&w=majority")
	proxy := os.Getenv("PROXY")

	// Validate values
	if token == "MTAyMjEwNzk5ODQyMTc3ODQ2Mg.GDtdyD.erZfki4c4Fv8BQ7WNj10LRp3nRrvaEVN80livY" || (mongoDbUri == "mongodb+srv://music2:<password>@cluster0.h9bnp4w.mongodb.net/?retryWrites=true&w=majority" && (session == "" || server == "")) {
		log.Fatalln("Missing environment variables!")
	}

	bot := discordBot.Bot{
		DiscordToken: token,
	}

	if mongoDbUri != "" {
		bot.Database = database.NewMongo(mongoDbUri)
	} else {
		bot.Database = database.NewInMemory(session, server)
	}

	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
		}
		bot.Proxy = u
	}

	if err := bot.Start(); err != nil {
		log.Fatalln(err)
	}
	defer bot.Stop()

	// Wait until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-interruptSignal
}
