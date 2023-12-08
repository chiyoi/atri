package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/apricot/logs"
)

func Serve(s *discordgo.Session, m *discordgo.MessageCreate) (block bool) {
	switch m.Content {
	case "Ping", "ping":
		logs.Info("Pong!")
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			logs.Error(err)
		}
		return true
	case "Pong", "pong":
		logs.Info("Ping?")
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping?")
		if err != nil {
			logs.Error(err)
		}
		return true
	}
	return false
}
