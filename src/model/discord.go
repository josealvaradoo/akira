package model

import (
	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	Name         string
	Event        func(user *discordgo.InteractionCreate) LotteryResponse
	ForEveryone  bool
	IsAttachment bool
}

type Command struct {
	Name         string
	Description  string
	Event        func(user *discordgo.InteractionCreate) LotteryResponse
	ForEveryone  bool
	IsAttachment bool
}
