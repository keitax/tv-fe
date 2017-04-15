package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/keitax/textvid/entity"
)

type PostDao interface {
	SelectOne(id int64) (*entity.Post, error)
	SelectByQuery(query *PostQuery) ([]*entity.Post, error)
	Insert(post *entity.Post) error
}

type PostQuery struct {
	Start   int
	Results int
}

type postDaoImpl struct {
	db *sql.DB
}

func NewPostDao(db *sql.DB) PostDao {
	return &postDaoImpl{db}
}

func (pd *postDaoImpl) SelectOne(id int64) (*entity.Post, error) {
	rows, err := pd.db.Query("SELECT ID, CREATED_AT, UPDATED_AT, URL_NAME, TITLE, BODY FROM POST WHERE ID = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil // not found
	}
	p := new(entity.Post)
	if err := rows.Scan(&p.Id, &p.CreatedAt, &p.UpdatedAt, &p.UrlName, &p.Title, &p.Body); err != nil {
		return nil, err
	}
	return p, nil
}

func (pd *postDaoImpl) SelectByQuery(query *PostQuery) ([]*entity.Post, error) {
	rows, err := pd.db.Query("SELECT ID, CREATED_AT, UPDATED_AT, URL_NAME, TITLE, BODY FROM POST")
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

func (pd *postDaoImpl) Insert(post *entity.Post) error {
	tx, err := pd.db.Begin()
	if err != nil {
		return err
	}
	post.Id, err = pd.issuePostId()
	if err != nil {
		return pd.rollback(tx, err)
	}
	now := time.Now()
	_, err = pd.db.Exec("INSERT INTO POST (ID, CREATED_AT, UPDATED_AT, URL_NAME, TITLE, BODY) VALUES (?, ?, ?, ?, ?, ?)", post.Id, now, now, post.UrlName, post.Title, post.Body)
	if err != nil {
		return pd.rollback(tx, err)
	}
	return nil
}

func (pd *postDaoImpl) rollback(tx *sql.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		panic(err)
	}
	return err
}

func (pd *postDaoImpl) issuePostId() (int64, error) {
	if _, err := pd.db.Exec("UPDATE LAST_ID SET POST_LAST_ID = LAST_INSERT_ID(POST_LAST_ID + 1)"); err != nil {
		return 0, err
	}
	rs, err := pd.db.Query("SELECT LAST_INSERT_ID() FROM DUAL")
	if err != nil {
		return 0, err
	}
	defer rs.Close()
	if !rs.Next() {
		return 0, errors.New("Failed to issue the id")
	}
	var postId int64
	if err := rs.Scan(&postId); err != nil {
		return 0, err
	}
	return postId, nil
}
