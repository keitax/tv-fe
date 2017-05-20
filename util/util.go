package util

import (
	"html/template"
	"os"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

var metadataRe = regexp.MustCompile(`(?ms)^---\s*$\n(.*?)^---\s*$\n(.*)`)

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

func ConvertToStringSlice(orig []interface{}) []string {
	dest := []string{}
	for _, elem := range orig {
		dest = append(dest, elem.(string))
	}
	return dest
}

func StripMetadata(content string) (map[string]interface{}, string) {
	ms := metadataRe.FindStringSubmatch(content)
	if len(ms) < 3 {
		return map[string]interface{}{}, content
	}
	if len(ms) > 3 {
		panic("BUG: must not happen")
	}
	metadataSection, bodySection := ms[1], ms[2]
	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(metadataSection), &metadata); err != nil {
		return map[string]interface{}{}, bodySection
	}
	return metadata, bodySection
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
