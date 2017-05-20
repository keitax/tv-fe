package repository

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
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
	path := filepath.Join(r.localGitRepoPath, "posts", key+".md")
	if !util.ExistsFile(path) {
		return nil
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	meta, body := util.StripMetadata(string(bs))
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
		Labels: util.ConvertToStringSlice(meta["labels"].([]interface{})),
	}
}

func (r *Repository) Fetch(pq *dao.PostQuery) []*entity.Post {
	ps := []*entity.Post{}
	if err := filepath.Walk(filepath.Join(r.localGitRepoPath, "posts"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !postFileRe.MatchString(path) {
			return nil
		}
		key := postFileRe.FindStringSubmatch(path)[1]
		ps = append(ps, r.FetchOne(key))
		return nil
	}); err != nil {
		panic(err)
	}
	sort.Sort(entity.SortPost(ps))
	return ps
}