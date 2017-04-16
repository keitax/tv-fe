package view

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/keitax/textvid/config"
	"github.com/Sirupsen/logrus"
)

type View interface {
	RenderIndex(out io.Writer) error
	Render500(out io.Writer) error
}

type view struct {
	config *config.Config
}

func New(config_ *config.Config) View {
	return &view{config_}
}

func (v *view) RenderIndex(out io.Writer) error {
	return v.renderTemplate("index.tmpl", out, map[string]interface{}{})
}

func (v *view) Render500(out io.Writer) error {
	if _, err := out.Write([]byte("500 Internal Server Error")); err != nil {
		logrus.Error(err)
	}
	return nil
}

func (v *view) renderTemplate(templateName string, out io.Writer, context map[string]interface{}) error {
	t, err := template.ParseFiles(filepath.Join(v.config.TemplateDir, templateName))
	if err != nil {
		return err
	}
	context_ := map[string]interface{}{
		"SiteTitle": v.config.SiteTitle,
	}
	for key, value := range context {
		context_[key] = value
	}
	if err := t.Execute(out, context_); err != nil {
		return err
	}
	return nil
}
