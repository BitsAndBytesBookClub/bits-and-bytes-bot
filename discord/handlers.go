package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/devanbenz/bits-and-bytes-bot/state"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

func GenerateCalenderPoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Name != "calender-poll" {
		slog.Warn("not calender-poll", "data", data.Name)
		return
	}

	eventName := data.Options[0].StringValue()
	slog.Info("creating a new poll", "eventName", eventName, "dates", data.Options[1].StringValue())
	dates := strings.Split(data.Options[1].StringValue(), ",")

	for _, v := range dates {
		// 2024-03-15T10:00:00,2024-03-15T11:00:00,2024-03-15T12:00:00
		const shortForm = "2006-01-02T15:04:05"
		_, err := time.Parse(shortForm, v)
		if err != nil {
			slog.Error("cannot parse date times", "error", err)
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "please provide a proper list of dates!",
				},
			},
			); err != nil {
				slog.Error("error performing calender-poll interaction", err)
				return
			}
		}
	}

	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "Vote for meeting time",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					Emoji: discordgo.ComponentEmoji{
						Name: "🗓️",
					},
					CustomID: "vote_meeting_time",
				},
			},
		},
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("Please vote for a meeting time for **%s**", eventName),
			Components: components,
		},
	},
	); err != nil {
		slog.Error("error performing calender-poll interaction", "error", err)
		return
	}
}

func insertEmailForVoting(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modal := discordgo.ModalSubmitInteractionData{
		CustomID: "user_email_input",
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "email_input",
					Label:       "Your Email",
					Style:       discordgo.TextInputShort,
					Placeholder: "Enter your email address",
					Required:    true,
				},
			}},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "meeting_vote_form",
			Title:      "Vote for meeting time",
			Components: modal.Components,
		},
	})
	if err != nil {
		slog.Error("error handling voting button", "error", err)
		return
	}
}

func voteForMeeting(s *discordgo.Session, i *discordgo.InteractionCreate, state *state.State) {
	if i.Type != discordgo.InteractionModalSubmit {
		slog.Warn("not an interaction message component", "interaction", discordgo.InteractionMessageComponent)
		return
	}

	var formOptions []discordgo.SelectMenuOption

	for _, v := range state.Dates {
		formOptions = append(formOptions, discordgo.SelectMenuOption{
			Label:   v,
			Value:   v,
			Default: false,
			Emoji: discordgo.ComponentEmoji{
				Name: "🗓️",
			},
		})
	}

	form := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "date_selection",
					Placeholder: "Select the date you would like to meet",
					Options:     formOptions,
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "meeting_vote_form",
			Title:      "Vote for meeting time",
			Content:    "Please select an available meeting time",
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: form,
		},
	})
	if err != nil {
		slog.Error("error handling date selection", "error", err)
		return
	}
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

func completeVoting(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Post voting",
			Content: "Thank you for voting!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("error handling post voting response", "error", err)
		return
	}
}
