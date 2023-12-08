package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/atri/handlers/chat"
	"github.com/chiyoi/atri/handlers/ping"
)

func MessageCreate() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		switch {
		case ping.Serve(s, m):
		case chat.Serve(s, m):
		}
	}
}
