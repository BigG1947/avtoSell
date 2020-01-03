package model

import "database/sql"

type Car struct {
	Id           int
	Model        string
	MiniDesc     string
	Description  string
	Category     Category
	Color        Color
	Manufacturer Manufacturer
	Price        int
	Images       string
	SecondImages []string
	Year         int
}

func (c *Car) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO car_list(model, category, color, price, mini_desc, description, images, second_images, manufacturer, year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		c.Model, c.Category.Id, c.Color.Id, c.Price, c.MiniDesc, c.Description, c.Images, c.SecondImages, c.Manufacturer.Id, c.Year); err != nil {
		return err
	}
	return nil
}
