package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/megawubs/calendar"
	"github.com/megawubs/wod_ical/renderers"
	"github.com/megawubs/wod_ical/wod"
	"github.com/uniplaces/carbon"
	"time"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	wods, err := wod.All(c.Param("apiKey"), carbon.Now().Time)
	if err != nil {
		c.Render(500, r.JSON(err))
	}
	nextMonthWods, err := wod.All(c.Param("apiKey"), carbon.Now().AddMonth().Time)
	if err != nil {
		c.Render(500, r.JSON(err))
	}
	wods = append(wods, nextMonthWods...)

	cal := calendar.Calendar{Version: "2.0", ProId: "wod_ical"}
	for _, w := range wods {
		format := "02-01-2006 15:04"
		start, err := time.Parse(format, w.DateStart)
		if err != nil {
			return fmt.Errorf("could not parse WOD start date: %s", err)
		}
		end, err := time.Parse(format, w.DateEnd)
		if err != nil {
			return fmt.Errorf("could not parse WOD end date: %s", err)
		}
		calendar.NewEvent(w.Id+2, "", start, end, w.Name)
	}
	err = wods.MarshallICalendar(&cal, time.Now().Location())
	if err != nil {
		return fmt.Errorf("could not marshal icalendar: %s", err)
	}

	return c.Render(200, renderers.ICAL(cal))
}
