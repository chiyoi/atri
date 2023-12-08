package pass

import (
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/atri/at"
)

const (
	Name        = "pass"
	Description = "Pass a message if you think it does not make sense or is not sent to you."
)

var Parameters = at.Parameters{
	Properties: struct{}{},
	Required:   []string{},
}

func Serve(o *[]at.ToolOutput, c *at.FunctionCall) {
	logs.Info("Passed.")
}
