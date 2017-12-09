package textvid

import (
	"html/template"
	"os"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

var metadataRe = regexp.MustCompile(`(?ms)^---\s*$\n(.*?)^---\s*$\n(.*)`)

// ExistsFile checks whether the file exists or not.
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

// ParseMarkdown parses the Markdown text and makes a HTML document.
func ParseMarkdown(text string) template.HTML {
	bs := blackfriday.Run([]byte(text))
	bs = bluemonday.UGCPolicy().SanitizeBytes(bs)
	return template.HTML(string(bs))
}

// ConvertToStringSlice converts the slice of interface{} to a slice of strings.
func ConvertToStringSlice(orig []interface{}) []string {
	dest := []string{}
	for _, elem := range orig {
		dest = append(dest, elem.(string))
	}
	return dest
}

// StripFrontMatter parses and strips the front matter section of the input text.
func StripFrontMatter(content string) (map[string]interface{}, string) {
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

// Max returns the max of x and y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the min of x and y.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
