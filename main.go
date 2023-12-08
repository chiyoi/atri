package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/atri/env"
	"github.com/chiyoi/atri/handlers"
)

func main() {
	s, err := discordgo.New("Bot " + env.TokenDiscordApplication)
	if err != nil {
		logs.Panic(err)
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages
	s.AddHandler(handlers.MessageCreate())

	logs.Info("Open discord session.")
	if err := s.Open(); err != nil {
		logs.Panic(err)
	}
	defer s.Close()

	neko.Block()
}
