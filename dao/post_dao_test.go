package dao

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"

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

	inserted, err := d.SelectOne(original.Id)
	if err != nil {
		t.Fatal(err)
	}

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

func TestSelectByQuery(t *testing.T) {
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

func prepareDao(t *testing.T) PostDao {
	db, err := sql.Open("mysql", "keitax/keitax@tcp(localhost:3306)/test?parseTime=True")
	if err != nil {
		t.Fatal(err)
	}
	return NewPostDao(db)
}

func (pd *postDao) cleanup(t *testing.T) {
	if _, err := pd.db.Exec("TRUNCATE TABLE POST"); err != nil {
		t.Fatal(err)
	}
	if _, err := pd.db.Exec("TRUNCATE TABLE LAST_ID"); err != nil {
		t.Fatal(err)
	}
	if _, err := pd.db.Exec("INSERT INTO LAST_ID (POST_LAST_ID) VALUES (0)"); err != nil {
		t.Fatal(err)
	}
}
