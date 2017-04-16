package util

import (
	"html/template"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	result := ParseMarkdown(`Hello
==

hello, world!
`)
	expected := template.HTML(`<h1>Hello</h1>

<p>hello, world!</p>
`)
	if expected != result {
		t.Errorf("Failed to parse:\nexpected:\n%v\nresult:\n%v", expected, result)
	}
}
