package model

import (
	"database/sql"
	"encoding/json"
)

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

type CarList []Car

type YearsList []string

func (c *Car) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO car_list(model, category, color, price, mini_desc, description, images, second_images, manufacturer, year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		c.Model, c.Category.Id, c.Color.Id, c.Price, c.MiniDesc, c.Description, c.Images, getSecondImagesJsonString(c.SecondImages), c.Manufacturer.Id, c.Year); err != nil {
		return err
	}
	return nil
}

func (c *Car) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM car_list WHERE id = $1", c.Id); err != nil {
		return err
	}
	return nil
}

func (c *Car) Get(db *sql.DB, id int) error {
	row := db.QueryRow("SELECT car_list.id, car_list.model, car_list.category, category.name, car_list.color, colors.name, car_list.price, car_list.mini_desc, car_list.description, car_list.images, car_list.second_images, car_list.manufacturer, manufacturer.name, car_list.year FROM car_list, manufacturer, category, colors WHERE car_list.id = $1 AND colors.id = car_list.color AND category.id = car_list.category AND manufacturer.id = car_list.manufacturer", id)

	var secondImages string
	if err := row.Scan(&c.Id, &c.Model, &c.Category.Id, &c.Category.Name, &c.Color.Id, &c.Color.Name, &c.Price, &c.MiniDesc, &c.Description, &c.Images, &secondImages, &c.Manufacturer.Id, &c.Manufacturer.Name, &c.Year); err != nil {
		return err
	}
	c.SecondImages = getSecondImagesSliceOfString(secondImages)
	return nil
}

func (c *Car) Edit(db *sql.DB) error {
	if _, err := db.Exec("UPDATE car_list SET model = $1, manufacturer = $2, color = $3, category = $4, price = $5, year = $6, images = $7, second_images = $8, description = $9, mini_desc = $10 WHERE id = $11;",
		c.Model, c.Manufacturer.Id, c.Color.Id, c.Category.Id, c.Price, c.Year, c.Images, getSecondImagesJsonString(c.SecondImages), c.Description, c.MiniDesc, c.Id); err != nil {
		return err
	}
	return nil
}

func (yl *YearsList) Get(db *sql.DB) error {
	rows, err := db.Query("SELECT car_list.year FROM car_list GROUP BY car_list.year")
	if err != nil {
		return err
	}
	for rows.Next() {
		var year string
		if err := rows.Scan(&year); err != nil {
			return err
		}
		*yl = append(*yl, year)
	}
	return nil
}

func (cl *CarList) GetAll(db *sql.DB) error {
	rows, err := db.Query(
		"SELECT car_list.id, car_list.model, car_list.category, category.name, car_list.color, colors.name, car_list.price, car_list.mini_desc, car_list.description, car_list.images, car_list.second_images, car_list.manufacturer, manufacturer.name, car_list.year FROM car_list, manufacturer, category, colors WHERE colors.id = car_list.color AND category.id = car_list.category AND manufacturer.id = car_list.manufacturer ORDER BY car_list.id DESC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var c Car
		var secondImages string
		if err := rows.Scan(&c.Id, &c.Model, &c.Category.Id, &c.Category.Name, &c.Color.Id, &c.Color.Name, &c.Price, &c.MiniDesc, &c.Description, &c.Images, &secondImages, &c.Manufacturer.Id, &c.Manufacturer.Name, &c.Year); err != nil {
			return err
		}
		c.SecondImages = getSecondImagesSliceOfString(secondImages)
		*cl = append(*cl, c)
	}
	return nil
}

func getSecondImagesJsonString(images []string) string {
	result, _ := json.Marshal(images)
	return string(result)
}

func getSecondImagesSliceOfString(images string) []string {
	var result []string
	_ = json.Unmarshal([]byte(images), &result)
	return result
}
