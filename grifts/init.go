package grifts

import (
	"github.com/Azure/spec-sla-bot/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
