package table

import (
	"html/template"
	"io"
	"net/url"
)

var tableTemplate *template.Template

func get() *template.Template {
	if tableTemplate == nil {
		tableTemplate = template.Must(
			template.ParseFiles("templates/table.template.html"))
	}
	return tableTemplate
}

// Render the form as a html table
func Render(form url.Values, w io.Writer) {
	get().Execute(w, form)
}
