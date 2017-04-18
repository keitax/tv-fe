package view

import (
	"html/template"
	"io"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
)

type View interface {
	RenderIndex(out io.Writer, posts []*entity.Post) error
	Render500(out io.Writer) error
}

type view struct {
	config *config.Config
}

func New(config_ *config.Config) View {
	return &view{config_}
}

func (v *view) RenderIndex(out io.Writer, posts []*entity.Post) error {
	return v.renderTemplate("post_list.tmpl", out, map[string]interface{}{
		"posts": posts,
	})
}

func (v *view) Render500(out io.Writer) error {
	if _, err := out.Write([]byte("500 Internal Server Error")); err != nil {
		logrus.Error(err)
	}
	return nil
}

func (v *view) renderTemplate(templateName string, out io.Writer, context map[string]interface{}) error {
	ts := template.New("root").Funcs(template.FuncMap{
		"RenderMarkdown": util.ParseMarkdown,
		"ShowTime": func(t time.Time) string {
			return t.Format("Jan. 02, 2006, 3:04 PM")
		},
	})
	ts = template.Must(ts.ParseFiles(
		filepath.Join(v.config.TemplateDir, "layout.tmpl"),
		filepath.Join(v.config.TemplateDir, templateName),
	))
	context_ := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
	}
	for key, value := range context {
		context_[key] = value
	}
	if err := ts.ExecuteTemplate(out, "layout.tmpl", context_); err != nil {
		return err
	}
	return nil
}
