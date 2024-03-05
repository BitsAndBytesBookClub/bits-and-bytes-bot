package main

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"regexp"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"calender-poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("creating a new poll")
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "meeting poll created!",
			},
		},
		); err != nil {
			slog.Error("error performing calender-poll interaction", err)
			return
		}
	},
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	slog.Info("message received", "message", m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	pattern := `(?i)\bbook\b.*?\bwe\b.*?\breading\b`
	r, err := regexp.Compile(pattern)
	if err != nil {
		slog.Error("Error compiling regex", err)
		return
	}

	if r.MatchString(m.Content) {
		response := "We're currently reading 'Designing Data Intensive Applications' by Martin Kleppmann. Join the discussion!"
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			slog.Error("Error sending message", err)
		}
	}
}
