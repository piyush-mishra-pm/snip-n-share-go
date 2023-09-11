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
	return 0, nil
}

func (m *SnipModel) Get(id int) (*Snip, error) {
	return nil, nil
}

func (m *SnipModel) Latest() ([]*Snip, error) {
	return nil, nil
}
