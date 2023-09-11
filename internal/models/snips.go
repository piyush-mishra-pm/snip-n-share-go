package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snip struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnipModel struct {
	DB *sql.DB
}

func (m *SnipModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snips (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnipModel) Get(id int) (*Snip, error) {
	stmt := `SELECT id, title, content, created, expires FROM snips
WHERE expires > UTC_TIMESTAMP() AND id = ?`
	s := &Snip{}

	if err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnipModel) Latest() ([]*Snip, error) {
	return nil, nil
}
