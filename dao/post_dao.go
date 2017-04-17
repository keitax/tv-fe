package dao

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

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
}

type postDao struct {
	db     *sql.DB
	config *config.Config
}

func NewPostDao(db *sql.DB, conf *config.Config) PostDao {
	return &postDao{
		db:     db,
		config: conf,
	}
}

func (pd *postDao) SelectOne(id int64) (*entity.Post, error) {
	sb := sq.Select("id", "created_at", "updated_at", "url_name", "title", "body").
		From("post").
		Where(sq.Eq{"ID": id})
	row := sb.RunWith(pd.db).QueryRow()
	p := new(entity.Post)
	if err := row.Scan(&p.Id, &p.CreatedAt, &p.UpdatedAt, &p.UrlName, &p.Title, &p.Body); err != nil {
		return nil, err
	}
	return p, nil
}

func (pd *postDao) SelectByQuery(query *PostQuery) ([]*entity.Post, error) {
	sb := sq.Select("id", "created_at", "updated_at", "url_name", "title", "body").From("post")
	if query.Year != 0 && query.Month != 0 {
		loc, err := time.LoadLocation(pd.config.Locale)
		if err != nil {
			return nil, err
		}
		startDateTime := time.Date(query.Year, query.Month, 1, 0, 0, 0, 0, loc)
		endDateTime := startDateTime.AddDate(0, 1, 0)
		sb = sb.Where(sq.GtOrEq{"created_at": startDateTime}).Where(sq.Lt{"created_at": endDateTime})
	}
	sb = sb.OrderBy("id desc").Limit(query.Results).Offset(query.Start - 1)

	rows, err := sb.RunWith(pd.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := []*entity.Post{}
	for rows.Next() {
		p := &entity.Post{}
		if err := rows.Scan(&p.Id, &p.CreatedAt, &p.UpdatedAt, &p.UrlName, &p.Title, &p.Body); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (pd *postDao) Insert(post *entity.Post) error {
	tx, err := pd.db.Begin()
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
	ib := sq.Insert("post").
		Columns("id", "created_at", "updated_at", "url_name", "title", "body").
		Values(post.Id, post.CreatedAt, post.UpdatedAt, post.UrlName, post.Title, post.Body)
	_, err = ib.RunWith(pd.db).Exec()
	if err != nil {
		return pd.rollback(tx, err)
	}
	return nil
}

func (pd *postDao) rollback(tx *sql.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		panic(err)
	}
	return err
}

func (pd *postDao) issuePostId() (int64, error) {
	ub := sq.Update("last_id").Set("post_last_id", sq.Expr("last_insert_id(post_last_id + 1)"))
	if _, err := ub.RunWith(pd.db).Exec(); err != nil {
		return 0, err
	}
	sb := sq.Select("last_insert_id()").From("dual")
	row := sb.RunWith(pd.db).QueryRow()
	var postId int64
	if err := row.Scan(&postId); err != nil {
		return 0, err
	}
	return postId, nil
}
