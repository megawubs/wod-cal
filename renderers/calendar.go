package renderers

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/megawubs/calendar"
	"io"
)

type icalRenderer struct {
	value calendar.Calendar
}

func (i icalRenderer) ContentType() string {
	return "text/calendar"
}

func (i icalRenderer) Render(w io.Writer, d render.Data) error {
	i.value.Write(w)
	return nil
}

func ICAL(c calendar.Calendar) render.Renderer {
	return icalRenderer{value: c}
}
