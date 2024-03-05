package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type State struct {
	Dates        []string
	Emails       []string
	Votes        map[string]int
	PollDuration time.Duration
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
		return
	}
	state := State{}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_API_SECRET")))
	if err != nil {
		slog.Error("error creating discord api connection", err)
		return
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println(fmt.Sprintf("dates: %s", state.Dates))
		fmt.Println(fmt.Sprintf("emails: %s", state.Emails))
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "vote_meeting_time" {
			insertEmailForVoting(s, i)
		} else if i.Type == discordgo.InteractionModalSubmit {
			email := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			state.Emails = append(state.Emails, email)
			slog.Info("email added to responses", "email", email)
			voteForMeeting(s, i, state)
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "date_selection" {
			date := i.MessageComponentData().Values[0]
			state.Votes[date] = state.Votes[date] + 1
			slog.Info("vote cast for next meeting", "date", date, "voteCount", state.Votes[date])
			completeVoting(s, i)
		} else if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			state.PollDuration = time.Duration(i.ApplicationCommandData().Options[2].IntValue()) * time.Second
			state.Dates = strings.Split(i.ApplicationCommandData().Options[1].StringValue(), ",")
			state.Votes = make(map[string]int, len(state.Dates))
			for _, v := range state.Dates {
				state.Votes[v] = 0
			}

			h(s, i)
		}
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
		_, err = discord.ApplicationCommandCreate(discord.State.User.ID, "1214625037244432465", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	time.AfterFunc(state.PollDuration, func() {
		fmt.Println(fmt.Sprintf("emails: %s", state.Emails))
		fmt.Println(fmt.Sprintf("dates: %s", state.Dates))
		fmt.Println(fmt.Sprintf("votes: %s", state.Votes))
		//reset the state
		state = State{}
	})

	slog.Info("discord bot listening...")
	terminateOnSignal()
}

func terminateOnSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
