package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
		return
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_API_SECRET")))
	if err != nil {
		slog.Error("error creating discord api connection", err)
		return
	}

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err = discord.Open(); err != nil {
		slog.Error("error opening websocket connection", err)
		return
	}
	defer func(discord *discordgo.Session) {
		err = discord.Close()
		if err != nil {
			slog.Error("error closing discord api", err)
		}
	}(discord)

	slog.Info("discord bot generating commands...")
	for _, v := range commands {
		_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "788203854038564886", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	slog.Info("discord bot listening...")
	terminateOnSignal()
}

func terminateOnSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
