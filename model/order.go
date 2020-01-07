package model

import "database/sql"

type Order struct {
	Id     int
	User   User
	Car    Car
	Date   string
	Status int
}

type OrderList []Order

func (o *Order) Add(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO order_list(user_id, car_id, date, status) VALUES ($1, $2, $3, $4)", o.User.Id, o.Car.Id, o.Date, o.Status); err != nil {
		return err
	}
	return nil
}

func (ol *OrderList) GetUserNewOrders(db *sql.DB, id int) error {
	var rows *sql.Rows
	var err error
	if id != 0 {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE user_id = $1 AND status = 1 ORDER BY date, id DESC", id)
	} else {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE status = 1 ORDER BY date, id DESC", id)
	}
	if err != nil {
		return err
	}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.Id, &o.User.Id, &o.Car.Id, &o.Date, &o.Status); err != nil {
			return err
		}
		o.Car.Get(db, o.Car.Id)
		o.User.GetById(db, o.User.Id)
		*ol = append(*ol, o)
	}
	return nil
}

func (ol *OrderList) GetUserCancelOrders(db *sql.DB, id int) error {
	var rows *sql.Rows
	var err error
	if id != 0 {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE user_id = $1 AND status = 3 ORDER BY date, id DESC", id)
	} else {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE status = 3 ORDER BY date, id DESC", id)
	}
	if err != nil {
		return err
	}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.Id, &o.User.Id, &o.Car.Id, &o.Date, &o.Status); err != nil {
			return err
		}
		o.Car.Get(db, o.Car.Id)
		o.User.GetById(db, o.User.Id)
		*ol = append(*ol, o)
	}
	return nil
}

func (ol *OrderList) GetUserCheckOrders(db *sql.DB, id int) error {
	var rows *sql.Rows
	var err error
	if id != 0 {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE user_id = $1 AND status = 2 ORDER BY date, id DESC", id)
	} else {
		rows, err = db.Query("SELECT id, user_id, car_id, date, status FROM order_list WHERE status = 2 ORDER BY date, id DESC", id)
	}
	if err != nil {
		return err
	}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.Id, &o.User.Id, &o.Car.Id, &o.Date, &o.Status); err != nil {
			return err
		}
		o.Car.Get(db, o.Car.Id)
		o.User.GetById(db, o.User.Id)
		*ol = append(*ol, o)
	}
	return nil
}

func (o *Order) Cancel(db *sql.DB, userId int) error {
	if userId == 0 {
		if _, err := db.Exec("UPDATE order_list SET status = 3 WHERE id = $1", o.Id); err != nil {
			return err
		}
		return nil
	} else {
		if _, err := db.Exec("UPDATE order_list SET status = 3 WHERE id = $1 AND user_id = $2", o.Id, userId); err != nil {
			return err
		}
		return nil
	}
}

func (o *Order) Check(db *sql.DB, userId int) error {
	if userId != 0 {
		if _, err := db.Exec("UPDATE order_list SET status = 2 WHERE id = $1 AND user_id = $2", o.Id, userId); err != nil {
			return err
		}
		return nil
	} else {
		if _, err := db.Exec("UPDATE order_list SET status = 2 WHERE id = $1", o.Id); err != nil {
			return err
		}
		return nil
	}
}
