package tvfe

import (
	"testing"
)

func TestLinkToTopPage(t *testing.T) {
	urlBuilder := &URLBuilder{config: &Config{
		BaseURL: "http://localhost/",
	}}
	url := urlBuilder.LinkToTopPage()
	expected := "http://localhost/"
	if url != expected {
		t.Errorf("urlBuilder.LinkToTopPage() = %q, expected: %q", url, expected)
	}
}

func TestLinkToPostPage(t *testing.T) {
	urlBuilder := &URLBuilder{config: &Config{
		BaseURL: "http://localhost/",
	}}
	url := urlBuilder.LinkToPostPage(&Post{
		Key: "2017/01/hello-world",
	})
	expected := "http://localhost/2017/01/hello-world.html"
	if url != expected {
		t.Errorf("urlBuilder.LinkToPostView(_) = %q, expected %q", url, expected)
	}
}

func TestLinkToPostListPage(t *testing.T) {
	ub := &URLBuilder{config: &Config{
		BaseURL: "http://localhost/",
	}}
	u := ub.LinkToPostListPage(&PostQuery{
		Start:   1,
		Results: 10,
	})
	expected := "http://localhost/posts/?start=1&results=10"
	if u != expected {
		t.Errorf("ub.LinkToPostListPage(_) = %q, expected %q", u, expected)
	}
}
