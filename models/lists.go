// License: AGPL-3.0-only
// (c) 2024 Dakota Walsh <kota@nilsu.org>
package models

import (
	"database/sql"
	"errors"
)

type List struct {
	Name string
	Body string
}

type ListModel struct {
	DB *sql.DB
}

var ErrNoRecord = errors.New("models: no matching record found")

// This will return a specific list based on its name.
func (m *ListModel) Get(name string) (*List, error) {
	stmt := "SELECT name, body FROM lists WHERE name = ?"
	row := m.DB.QueryRow(stmt, name)
	l := &List{}
	err := row.Scan(&l.Name, &l.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return l, nil
}

// This will insert a new list into the database.
func (m *ListModel) Insert(name string, body string) error {
	stmt := "INSERT INTO lists (name, body) VALUES(?, ?)"
	_, err := m.DB.Exec(stmt, name, body)
	if err != nil {
		return err
	}
	return nil
}

// This will update a list in the database.
func (m *ListModel) Update(name string, body string) error {
	stmt := "UPDATE lists SET body = ? WHERE name = ?"
	_, err := m.DB.Exec(stmt, body, name)
	if err != nil {
		return err
	}
	return nil
}
