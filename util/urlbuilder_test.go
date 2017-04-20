package util

import (
	"testing"
	"time"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
)

func TestLinkToPostPage(t *testing.T) {
	urlBuilder := &UrlBuilder{config: &config.Config{
		BaseUrl: "http://localhost/",
	}}
	createdAt, err := time.Parse(time.RFC3339, "2017-01-01T00:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	url := urlBuilder.LinkToPostPage(&entity.Post{
		CreatedAt: &createdAt,
		UrlName:   "hello-world",
	})
	expected := "http://localhost/2017/01/hello-world.html"
	if url != expected {
		t.Errorf("urlBuilder.LinkToPostView(_) = %q, expected %q", url, expected)
	}
}
