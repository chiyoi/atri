package env

import (
	"os"

	"github.com/chiyoi/apricot/logs"
)

const (
	EndpointCosmos = "https://neko03cosmos.documents.azure.com:443/"
)

var (
	AssistantIDAtri string
	CategoryChat    string

	Database                string
	TokenDiscordApplication string
	TokenOpenAI             string
)

func init() {
	switch os.Getenv("ENV") {
	case "prod":
		Prod()
	default:
		Dev()
	}
}

func common() {
	TokenDiscordApplication = os.Getenv("TOKEN_DISCORD_APPLICATION")
	TokenOpenAI = os.Getenv("OPENAI_API_KEY")
	AssistantIDAtri = "asst_aTI20AjVAwpCli9qAn7uceNs"
}

func Dev() {
	common()
	os.Setenv("VERSION", "Yuxiu v0.1.0-dev")
	Database = "neko0001"
	// Actually Category Neko0001
	CategoryChat = "1180855047073054750"
	logs.SetLevel(logs.LevelDebug)
}

func Prod() {
	common()
	os.Setenv("VERSION", "Yuxiu v0.1.0")
	Database = "atri"
	CategoryChat = "1180508717167419432"
}
