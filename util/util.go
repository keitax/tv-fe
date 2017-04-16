package util

import (
	"html/template"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func ParseMarkdown(text string) template.HTML {
	bs := blackfriday.MarkdownBasic([]byte(text))
	bs = bluemonday.UGCPolicy().SanitizeBytes(bs)
	return template.HTML(string(bs))
}
