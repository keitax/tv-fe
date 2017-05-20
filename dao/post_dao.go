package dao

import (
	"time"

	"github.com/gocraft/dbr"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/repository"
)

type PostDao interface {
	SelectOne(id int64) *entity.Post
	SelectByQuery(query *repository.PostQuery) []*entity.Post
	Insert(post *entity.Post)
	Update(post *entity.Post)
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

func (pd *postDao) SelectOne(id int64) *entity.Post {
	sess := pd.conn.NewSession(nil)
	var ps []*entity.Post
	sb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where("id = ?", id)
	cnt, err := sb.Load(&ps)
	if err != nil {
		panic(err)
	}
	if cnt <= 0 {
		return nil
	}
	p := ps[0]
	var nps []*entity.Post
	nsb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where("id = ?", dbr.Select("min(id)").From("post").Where("id > ?", id))
	nCnt, err := nsb.Load(&nps)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	if pCnt > 0 {
		p.PreviousPost = pps[0]
	}
	return p
}

func (pd *postDao) SelectByQuery(query *repository.PostQuery) []*entity.Post {
	sess := pd.conn.NewSession(nil)
	sb := sess.Select("id", "created_at", "updated_at", "url_name", "title", "body").From("post")
	if query.Year != 0 && query.Month != 0 {
		loc, err := time.LoadLocation(pd.config.Locale)
		if err != nil {
			panic(err)
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
		panic(err)
	}
	return ps
}

func (pd *postDao) Insert(post *entity.Post) {
	sess := pd.conn.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		panic(err)
	}
	post.Id = pd.issuePostId()
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
		pd.rollback(tx, err)
		panic(err)
	}
}

func (pd *postDao) Update(post *entity.Post) {
	s := pd.conn.NewSession(nil)
	tx, err := s.Begin()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	if post.UpdatedAt == nil {
		post.UpdatedAt = &now
	}
	ub := s.Update("post").SetMap(map[string]interface{}{
		"updated_at": post.UpdatedAt,
		"url_name":   post.UrlName,
		"title":      post.Title,
		"body":       post.Body,
	}).Where("id = ?", post.Id)
	if _, err := ub.Exec(); err != nil {
		pd.rollback(tx, err)
		panic(err)
	}
}

func (pd *postDao) rollback(tx *dbr.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		panic(err)
	}
	return err
}

func (pd *postDao) issuePostId() int64 {
	sess := pd.conn.NewSession(nil)
	ub := sess.Update("last_id").Set("post_last_id", dbr.Expr("last_insert_id(post_last_id + 1)"))
	if _, err := ub.Exec(); err != nil {
		panic(err)
	}
	var ids []int64
	sb := sess.Select("last_insert_id()").From("dual")
	if _, err := sb.Load(&ids); err != nil {
		panic(err)
	}
	if len(ids) <= 0 {
		panic("failed to issue a post id")
	}
	return ids[0]
}
