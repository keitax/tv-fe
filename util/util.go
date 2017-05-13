package util

import (
	"html/template"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

const metadataSeparator = "---"

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
	lineno := 0
	lines := strings.Split(content, "\n")

	isEof := func() bool {
		return lineno >= len(lines)
	}
	parseMetadataSeparator := func() bool {
		if isEof() {
			return false
		}
		if lines[lineno] == metadataSeparator {
			lineno++
			return true
		}
		return false
	}
	parseMetadataBody := func() map[string]interface{} {
		ls := []string{}
		for !isEof() {
			if lines[lineno] == metadataSeparator {
				break
			}
			ls = append(ls, lines[lineno])
			lineno++
		}
		metadata := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(strings.Join(ls, "\n")), &metadata); err != nil {
			return metadata
		}
		return metadata
	}
	parseBody := func() string {
		ls := []string{}
		for !isEof() {
			ls = append(ls, lines[lineno])
			lineno++
		}
		return strings.Join(ls, "\n")
	}

	emptyMetadata := map[string]interface{}{}
	if !parseMetadataSeparator() {
		return emptyMetadata, parseBody()
	}
	metadata := parseMetadataBody()
	if !parseMetadataSeparator() {
		return emptyMetadata, parseBody()
	}
	return metadata, parseBody()
}
