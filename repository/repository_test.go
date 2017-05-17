package repository

import (
	"reflect"
	"testing"
	"time"
)

func TestFetchOne(t *testing.T) {
	r := New("./test-repo", "")
	p := r.FetchOne("2017/01/test-post")
	if p == nil {
		t.Fatal("Failed to fetch the post")
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	testCases := []struct {
		descr    string
		expected interface{}
		actual   interface{}
	}{
		{"Key", "2017/01/test-post", p.Key},
		{"Date", time.Date(2017, 1, 1, 0, 0, 0, 0, jst), *p.Date},
		{"Title", "Test Post", p.Title},
	}
	for _, tc := range testCases {
		matched := false
		switch tc.expected.(type) {
		case time.Time:
			matched = tc.expected.(time.Time).Equal(tc.actual.(time.Time))
		default:
			matched = reflect.DeepEqual(tc.expected, tc.actual)
		}
		if !matched {
			t.Errorf("%s: FetchOne() = %#v, exepcted %#v", tc.descr, tc.actual, tc.expected)
		}
	}
}

func TestFetch(t *testing.T) {
	r := New("./test-repo", "")
	ps := r.Fetch(nil)
	if len(ps) != 1 {
		t.Fatalf("len(ps) = %d, expected %d", len(ps), 1)
	}
	expected := "2017/01/test-post"
	if ps[0].Key != "2017/01/test-post" {
		t.Errorf("ps[0].Key = %#v, expected %#v", ps[0], expected)
	}
}
