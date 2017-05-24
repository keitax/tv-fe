package textvid

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/keitax/textvid/entity"
	"gopkg.in/src-d/go-git.v4"
)

var postFileRe = regexp.MustCompile(`^.*([0-9][0-9][0-9][0-9]/[0-9][0-9]/.+)\.md$`)

type Repository struct {
	localGitRepoPath  string
	remoteGitRepoPath string
	gitRepo           *git.Repository
}

func New(localGitRepoPath, remoteGitRepoPath string) *Repository {
	return &Repository{
		localGitRepoPath:  localGitRepoPath,
		remoteGitRepoPath: remoteGitRepoPath,
	}
}

func (r *Repository) FetchOne(key string) *entity.Post {
	ps := r.Fetch(&PostQuery{
		Start:   1,
		Results: 0,
	})

	var found *entity.Post
	var foundIdx int
	for i, p := range ps {
		if key == p.Key {
			found = p
			foundIdx = i
			break
		}
	}
	if found == nil {
		return nil
	}

	nextIdx := foundIdx - 1
	prevIdx := foundIdx + 1
	if 0 <= nextIdx && nextIdx < len(ps) {
		found.NextPost = ps[nextIdx]
	}
	if 0 <= prevIdx && prevIdx < len(ps) {
		found.PreviousPost = ps[prevIdx]
	}

	return found
}

func (r *Repository) Fetch(pq *PostQuery) []*entity.Post {
	ps := []*entity.Post{}

	if err := filepath.Walk(filepath.Join(r.localGitRepoPath, "posts"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !postFileRe.MatchString(path) {
			return nil
		}
		key := postFileRe.FindStringSubmatch(path)[1]
		ps = append(ps, r.loadPost(key))
		return nil
	}); err != nil {
		panic(err)
	}

	sort.Sort(entity.SortPost(ps))

	start := Min(len(ps), Max(0, int(pq.Start)-1))
	ps = ps[start:]
	if pq.Results >= 1 {
		end := Min(len(ps), Max(1, int(pq.Results)))
		ps = ps[:end]
	}
	return ps
}

func (r *Repository) Commit(p *entity.Post) {
	panic("not implemented")
}

func (r *Repository) loadPost(key string) *entity.Post {
	path := filepath.Join(r.localGitRepoPath, "posts", key+".md")
	if !ExistsFile(path) {
		panic("Faile to load the post")
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	meta, body := StripFrontMatter(string(bs))
	d_, err := time.Parse("2006-01-02 15:04:05 Z07:00", meta["date"].(string))
	d := &d_
	if err != nil {
		d = nil
	}
	return &entity.Post{
		Key:    key,
		Date:   d,
		Title:  meta["title"].(string),
		Body:   body,
		Labels: ConvertToStringSlice(meta["labels"].([]interface{})),
	}
}
