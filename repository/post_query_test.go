package repository

import (
	"reflect"
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	iq := &PostQuery{
		Start:   6,
		Results: 5,
		Year:    2017,
		Month:   time.January,
		UrlName: "test",
	}
	oq := iq.Next()
	expected := &PostQuery{
		Start:   1,
		Results: 5,
		Year:    2017,
		Month:   time.January,
		UrlName: "test",
	}
	if !reflect.DeepEqual(oq, expected) {
		t.Errorf("iq.Next() = %v, expected %v", oq, expected)
	}
}

func TestPrevious(t *testing.T) {
	iq := &PostQuery{
		Start:   6,
		Results: 5,
		Year:    2017,
		Month:   time.January,
		UrlName: "test",
	}
	oq := iq.Previous()
	expected := &PostQuery{
		Start:   11,
		Results: 5,
		Year:    2017,
		Month:   time.January,
		UrlName: "test",
	}
	if !reflect.DeepEqual(oq, expected) {
		t.Errorf("iq.Previous() = %v, expected %v", oq, expected)
	}
}
