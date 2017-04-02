package view

import (
	"html/template"
	"io"

	"github.com/labstack/gommon/log"

	"github.com/keitax/textvid/config"
)

type View interface {
	RenderIndex(out io.Writer) error
	Render500(out io.Writer) error
}

type ViewImpl struct {
	config *config.Config
}

func New(config_ *config.Config) View {
	return &ViewImpl{config_}
}

func (v *ViewImpl) RenderIndex(out io.Writer) error {
	t, err := template.ParseFiles("view/templates/index.tmpl")
	if err != nil {
		return err
	}
	if err := t.Execute(out, map[string]interface{}{
	}); err != nil {
		return err
	}
	return nil
}

func (v *ViewImpl) Render500(out io.Writer) error {
	if _, err := out.Write([]byte("500 Internal Server Error")); err != nil {
		log.Error(err)
	}
	return nil
}
