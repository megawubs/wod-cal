package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/megawubs/wod_ical/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
