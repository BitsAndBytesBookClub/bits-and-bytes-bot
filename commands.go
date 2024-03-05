package main

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
				Name:        "date",
				Description: "Date that meeting should take place; ex: 01/02/2024",
				Required:    true,
			},
		},
	},
}
