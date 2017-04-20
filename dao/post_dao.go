package dao

import (
	"errors"
	"time"

	"github.com/gocraft/dbr"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
)

type PostDao interface {
	SelectOne(id int64) (*entity.Post, error)
	SelectByQuery(query *PostQuery) ([]*entity.Post, error)
	Insert(post *entity.Post) error
}

type PostQuery struct {
	Start   uint64
	Results uint64
	Year    int
	Month   time.Month
	UrlName string
}

type postDao struct {
	conn   *dbr.Connection
	config *config.Config
}

func NewPostDao(conn *dbr.Connection, conf *config.Config) PostDao {
	return &postDao{
		conn:   conn,
		config: conf,
	}
}

func (pd *postDao) SelectOne(id int64) (*entity.Post, error) {
	sess := pd.conn.NewSession(nil)
	var ps []*entity.Post
	sb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where("id = ?", id)
	cnt, err := sb.Load(&ps)
	if err != nil {
		return nil, err
	}
	if cnt <= 0 {
		return nil, nil
	}
	p := ps[0]
	var nps []*entity.Post
	nsb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where("id = ?", dbr.Select("min(id)").From("post").Where("id > ?", id))
	nCnt, err := nsb.Load(&nps)
	if err != nil {
		return nil, err
	}
	if nCnt > 0 {
		p.NextPost = nps[0]
	}
	var pps []*entity.Post
	psb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where("id = ?", dbr.Select("max(id)").From("post").Where("id < ?", id))
	pCnt, err := psb.Load(&pps)
	if err != nil {
		return nil, err
	}
	if pCnt > 0 {
		p.PreviousPost = pps[0]
	}
	return p, nil
}

func (pd *postDao) SelectByQuery(query *PostQuery) ([]*entity.Post, error) {
	sess := pd.conn.NewSession(nil)
	sb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").From("post")
	if query.Year != 0 && query.Month != 0 {
		loc, err := time.LoadLocation(pd.config.Locale)
		if err != nil {
			return nil, err
		}
		startDateTime := time.Date(query.Year, query.Month, 1, 0, 0, 0, 0, loc)
		endDateTime := startDateTime.AddDate(0, 1, 0)
		sb = sb.Where("created_at >= ? and created_at < ?", startDateTime, endDateTime)
	}
	if len(query.UrlName) > 0 {
		sb = sb.Where("url_name = ?", query.UrlName)
	}
	sb = sb.OrderBy("id desc").Limit(query.Results).Offset(query.Start - 1)
	var ps []*entity.Post
	if _, err := sb.Load(&ps); err != nil {
		return nil, err
	}
	return ps, nil
}

func (pd *postDao) Insert(post *entity.Post) error {
	sess := pd.conn.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		return err
	}
	post.Id, err = pd.issuePostId()
	if err != nil {
		return pd.rollback(tx, err)
	}
	now := time.Now()
	if post.CreatedAt == nil {
		post.CreatedAt = &now
	}
	if post.UpdatedAt == nil {
		post.UpdatedAt = &now
	}
	ib := sess.InsertInto("post").
		Columns("id", "created_at", "updated_at", "url_name", "title", "body").
		Values(post.Id, post.CreatedAt, post.UpdatedAt, post.UrlName, post.Title, post.Body)
	if _, err := ib.Exec(); err != nil {
		return pd.rollback(tx, err)
	}
	return nil
}

func (pd *postDao) rollback(tx *dbr.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		panic(err)
	}
	return err
}

func (pd *postDao) issuePostId() (int64, error) {
	sess := pd.conn.NewSession(nil)
	ub := sess.Update("last_id").Set("post_last_id", dbr.Expr("last_insert_id(post_last_id + 1)"))
	if _, err := ub.Exec(); err != nil {
		return 0, err
	}
	var ids []int64
	sb := sess.Select("last_insert_id()").From("dual")
	if _, err := sb.Load(&ids); err != nil {
		return 0, err
	}
	if len(ids) <= 0 {
		return 0, errors.New("failed to issue a post id")
	}
	return ids[0], nil
}
