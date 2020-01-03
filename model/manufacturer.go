package model

import "database/sql"

type Manufacturer struct {
	Id   int
	Name string
}

type ManufacturerList []Manufacturer

func (ml *ManufacturerList) GetAll(db *sql.DB) error {
	rows, err := db.Query("SELECT id, name FROM manufacturer ORDER BY id DESC;")
	if err != nil {
		return err
	}
	for rows.Next() {
		var m Manufacturer
		if err := rows.Scan(&m.Id, &m.Name); err != nil {
			return err
		}
		*ml = append(*ml, m)
	}
	return nil
}

func (m *Manufacturer) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO manufacturer(name) VALUES ($1);", m.Name); err != nil {
		return err
	}
	return nil
}

func (m *Manufacturer) Edit(db *sql.DB) error {
	if _, err := db.Exec("UPDATE manufacturer SET name = $1 WHERE id = $2", m.Name, m.Id); err != nil {
		return err
	}
	return nil
}

func (m *Manufacturer) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM manufacturer WHERE id = $1", m.Id); err != nil {
		return err
	}
	return nil
}

func (m *Manufacturer) Get(db *sql.DB, id int) error {
	if err := db.QueryRow("SELECT id, name FROM manufacturer WHERE id = $1", id).Scan(&m.Id, &m.Name); err != nil {
		return err
	}
	return nil
}
