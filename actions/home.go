package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/megawubs/calendar"
	"github.com/megawubs/go-wod/wod"
	"github.com/megawubs/wod_ical/renderers"
	"time"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	wods, err := wod.All(c.Param("apiKey"))
	if err != nil {
		c.Render(500, r.JSON(err))
	}
	cal := calendar.Calendar{Version: "2.0", ProId: "wod_ical"}
	wods.MarshallICalendar(&cal, time.Now().Location())

	return c.Render(200, renderers.ICAL(cal))
}
