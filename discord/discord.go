package discord

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/devanbenz/bits-and-bytes-bot/database"
	"github.com/devanbenz/bits-and-bytes-bot/state"
)

type Bot struct {
	*discordgo.Session
	db database.Db
}

func NewDiscordBot(db database.Db) *Bot {
	botSession, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_API_SECRET")))
	if err != nil {
		slog.Error("error creating botSession api connection", err)
	}

	return &Bot{botSession, db}
}

func (botSession *Bot) StartDiscordBot() error {
	if err := botSession.Open(); err != nil {
		slog.Error("error opening websocket connection", err)
		return fmt.Errorf("error opening websocket connection %s", err)
	}

	slog.Info("botSession bot generating commands...")
	for _, v := range commands {
		_, err := botSession.ApplicationCommandCreate(botSession.State.User.ID, os.Getenv("GUILD_ID"), v)
		if err != nil {
			slog.Error("Cannot creating command", "name", v.Name, "error", err)
			return fmt.Errorf("error generating commands %s", err)
		}
	}

	slog.Info("discord bot listening...")

	return nil
}

func (botSession *Bot) CloseDiscordBot() {
	err := botSession.Close()
	if err != nil {
		slog.Error("error closing botSession api", err)
	}
}

func (botSession *Bot) AddDiscordHandlers(state *state.State) {
	botSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	botSession.AddHandler(messageCreate)
	botSession.Identify.Intents = discordgo.IntentsGuildMessages
	botSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "vote_meeting_time" {
			user := i.Member.User

			if slices.Contains(state.Voters, user.Username) {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You have already voted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})

				if err != nil {
					slog.Error("error handling voting button", "error", err)
				}

				return
			}

			insertEmailForVoting(s, i)
		} else if i.Type == discordgo.InteractionModalSubmit {
			email := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

			if slices.Contains(state.Emails, email) {
				slog.Info("email already exists", "email", email)
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This email already voted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})

				if err != nil {
					slog.Error("error performing email interaction", err)
				}

				return
			}

			state.Emails = append(state.Emails, email)

			slog.Info("email added to responses", "email", email)

			voteForMeeting(s, i, state)
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "date_selection" {
			date := i.MessageComponentData().Values[0]
			state.Votes[date] = state.Votes[date] + 1
			state.Voters = append(state.Voters, i.Member.User.Username)

			slog.Info("vote cast for next meeting", "date", date, "voteCount", state.Votes[date])

			completeVoting(s, i)
		} else if i.ApplicationCommandData().Name == "calendar-poll" {
			state.PollDuration = time.Duration(i.ApplicationCommandData().Options[2].IntValue()) * time.Second
			state.Dates = strings.Split(i.ApplicationCommandData().Options[1].StringValue(), ",")
			state.Votes = make(map[string]int, len(state.Dates))
			slog.Info("poll created with duration", "duration", state.PollDuration)

			for _, v := range state.Dates {
				state.Votes[v] = 0
			}

			if state.PollDuration > 0 {
				slog.Info("sending poll timer start channel", "duration", state.PollDuration)
				state.StartPollTimer <- true
			}

			GenerateCalendarPoll(s, i)
		}
	})
}

func PollFinished(state *state.State) {
	go func() {
		for {
			enabled := <-state.StartPollTimer
			if enabled && state.PollDuration > 0 {
				slog.Info("poll timer has been enabled", "enabled", state.StartPollTimer, "duration", state.PollDuration)
				time.AfterFunc(state.PollDuration, func() {
					slog.Info("sending payload", "emails", state.Emails, "dates", state.Dates, "votes", state.Votes)
					// TODO: Send data to data store and calendly API
					// TODO: Create a botSession event when poll finishes
					fmt.Println(fmt.Sprintf("emails: %s", state.Emails))
					fmt.Println(fmt.Sprintf("dates: %s", state.Dates))
					fmt.Println(fmt.Sprintf("votes: %v", state.Votes))
					fmt.Println(fmt.Sprintf("voters: %s", state.Voters))

					slog.Info("resetting state")
					state.ResetState()

					fmt.Println(fmt.Sprintf("emails: %s", state.Emails))
					fmt.Println(fmt.Sprintf("dates: %s", state.Dates))
					fmt.Println(fmt.Sprintf("votes: %v", state.Votes))
					fmt.Println(fmt.Sprintf("voters: %s", state.Voters))
				})
			}
		}
	}()
}

func TerminateOnSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
