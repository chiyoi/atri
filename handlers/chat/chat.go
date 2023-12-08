package chat

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/apricot/kitsune"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/atri/at"
	"github.com/chiyoi/atri/containers"
	"github.com/chiyoi/atri/contexts"
	"github.com/chiyoi/atri/env"
	"github.com/chiyoi/atri/functions"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/iter/opt"
	"github.com/chiyoi/iter/res"
)

func Serve(s *discordgo.Session, m *discordgo.MessageCreate) (block bool) {
	logs.Info("Ask chat assistant.")
	block = true
	channelID := m.ChannelID
	var reply *discordgo.Message
	reply, err := s.ChannelMessageSend(channelID, "[Auto Reply]アトリ、検索中ーー")
	if err != nil {
		logs.Error(err)
		return
	}
	defer func() {
		if err != nil {
			_, err := s.ChannelMessageEdit(channelID, reply.ID, "[Auto Reply]エラー発生。")
			if err != nil {
				logs.Error(err)
			}
		}
	}()

	c, err := s.Channel(channelID)
	if err != nil {
		logs.Error(err)
		return
	}
	if c.ParentID != env.CategoryChat {
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, at.ContextKeyAPIKey, env.TokenOpenAI)
	ctx = contexts.WithDiscordSessionMessage(ctx, s, m.Message)

	var threadID string
	cs, err := containers.Client(containers.Channels)
	err = res.C(cs, err, cosmos.PointRead(ctx, channelID, &threadID))
	if az.IsNotFound(err) {
		resp, err1 := at.ThreadsCreate(ctx, nil)
		threadID = resp.ID
		err = res.C(cs, err1, cosmos.PointUpsert(ctx, channelID, threadID))
	}
	if err != nil {
		logs.Error(err)
		return
	}

	_, err = at.ThreadsMessagesCreate(ctx, threadID, at.RoleUser, m.Content, nil)
	var re *kitsune.ResponseError
	if errors.As(err, &re) && re.StatusCode == http.StatusBadRequest {
		_, err = s.ChannelMessageEdit(channelID, reply.ID, "[Auto Reply]リクエスト出来なっかた。チャットちゃんお忙しい中かも。（該当メッセージは後ほど削除されます）")
		if err != nil {
			logs.Error(err)
		}
		time.Sleep(time.Second * 10)
		s.ChannelMessagesBulkDelete(channelID, []string{m.ID, reply.ID})
		if err != nil {
			logs.Error(err)
		}
		return
	}
	if err != nil {
		logs.Error(err)
		return
	}

	run, err := at.ThreadsRunsCreate(ctx, threadID, env.AssistantIDAtri, nil)
	if err != nil {
		logs.Error(err)
		return
	}

	replyLoading := func() func() {
		dots := func() func() string {
			count := 1
			return func() string {
				dots := strings.Repeat(".", count)
				if count == 3 {
					count = 1
				} else {
					count++
				}
				return dots
			}
		}()
		return func() {
			reply, err = s.ChannelMessageEdit(channelID, reply.ID, "[Auto Reply]チャットちゃん思考中"+dots())
			if err != nil {
				logs.Error(err)
			}
		}
	}()

	for {
		switch run.Status {
		case at.RunStatusQueued, at.RunStatusInProgress:
			replyLoading()
		case at.RunStatusRequiresAction:
			replyLoading()
			var step at.RunStep
			step, err = getLatestRunStep(ctx, threadID, run.ID)
			err = res.C(step.ID, err, handleRequireAction(ctx, threadID, run.ID))
			if err != nil {
				logs.Error(err)
				return
			}
		case at.RunStatusCompleted:
			var message at.Message
			message, err = getLatestMessage(ctx, threadID)
			if err != nil {
				logs.Error(err)
				return
			}
			reply, err = s.ChannelMessageEdit(channelID, reply.ID, message.Plaintext())
			if err != nil {
				logs.Error(err)
			}
			return
		case at.RunStatusCancelled, at.RunStatusCancelling:
			err := s.ChannelMessageDelete(m.ChannelID, reply.ID)
			if err != nil {
				logs.Error(err)
			}
			return
		case at.RunStatusExpired, at.RunStatusFailed:
			logs.Warning("Run stopped abnormally.")
			s.ChannelMessageEdit(channelID, reply.ID, "[Auto Reply]エラー発生。")
			return
		}

		time.Sleep(time.Second)
		run, err = at.ThreadsRunsRetrieve(ctx, threadID, run.ID)
		if err != nil {
			logs.Error(err)
			return
		}
	}
}

func handleRequireAction(ctx context.Context, threadID, runID string) func(stepID string) (err error) {
	return func(stepID string) (err error) {
		step, err := getLatestRunStep(ctx, threadID, runID)
		if err != nil {
			return
		}

		queued, err := containers.Client(containers.QueuedSteps)
		exist, err := res.R(queued, err, cosmos.PointExist(ctx, stepID))
		if err != nil {
			return
		}
		if exist {
			return
		}

		err = cosmos.PointUpsert(ctx, step.ID, nil)(queued)
		if err != nil {
			return
		}

		var outputs []at.ToolOutput
		for _, call := range step.StepDetails.ToolCalls {
			c := at.NewFunctionCall(ctx, threadID, runID, call)
			functions.Serve(&outputs, c)
		}
		if len(outputs) > 0 {
			_, err = at.ThreadsRunsToolOutputsSubmit(ctx, threadID, runID, outputs)
		} else {
			_, err = at.ThreadsRunsCancel(ctx, threadID, runID)
		}
		return
	}
}

func getLatestRunStep(ctx context.Context, threadID, runID string) (step at.RunStep, err error) {
	pager, err := at.ThreadsRunsStepsList(ctx, threadID, runID, at.QueryPager{
		Limit: 1,
	})
	if len(pager.Data) != 1 {
		err = opt.Or(err, errors.New("get latest run step error (len(pager.Data) != 1)"))
		return
	}
	step = pager.Data[0]
	return
}

func getLatestMessage(ctx context.Context, threadID string) (message at.Message, err error) {
	pager, err := at.ThreadsMessagesList(ctx, threadID, at.QueryPager{
		Limit: 1,
	})
	if len(pager.Data) != 1 {
		err = opt.Or(err, errors.New("get latest message error (len(pager.Data) != 1)"))
		return
	}
	message = pager.Data[0]
	return
}
