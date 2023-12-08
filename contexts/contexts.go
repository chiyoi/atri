package contexts

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func WithDiscordSessionMessage(ctx context.Context, s *discordgo.Session, m *discordgo.Message) context.Context {
	ctx = context.WithValue(ctx, ContextKeyDiscordSession, s)
	ctx = context.WithValue(ctx, ContextKeyGuildMessage, m)
	return ctx
}

func GetDiscordSessionMessage(ctx context.Context) (s *discordgo.Session, m *discordgo.Message, ok bool) {
	s, ok1 := ctx.Value(ContextKeyDiscordSession).(*discordgo.Session)
	m, ok2 := ctx.Value(ContextKeyGuildMessage).(*discordgo.Message)
	ok = ok1 && ok2
	return
}

type ContextKey string

const (
	ContextKeyDiscordSession ContextKey = "Discord Session"
	ContextKeyGuildMessage   ContextKey = "Guild Message"
)
