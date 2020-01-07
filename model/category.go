package model

import "database/sql"

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryList []Category

func (cl *CategoryList) GetAll(db *sql.DB) error {
	rows, err := db.Query("SELECT id, name FROM category ORDER BY id DESC;")
	if err != nil {
		return err
	}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return err
		}
		*cl = append(*cl, c)
	}
	return nil
}

func (c *Category) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO category(name) VALUES ($1);", c.Name); err != nil {
		return err
	}
	return nil
}

func (c *Category) Edit(db *sql.DB) error {
	if _, err := db.Exec("UPDATE category SET name = $1 WHERE id = $2", c.Name, c.Id); err != nil {
		return err
	}
	return nil
}

func (c *Category) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM category WHERE id = $1", c.Id); err != nil {
		return err
	}
	return nil
}

func (c *Category) Get(db *sql.DB, id int) error {
	if err := db.QueryRow("SELECT id, name FROM category WHERE id = $1", id).Scan(&c.Id, &c.Name); err != nil {
		return err
	}
	return nil
}
