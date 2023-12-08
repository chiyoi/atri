package functions

import (
	"github.com/chiyoi/atri/at"
	"github.com/chiyoi/atri/functions/pass"
	"github.com/chiyoi/atri/functions/pick_loli"
)

func Serve(o *[]at.ToolOutput, c *at.FunctionCall) {
	switch c.ToolCall.Function.Name {
	case pass.Name:
		pass.Serve(o, c)
	case pick_loli.Name:
		pick_loli.Serve(o, c)
	}
}
