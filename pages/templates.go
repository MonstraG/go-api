package pages

import (
	"embed"
	"html/template"
)

//go:embed *.gohtml **/*.gohtml
var templates embed.FS

func mustParseTemplates(patterns ...string) *template.Template {
	return template.Must(template.ParseFS(templates, patterns...))
}

func ParsePage(patterns ...string) *template.Template {
	all := append([]string{"base.gohtml"}, patterns...)
	return mustParseTemplates(all...)
}

func ParsePartial(patterns ...string) *template.Template {
	return mustParseTemplates(patterns...)
}
