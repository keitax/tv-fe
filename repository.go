package textvid

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var postFileRe = regexp.MustCompile(`^.*([0-9][0-9][0-9][0-9]/[0-9][0-9]/.+)\.md$`)

type Repository struct {
	gitRepo   *git.Repository
	postCache map[string]*Post
	mutex     *sync.RWMutex
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
		gitRepo: r,
		mutex:   &sync.RWMutex{},
	}, nil
}

func (r *Repository) SynchronizeRemote() {
	logrus.Info("Pull the remote repository.")
	if err := r.gitRepo.Pull(&git.PullOptions{}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return
		}
		panic(err)
	}
	r.UpdateCache()
}

func (r *Repository) UpdateCache() {
	logrus.Info("Update post cache.")

	fi, err := r.getHeadCommit().Files()
	if err != nil {
		panic(err)
	}

	gitAddedMap := r.collectGitAdded()
	cache := map[string]*Post{}
	for f, err := fi.Next(); err != io.EOF; f, err = fi.Next() {
		if err != nil {
			panic(err)
		}
		if postFileRe.MatchString(f.Name) {
			key := postFileRe.FindStringSubmatch(f.Name)[1]
			p := r.loadPost(key)
			if p.Date == nil {
				p.Date = gitAddedMap[postKeyToFilePath(key)]
			}
			cache[key] = p
		}
	}

	r.mutex.Lock()
	r.postCache = cache
	r.mutex.Unlock()

	var np *Post
	for _, p := range r.getPostList() {
		if np != nil {
			np.PreviousPost = p
		}
		p.NextPost = np
		np = p
	}
}

func (r *Repository) FetchOne(key string) *Post {
	return r.refPostCache()[key]
}

func (r *Repository) Fetch(pq *PostQuery) []*Post {
	ps := r.getPostList()
	start := Min(len(ps), Max(0, int(pq.Start) - 1))
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
	cs := r.readFile(postKeyToFilePath(key))
	fm, body := StripFrontMatter(cs)
	p := &Post{
		Key:  key,
		Body: body,
	}
	if fm["title"] != nil {
		p.Title = fm["title"].(string)
	}
	if fm["labels"] != nil {
		ls, ok := fm["labels"].([]string)
		if ok {
			p.Labels = ls
		}
	}
	if fm["date"] != nil {
		d, err := time.Parse("2006-01-02 15:04:05 Z07:00", fm["date"].(string))
		if err == nil {
			p.Date = &d
		}
	}
	return p
}

func (r *Repository) refPostCache() map[string]*Post {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.postCache
}

func (r *Repository) getPostList() []*Post {
	ps := []*Post{}
	for _, p := range r.refPostCache() {
		ps = append(ps, p)
	}
	sort.Sort(PostList(ps))
	return ps
}

func (r *Repository) getHeadCommit() *object.Commit {
	ref, err := r.gitRepo.Head()
	if err != nil {
		panic(err)
	}
	c, err := r.gitRepo.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}
	return c
}

func (r *Repository) readFile(filePath string) string {
	ref, err := r.gitRepo.Head()
	if err != nil {
		panic(err)
	}
	c, err := r.gitRepo.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}
	f, err := c.File(filePath)
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
	return cs
}

func (r *Repository) collectGitAdded() map[string]*time.Time {
	ts := map[string]*time.Time{}
	cc := make(chan *object.Commit, 2)
	cc <- r.getHeadCommit()
	for len(cc) > 0 {
		c := <-cc
		fi, err := c.Files()
		if err != nil {
			panic(err)
		}
		t := c.Author.When
		for f, err := fi.Next(); err != io.EOF; f, err = fi.Next() {
			if err != nil {
				panic(err)
			}
			ts[f.Name] = &t
		}
		pi := c.Parents()
		for p, err := pi.Next(); err != io.EOF; p, err = pi.Next() {
			if err != nil {
				panic(err)
			}
			cc <- p
		}
	}
	return ts
}

func postKeyToFilePath(key string) string {
	return fmt.Sprintf("posts/%s.md", key)
}
