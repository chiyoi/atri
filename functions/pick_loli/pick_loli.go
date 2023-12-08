package pick_loli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/chiyoi/apricot/kitsune"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/atri/at"
	"github.com/chiyoi/atri/contexts"
)

const (
	EndpointLolicon = "https://api.lolicon.app/setu/v2"
)

const (
	Name        = "pick_loli"
	Description = "Randomly pick a cute girl (loli) illustration, returns a description containing its URL, title and author."
)

var Parameters = at.Parameters{
	Properties: struct{}{},
	Required:   []string{},
}

func Serve(o *[]at.ToolOutput, c *at.FunctionCall) {
	logs.Info("Pick loli!")
	ctx := c.Context()
	image, err := Pick(ctx)
	if err != nil {
		logs.Error(err)
		return
	}
	*o = append(*o, at.ToolOutput{
		ToolCallID: c.ToolCall.ID,
		Output:     image.Description,
	})
	handleGuildReply(ctx, image)
}

func handleGuildReply(ctx context.Context, image Image) {
	defer image.Body.Close()
	s, m, ok := contexts.GetDiscordSessionMessage(ctx)
	if !ok {
		return
	}
	_, err := s.ChannelFileSend(m.ChannelID, image.Filename, image.Body)
	if err != nil {
		logs.Error(err)
	}
}

func Pick(ctx context.Context) (image Image, err error) {
	attachLoliconQuery := func(r *http.Request) (*http.Request, error) {
		q := r.URL.Query()
		q.Set("r18", "2")
		q.Set("proxy", "0")
		r.URL.RawQuery = q.Encode()
		return r, nil
	}

	attachPixivReferer := func(r *http.Request) (*http.Request, error) {
		r.Header.Set("Referer", "https://www.pixiv.net/")
		return r, nil
	}

	var resp struct {
		Data []struct {
			Title  string `json:"title"`
			PID    int    `json:"pid"`
			Author string `json:"author"`
			UID    int    `json:"uid"`
			URLs   struct {
				Original string `json:"original"`
			} `json:"urls"`
		} `json:"data"`
	}
	err = kitsune.Get(ctx, &resp, attachLoliconQuery)(EndpointLolicon)
	switch {
	case err != nil:
		return
	case len(resp.Data) == 0:
		err = errors.New("empty response from lolicon api")
		return
	}

	data := resp.Data[0]
	u := data.URLs.Original
	image.Filename = path.Base(u)
	image.Description = fmt.Sprintf("URL: %s\nTitle: %s[%d]\nAuthor: %s[%d]", u, data.Title, data.PID, data.Author, data.UID)
	body, err := kitsune.GetStream(ctx, attachPixivReferer)(u)
	image.Body = body
	return
}

type Image struct {
	Body        io.ReadCloser
	Filename    string
	Description string
}
