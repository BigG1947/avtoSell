package model

import (
	"database/sql"
)

type News struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	MiniDesc    string `json:"mini_desc"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type NewsList []News

func (n *News) Add(db *sql.DB) error {
	if _, err := db.Exec(
		"INSERT INTO news_list(title, miniDesk, Description, Images) VALUES ($1, $2, $3, $4)",
		n.Title, n.MiniDesc, n.Description, n.Image); err != nil {
		return err
	}
	return nil
}

func (n *News) Get(db *sql.DB, id int) error {
	if err := db.QueryRow("SELECT id, title, miniDesk, Description, Images FROM news_list WHERE id = $1", id).Scan(&n.Id, &n.Title, &n.MiniDesc, &n.Description, &n.Image); err != nil {
		return err
	}
	return nil
}

func (n *News) Edit(db *sql.DB) error {
	if _, err := db.Exec("UPDATE news_list SET title = $1, miniDesk = $2, Description = $3, Images = $4 WHERE id = $5", n.Title, n.MiniDesc, n.Description, n.Image, n.Id); err != nil {
		return err
	}
	return nil
}

func (n *News) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM news_list WHERE id = $1", n.Id); err != nil {
		return err
	}
	return nil
}

func (nl *NewsList) GetAll(db *sql.DB) error {
	rows, err := db.Query("SELECT id, title, miniDesk, Description, Images FROM news_list ORDER BY id DESC")
	if err != nil {
		return err
	}

	for rows.Next() {
		var n News
		if err := rows.Scan(&n.Id, &n.Title, &n.MiniDesc, &n.Description, &n.Image); err != nil {
			return err
		}
		*nl = append(*nl, n)
	}
	return nil
}

func (nl *NewsList) GetLatest(db *sql.DB) error {
	rows, err := db.Query("SELECT id, title, miniDesk, Description, Images FROM news_list ORDER BY id DESC LIMIT 3")
	if err != nil {
		return err
	}

	for rows.Next() {
		var n News
		if err := rows.Scan(&n.Id, &n.Title, &n.MiniDesc, &n.Description, &n.Image); err != nil {
			return err
		}
		*nl = append(*nl, n)
	}
	return nil
}
