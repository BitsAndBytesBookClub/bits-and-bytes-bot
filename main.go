package main

import (
	"log/slog"
	"os"

	"github.com/devanbenz/bits-and-bytes-bot/database"
	"github.com/devanbenz/bits-and-bytes-bot/discord"
	"github.com/devanbenz/bits-and-bytes-bot/state"
	"github.com/joho/godotenv"
)

func main() {
	bootStrapLocalDev()

	state := state.NewState()
	pgDb := database.NewPostgresDb()
	bot := discord.NewDiscordBot(pgDb)

	bot.AddDiscordHandlers(state)
	if err := bot.StartDiscordBot(); err != nil {
		slog.Error("error opening botSession api connection", "error", err)
		return
	}
	defer pgDb.CloseDb()
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
