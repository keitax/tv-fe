package textvid

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var postFileRe = regexp.MustCompile(`^.*([0-9][0-9][0-9][0-9]/[0-9][0-9]/.+)\.md$`)

type Repository struct {
	localGitRepoPath  string
	remoteGitRepoPath string
	gitRepo           *git.Repository
}

func OpenRepository(localGitRepoPath, remoteGitRepoPath string) (*Repository, error) {
	logrus.Infof("Try to open the local repository %s.", localGitRepoPath)
	r, err := git.PlainOpen(localGitRepoPath)
	if err == git.ErrRepositoryNotExists {
		logrus.Infof("There are no local repository, clone the remote repository: %s -> %s", remoteGitRepoPath, localGitRepoPath)
		var err error
		r, err = git.PlainClone(localGitRepoPath, false, &git.CloneOptions{
			URL: remoteGitRepoPath,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to clone the remote repository %s: %s", remoteGitRepoPath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Failed to open the local repository: %s", err)
	}
	logrus.Infof("Succeeded to open the repository.")
	return &Repository{
		localGitRepoPath:  localGitRepoPath,
		remoteGitRepoPath: remoteGitRepoPath,
		gitRepo:           r,
	}, nil
}

func (r *Repository) FetchOne(key string) *Post {
	ps := r.Fetch(&PostQuery{
		Start:   1,
		Results: 0,
	})

	var found *Post
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

func (r *Repository) Fetch(pq *PostQuery) []*Post {
	ps := []*Post{}

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

	sort.Sort(SortPost(ps))

	start := Min(len(ps), Max(0, int(pq.Start)-1))
	ps = ps[start:]
	if pq.Results >= 1 {
		end := Min(len(ps), Max(1, int(pq.Results)))
		ps = ps[:end]
	}
	return ps
}

func (r *Repository) Commit(p *Post) {
	panic("not implemented")
}

func (r *Repository) loadPost(key string) *Post {
	ref, err := r.gitRepo.Head()
	if err != nil {
		panic(err)
	}
	c, err := r.gitRepo.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}
	f, err := c.File(fmt.Sprintf("posts/%s.md", key))
	if err == object.ErrFileNotFound {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	cs, err := f.Contents()
	if err != nil {
		panic(err)
	}
	meta, body := StripFrontMatter(cs)
	d_, err := time.Parse("2006-01-02 15:04:05 Z07:00", meta["date"].(string))
	d := &d_
	if err != nil {
		d = nil
	}
	return &Post{
		Key:    key,
		Date:   d,
		Title:  meta["title"].(string),
		Body:   body,
		Labels: ConvertToStringSlice(meta["labels"].([]interface{})),
	}
}
