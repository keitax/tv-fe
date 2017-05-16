package repository

import (
	"io/ioutil"
	"path/filepath"
	"time"
	"gopkg.in/src-d/go-git.v4"

	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
)

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
	return nil
}
