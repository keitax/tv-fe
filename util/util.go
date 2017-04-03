package util

import (
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

func ParseMarkdown(text string) string {
	bs := blackfriday.MarkdownBasic([]byte(text))
	bs = bluemonday.UGCPolicy().SanitizeBytes(bs)
	return string(bs)
}
