package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/devanbenz/bits-and-bytes-bot/discord"
	"github.com/devanbenz/bits-and-bytes-bot/models"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
		return
	}

	state := models.NewState()

	botSession, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_API_SECRET")))
	if err != nil {
		slog.Error("error creating botSession api connection", err)
		return
	}

	discord.AddDiscordHandlers(botSession, state)
	if err = discord.StartDiscordBot(botSession); err != nil {
		slog.Error("error opening botSession api connection", "error", err)
		return
	}

	slog.Info("discord bot listening...")

	discord.PollFinished(state)
	discord.TerminateOnSignal()
}
