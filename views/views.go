package views

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/views/components"
)

//go:embed layouts/*
var htmlTemplates embed.FS

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"inputField": func(label, name, fieldType, placeholder string) components.InputField {
		return components.InputField{Label: label, Name: name, Type: fieldType, Placeholder: placeholder}
	},
	"options": func(args ...string) []components.Option {
		options := make([]components.Option, 0, len(args)/2)
		for i := 0; i < len(args); i += 2 {
			options = append(options, components.Option{Key: args[i], Value: args[i+1]})
		}
		return options
	},
	"selectBox": func(label, name string, options []components.Option) components.SelectBox {
		return components.SelectBox{Label: label, Name: name, Options: options}
	},
	"autocomplete": func(label, name, fieldType, placeholder, getUrl string) components.Autocomplete {
		return components.Autocomplete{Label: label, Name: name, Type: fieldType, Placeholder: placeholder, GetUrl: getUrl}
	},
	"primaryButton": func(label string) components.PrimaryButton {
		return components.PrimaryButton{Label: label}
	},
	"emptySlice": func() []string {
		return []string{}
	},
}).ParseFS(htmlTemplates, "layouts/*.html"))

func NewView(layout string) *View {
	return &View{
		layout: layout,
	}
}

type View struct {
	layout string
}

func (v *View) Render(ctx context.Context, w http.ResponseWriter, data any) {
	err := tmpl.ExecuteTemplate(w, v.layout, data)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("could not render template %v", v.Layout()), "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *View) Layout() string {
	return v.layout
}
