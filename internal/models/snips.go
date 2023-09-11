package models

import (
	"database/sql"
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
	return nil, nil
}

func (m *SnipModel) Latest() ([]*Snip, error) {
	return nil, nil
}
