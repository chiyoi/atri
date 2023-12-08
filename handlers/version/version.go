package version

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/atri/env"
)

func Serve(s *discordgo.Session, m *discordgo.MessageCreate) (block bool) {
	switch m.Content {
	case "Version", "version":
		logs.Info("Get version.")
		_, err := s.ChannelMessageSend(m.ChannelID, "[Auto Reply]"+env.Version)
		if err != nil {
			logs.Error(err)
		}
	}
	return true
}
