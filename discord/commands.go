package discord

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "calender-poll",
		Description: "Creates a new calender poll",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "event-name",
				Description: "Name of the meeting.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date-times",
				Description: "List of date/time options seperated by a comma. (format: 2024-03-05T08:00:00).",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "poll-duration",
				Description: "Duration of the poll.",
				Required:    true,
			},
		},
	},
}
