package database

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
)

func TestSelectOne(t *testing.T) {
	db := prepareTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	post, err := db.SelectOne(99)
	if err != nil {
		t.Fatal(err)
	}
	if post != nil {
		t.Error("Post 99 doesn't exist but a post was selected")
	}
}

func TestInsert(t *testing.T) {
	db := prepareTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	if err := db.Insert(&entity.Post{
		UrlLabel: "test-url-label",
		Title:    "test-title",
		Body:     "test-body",
		Labels:   []string{"test-label"},
	}); err != nil {
		t.Fatal(err)
	}

	inserted, err := db.SelectOne(1)
	if err != nil {
		t.Fatal(err)
	}

	expected := &entity.Post{
		Id:       1,
		UrlLabel: "test-url-label",
		Title:    "test-title",
		Body:     "test-body",
		Labels:   []string{"test-label"},
	}
	if !reflect.DeepEqual(expected, inserted) {
		t.Errorf("Failed to insert the post: expected: %v, inserted: %v", expected, inserted)
	}
}

func TestUpdate(t *testing.T) {
	db := prepareTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	originalPost := &entity.Post{
		UrlLabel: "original-url-label",
		Title:    "original-title",
		Body:     "original-body",
		Labels:   []string{"original-label"},
	}
	if err := db.Insert(originalPost); err != nil {
		t.Fatal(err)
	}

	if err := db.Update(&entity.Post{
		Id:       originalPost.Id,
		UrlLabel: "updated-url-label",
		Title:    "updated-title",
		Body:     "updated-body",
		Labels:   []string{"updated-label"},
	}); err != nil {
		t.Fatal(err)
	}

	updated, err := db.SelectOne(originalPost.Id)
	if err != nil {
		t.Fatal(err)
	}
	expected := &entity.Post{
		Id:       originalPost.Id,
		UrlLabel: "updated-url-label",
		Title:    "updated-title",
		Body:     "updated-body",
		Labels:   []string{"updated-label"},
	}
	if !reflect.DeepEqual(expected, updated) {
		t.Errorf("failed to update the post: expected: %v, updated: %v", expected, updated)
	}
}

func prepareTestDatabase(t *testing.T) Database {
	dir, err := ioutil.TempDir(os.TempDir(), "database_test-")
	if err != nil {
		t.Fatal(err)
	}
	dbDir, err := Init(&config.Config{DatabaseDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	return dbDir
}

func cleanupTestDatabase(t *testing.T, db Database) {
	if err := os.RemoveAll(db.(*DatabaseImpl).config.DatabaseDir); err != nil {
		t.Fatal(err)
	}
}
