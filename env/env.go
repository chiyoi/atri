package env

import (
	"os"
)

const (
	Version = "Atri v0.1.0"
)

var (
	Database        = os.Getenv("DATABASE")
	Category        = os.Getenv("CATEGORY")
	AssistantIDAtri = os.Getenv("ASSISTANT_ID_ATRI")

	TokenDiscordApplication = os.Getenv("TOKEN_DISCORD_APPLICATION")
	OpenAiAPIKey            = os.Getenv("OPENAI_API_KEY")

	EndpointCosmos = os.Getenv("ENDPOINT_COSMOS")
)
