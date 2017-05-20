package util

import (
	"html/template"
	"reflect"
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

func TestStripMetadata(t *testing.T) {
	testCases := []struct {
		descr            string
		input            string
		expectedMetadata map[string]interface{}
		expectedBody     string
	}{
		{
			descr: "valid content",
			input: `---
title: Hello, World
labels:
  - hello
  - world
---
hello, world
`,
			expectedMetadata: map[string]interface{}{
				"title": "Hello, World",
				"labels": []interface{}{
					"hello",
					"world",
				},
			},
			expectedBody: "hello, world\n",
		},
		{
			descr: "broken metadata (no first separator)",
			input: `title: Hello, World
---
hello, world
`,
			expectedMetadata: map[string]interface{}{},
			expectedBody: `title: Hello, World
---
hello, world
`,
		},
		{
			descr: "broken metadata (no second separator)",
			input: `---
title: Hello, World
hello, world
`,
			expectedMetadata: map[string]interface{}{},
			expectedBody: `---
title: Hello, World
hello, world
`,
		},
		{
			descr: "broken metadata (broken yaml)",
			input: `---
title: Hello
 labels:
---
hello, world
`,
			expectedMetadata: map[string]interface{}{},
			expectedBody:     "hello, world\n",
		},
	}

	for _, tc := range testCases {
		metadata, body := StripMetadata(tc.input)
		if !reflect.DeepEqual(metadata, tc.expectedMetadata) {
			t.Errorf("%s: StripMetadata() = %#v, _, expected %#v", tc.descr, metadata, tc.expectedMetadata)
		}
		if body != tc.expectedBody {
			t.Errorf("%s: StripMetadata() = _, %#v, expected %#v", tc.descr, body, tc.expectedBody)
		}
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		descr    string
		x        int
		y        int
		expected int
	}{
		{"equal", 5, 5, 5},
		{"x > y", 6, 5, 6},
		{"x < y", 5, 6, 6},
	}
	for _, tc := range testCases {
		actual := Max(tc.x, tc.y)
		if actual != tc.expected {
			t.Errorf("Max(%d, %d) = %d, expected %d", tc.x, tc.y, actual, tc.expected)
		}
	}
}
