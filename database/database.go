package database

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
)

type Database interface {
	SelectOne(id int) (*entity.Post, error)
	SelectNeighbors(id string) (*entity.Post, error)
	Select(selector *Selector) ([]*entity.Post, error)
	Insert(post *entity.Post) error
	Update(post *entity.Post) error
	Delete(post *entity.Post) error
}

type Selector struct {
	Label string
}

type DatabaseImpl struct {
	config *config.Config
}

func Init(config_ *config.Config) (Database, error) {
	postsDir := filepath.Join(config_.DatabaseDir, "posts")
	if !util.ExistsFile(postsDir) {
		if err := os.MkdirAll(postsDir, 0755); err != nil {
			return nil, err
		}
	}
	return &DatabaseImpl{config_}, nil
}

func (db *DatabaseImpl) SelectOne(id int) (*entity.Post, error) {
	postFile := filepath.Join(db.config.DatabaseDir, "posts", strconv.FormatInt(int64(id), 10))
	if !util.ExistsFile(postFile) {
		return nil, nil
	}
	body, err := ioutil.ReadFile(filepath.Join(db.config.DatabaseDir, "posts", strconv.FormatInt(int64(id), 10)))
	if err != nil {
		return nil, err
	}
	var result entity.Post
	if err := yaml.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (db *DatabaseImpl) SelectNeighbors(id string) (*entity.Post, error) {
	return nil, nil
}

func (db *DatabaseImpl) Select(selector *Selector) ([]*entity.Post, error) {
	return nil, nil
}

func (db *DatabaseImpl) Insert(post *entity.Post) error {
	lastId, err := db.getLastId()
	if err != nil {
		return err
	}
	post.Id = lastId + 1
	if err := db.Update(post); err != nil {
		return err
	}
	return nil
}

func (db *DatabaseImpl) Update(post *entity.Post) error {
	if post.Id == 0 {
		return errors.New("the post id must be set")
	}
	bs, err := yaml.Marshal(post)
	if err != nil {
		return err
	}
	path_ := filepath.Join(db.config.DatabaseDir, "posts", strconv.FormatInt(int64(post.Id), 10))
	if err := ioutil.WriteFile(path_, bs, 0644); err != nil {
		return err
	}
	return nil
}

func (db *DatabaseImpl) Delete(post *entity.Post) error {
	return nil
}

func (db *DatabaseImpl) getLastId() (int, error) {
	postIds, err := db.getPostIds()
	if err != nil {
		return 0, err
	}
	if len(postIds) <= 0 {
		return 0, nil
	}
	return postIds[0], nil
}

func (db *DatabaseImpl) getPostIds() ([]int, error) {
	postIds := []int{}
	files, err := ioutil.ReadDir(path.Join(db.config.DatabaseDir, "posts"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		id, err := strconv.ParseInt(file.Name(), 10, 0)
		if err != nil {
			return nil, err
		}
		postIds = append(postIds, int(id))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(postIds)))
	return postIds, nil
}
