package main

import (
	"github.com/devanbenz/bits-and-bytes-bot/discord"
	"github.com/devanbenz/bits-and-bytes-bot/models"
	"github.com/joho/godotenv"
	"log/slog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
		return
	}

	state := models.NewState()
	bot := discord.NewDiscordBot()

	bot.AddDiscordHandlers(state)
	if err = bot.StartDiscordBot(); err != nil {
		slog.Error("error opening botSession api connection", "error", err)
		return
	}
	defer bot.CloseDiscordBot()

	discord.PollFinished(state)
	discord.TerminateOnSignal()
}
