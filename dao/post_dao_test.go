package dao

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
)

func TestInsertAndSelectOne(t *testing.T) {
	d := prepareDao(t)
	defer d.(*postDao).cleanup(t)

	original := &entity.Post{
		UrlName: "test-url-name",
		Labels:  []string{"test-label"},
		Title:   "test-title",
		Body:    "test-body",
	}
	if err := d.Insert(original); err != nil {
		t.Fatal(err)
	}

	inserted := d.SelectOne(original.Id)

	testCases := []struct {
		descr    string
		expected interface{}
		actual   interface{}
	}{
		{"Id", int64(1), inserted.Id},
		{"UrlLabel", "test-url-name", inserted.UrlName},
		{"Title", "test-title", inserted.Title},
		{"Body", "test-body", inserted.Body},
	}

	for _, tc := range testCases {
		if !reflect.DeepEqual(tc.expected, tc.actual) {
			t.Errorf("Failed to insert or to select: %s: expected: %v, actual: %v", tc.descr, tc.expected, tc.actual)
		}
	}
}

func TestSelectOneWithNeighbors(t *testing.T) {
	d := prepareDao(t)
	defer d.(*postDao).cleanup(t)

	for i := 0; i < 3; i++ {
		if err := d.Insert(&entity.Post{
			UrlName: "test-url-name",
			Title:   "test-title",
		}); err != nil {
			t.Fatal(err)
		}
	}

	p2 := d.SelectOne(2)
	expectedNextId := int64(3)
	if p2.NextPost.Id != expectedNextId {
		t.Errorf("p1.NextPost.Id = %d, expected %d", p2.NextPost.Id, expectedNextId)
	}
	expectedPreviousId := int64(1)
	if p2.PreviousPost.Id != expectedPreviousId {
		t.Errorf("p1.PreviousPost.Id = %d, expected %d", p2.NextPost.Id, expectedNextId)
	}

	p1 := d.SelectOne(1)
	if p1.PreviousPost != nil {
		t.Errorf("p1.PreviousPost = %v, expected nil", p1.PreviousPost)
	}
	if p1.NextPost == nil {
		t.Errorf("p1.NextPost = %v, expected non-nil", p1.NextPost)
	}

	p3 := d.SelectOne(3)
	if p3.PreviousPost == nil {
		t.Errorf("p3.PreviousPost = %v, expected non-nil", p3.PreviousPost)
	}
	if p3.NextPost != nil {
		t.Errorf("p3.NextPost = %v, expected nil", p3.NextPost)
	}
}

func TestSelectByQueryToSelectByRange(t *testing.T) {
	d := prepareDao(t)
	defer d.(*postDao).cleanup(t)

	for i := 0; i < 3; i++ {
		if err := d.Insert(&entity.Post{
			UrlName: "test-url-name",
			Title:   "test-title",
			Body:    "test-body",
		}); err != nil {
			t.Fatal(err)
		}
	}

	ps, err := d.SelectByQuery(&PostQuery{
		Start:   1,
		Results: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(ps) != 2 {
		t.Errorf("Selected post count isn't expected: expected: 2, actual: %d", len(ps))
	}
	if ps[0].Id != 3 {
		t.Errorf("First post id isn't expected: expected: 3, actual: %d", ps[0].Id)
	}
	if ps[1].Id != 2 {
		t.Errorf("Second post id isn't expected: expected: 2, actual: %d", ps[1].Id)
	}
}

func TestSelectByQueryToSelectByMonth(t *testing.T) {
	d := prepareDao(t)
	defer d.(*postDao).cleanup(t)

	createdAtList := []string{
		"2016-12-31T23:59:59+00:00",
		"2017-01-01T00:00:00+00:00",
		"2017-01-31T23:59:59+00:00",
		"2017-02-01T00:00:00+00:00",
	}
	for _, createdAt := range createdAtList {
		if err := d.Insert(makePostSpecificCreatedAt(t, createdAt)); err != nil {
			t.Fatal(err)
		}
	}

	ps, err := d.SelectByQuery(&PostQuery{
		Start:   1,
		Results: 10,
		Year:    2017,
		Month:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 2 {
		t.Fatalf("len(ps) = %d", len(ps))
	}
	for _, p := range ps {
		if !(p.CreatedAt.Year() == 2017 && p.CreatedAt.Month() == 1) {
			t.Errorf("p.CreatedAt = %d/%d", p.CreatedAt.Year(), p.CreatedAt.Month())
		}
	}
}

func TestSelectByQueryToSelectByUrlName(t *testing.T) {
	d := prepareDao(t)
	defer d.(*postDao).cleanup(t)

	urlNames := []string{
		"foo",
		"bar",
		"foobar",
	}
	for _, urlName := range urlNames {
		if err := d.Insert(&entity.Post{
			UrlName: urlName,
		}); err != nil {
			t.Fatal(err)
		}
	}

	ps, err := d.SelectByQuery(&PostQuery{
		Start:   1,
		Results: 10,
		UrlName: "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 {
		t.Fatalf("len(ps) = %d", len(ps))
	}
	if !reflect.DeepEqual("bar", ps[0].UrlName) {
		t.Errorf("%q.UrlName = %q", ps[0], ps[0].UrlName)
	}
}

func makePostSpecificCreatedAt(t *testing.T, createdAt string) *entity.Post {
	ti, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		t.Fatal(err)
	}
	return &entity.Post{
		CreatedAt: &ti,
		Title:     "Some Title",
	}
}

func prepareDao(t *testing.T) PostDao {
	conn, err := dbr.Open("mysql", "keitax/keitax@tcp(localhost:3306)/test?parseTime=True&loc=UTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	return NewPostDao(conn, &config.Config{
		Locale: "UTC",
	})
}

func (pd *postDao) cleanup(t *testing.T) {
	if _, err := pd.conn.Exec("TRUNCATE TABLE POST"); err != nil {
		t.Fatal(err)
	}
	if _, err := pd.conn.Exec("TRUNCATE TABLE LAST_ID"); err != nil {
		t.Fatal(err)
	}
	if _, err := pd.conn.Exec("INSERT INTO LAST_ID (POST_LAST_ID) VALUES (0)"); err != nil {
		t.Fatal(err)
	}
}
