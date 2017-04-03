package util

import "testing"

func TestParseMarkdown(t *testing.T) {
	result := ParseMarkdown(`Hello
==

hello, world!
`)
	expected := `<h1>Hello</h1>

<p>hello, world!</p>
`
	if expected != result {
		t.Errorf("Failed to parse:\nexpected:\n%v\nresult:\n%v", expected, result)
	}
}
