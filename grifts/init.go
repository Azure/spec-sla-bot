package grifts

import (
	"github.com/Azure/spec_sla_bot/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
