package model

import "database/sql"

type Color struct {
	Id   int
	Name string
}

type ColorList []Color

func (cl *ColorList) GetAll(db *sql.DB) error {
	rows, err := db.Query("SELECT id, name FROM colors ORDER BY id DESC;")
	if err != nil {
		return err
	}
	for rows.Next() {
		var c Color
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return err
		}
		*cl = append(*cl, c)
	}
	return nil
}

func (c *Color) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO colors(name) VALUES ($1);", c.Name); err != nil {
		return err
	}
	return nil
}

func (c *Color) Edit(db *sql.DB) error {
	if _, err := db.Exec("UPDATE colors SET name = $1 WHERE id = $2", c.Name, c.Id); err != nil {
		return err
	}
	return nil
}

func (c *Color) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM colors WHERE id = $1", c.Id); err != nil {
		return err
	}
	return nil
}

func (c *Color) Get(db *sql.DB, id int) error {
	if err := db.QueryRow("SELECT id, name FROM colors WHERE id = $1", id).Scan(&c.Id, &c.Name); err != nil {
		return err
	}
	return nil
}
