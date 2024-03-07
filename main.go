package main

import (
	"github.com/devanbenz/bits-and-bytes-bot/discord"
	"github.com/devanbenz/bits-and-bytes-bot/state"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	bootStrapLocalDev()

	state := state.NewState()
	bot := discord.NewDiscordBot()

	bot.AddDiscordHandlers(state)
	if err := bot.StartDiscordBot(); err != nil {
		slog.Error("error opening botSession api connection", "error", err)
		return
	}
	defer bot.CloseDiscordBot()

	discord.PollFinished(state)
	discord.TerminateOnSignal()
}

func bootStrapLocalDev() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading .env file", err)
			return
		}
	}
	return
}
