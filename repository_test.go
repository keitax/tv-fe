package textvid

import (
	"os"
	"reflect"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func prepareTestRepository(t *testing.T) *Repository {
	g, err := git.PlainInit("./test-repo", false)
	if err == git.ErrRepositoryAlreadyExists {
		return openTestRepository(t)
	}
	if err != nil {
		t.Fatal(err)
	}
	w, err := g.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "textvid-test",
			Email: "textvid-test@dummy.jp",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return openTestRepository(t)
}

func openTestRepository(t *testing.T) *Repository {
	r, err := OpenRepository("./test-repo", "")
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func cleanupTestRepository(t *testing.T) {
	repoData := "./test-repo/.git"
	_, err := os.Stat(repoData)
	if err == os.ErrNotExist {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(repoData); err != nil {
		t.Fatal(err)
	}
}

func TestFetchOne(t *testing.T) {
	r := prepareTestRepository(t)
	defer cleanupTestRepository(t)
	p := r.FetchOne("2017/01/test-post-01")
	if p == nil {
		t.Fatal("Failed to fetch the post")
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	testCases := []struct {
		descr    string
		expected interface{}
		actual   interface{}
	}{
		{"Key", "2017/01/test-post-01", p.Key},
		{"Date", time.Date(2017, 1, 1, 0, 0, 0, 0, jst), *p.Date},
		{"Title", "Test Post 1", p.Title},
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

func TestFetchOneGetsNeighbors(t *testing.T) {
	r := prepareTestRepository(t)
	defer cleanupTestRepository(t)
	testCases := []struct {
		descr    string
		iKey     string
		oNextKey string
		oPrevKey string
	}{
		{"Get neightbors of the post in the middle", "2017/01/test-post-02", "2017/01/test-post-03", "2017/01/test-post-01"},
		{"Get neightbors of the first post", "2017/01/test-post-03", "", "2017/01/test-post-02"},
		{"Get neightbors of the last post", "2017/01/test-post-01", "2017/01/test-post-02", ""},
	}
	for _, tc := range testCases {
		p := r.FetchOne(tc.iKey)
		var nk, pk string
		if p.NextPost != nil {
			nk = p.NextPost.Key
		}
		if p.PreviousPost != nil {
			pk = p.PreviousPost.Key
		}
		if nk != tc.oNextKey {
			t.Errorf("%s: NextPost.Key = %v, expected %s", tc.descr, nk, tc.oNextKey)
		}
		if pk != tc.oPrevKey {
			t.Errorf("%s: PreviousPost.Key = %v, expected %s", tc.descr, pk, tc.oPrevKey)
		}
	}
}

func TestFetchAcceptsRangeQuery(t *testing.T) {
	r := prepareTestRepository(t)
	defer cleanupTestRepository(t)
	testCases := []struct {
		descr    string
		iStart   uint64
		iResults uint64
		oKeys    []string
	}{
		{"Fetch first two posts", 1, 2, []string{"2017/01/test-post-03", "2017/01/test-post-02"}},
		{"Fetch last two posts", 2, 2, []string{"2017/01/test-post-02", "2017/01/test-post-01"}},
		{"Start is too large", 99, 2, []string{}},
		{"Start is too small", 0, 1, []string{"2017/01/test-post-03"}},
		{"Results is too large", 1, 99, []string{"2017/01/test-post-03", "2017/01/test-post-02", "2017/01/test-post-01"}},
		{"Results is zero (all results)", 1, 0, []string{"2017/01/test-post-03", "2017/01/test-post-02", "2017/01/test-post-01"}},
	}
	for _, tc := range testCases {
		ps := r.Fetch(&PostQuery{
			Start:   tc.iStart,
			Results: tc.iResults,
		})
		keys := []string{}
		for _, p := range ps {
			keys = append(keys, p.Key)
		}
		if !reflect.DeepEqual(keys, tc.oKeys) {
			t.Errorf("%s: keys = %v, expected %v", tc.descr, keys, tc.oKeys)
		}
	}
}
