package model

import (
	"crypto/sha256"
	"database/sql"
)

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash []byte `json:"password"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	IsAdmin      bool   `json:"is_admin"`
}

type UserList []User

func GeneratePasswordHash(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}

func CheckUserPhoneExist(phone string, db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_list WHERE phone = ?;", phone).Scan(&count)
	if count > 0 {
		return false
	}
	return true
}

func CheckUserEmailExist(email string, db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_list WHERE email = ?;", email).Scan(&count)
	if count > 0 {
		return false
	}
	return true
}

func CheckUserLoginExist(login string, db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_list WHERE login = ?;", login).Scan(&count)
	if count > 0 {
		return false
	}
	return true
}

func (u *User) Registration(db *sql.DB) (bool, error) {
	_, err := db.Exec("INSERT INTO user_list(login, first_name, last_name, password, email, phone) VALUES (?,?,?,?,?,?);", u.Login, u.FirstName, u.LastName, u.PasswordHash, u.Email, u.Phone)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *User) CheckUser(db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_list WHERE email = ? AND password = ?", u.Email, u.PasswordHash).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}

func (u *User) CheckAdmin(db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_list WHERE login = ? AND password = ? AND is_admin = 1", u.Login, u.PasswordHash).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}

func (u *User) GetById(db *sql.DB, id int) error {
	row := db.QueryRow("SELECT id, email, first_name, last_name, password, phone, login, is_admin FROM user_list WHERE id = ? LIMIT 1", id)
	err := row.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Phone, &u.Login, &u.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetByEmail(db *sql.DB, email string) error {
	row := db.QueryRow("SELECT id, email, first_name, last_name, password, phone, login, is_admin FROM user_list WHERE email = ? LIMIT 1", email)
	err := row.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Phone, &u.Login, &u.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetByLogin(db *sql.DB, login string) error {
	row := db.QueryRow("SELECT id, email, first_name, last_name, password, phone, login, is_admin FROM user_list WHERE login = ? LIMIT 1", login)
	err := row.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Phone, &u.Login, &u.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (userList *UserList) GetAllUsers(db *sql.DB) error {
	rows, err := db.Query("SELECT id, email, first_name, last_name, password, phone, login, is_admin FROM user_list")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Phone, &u.Login, &u.IsAdmin)
		if err != nil {
			return err
		}
		*userList = append(*userList, u)
	}
	return nil
}

func (u *User) EditPassword(db *sql.DB) error {
	if _, err := db.Exec("UPDATE user_list SET password = $1 WHERE id = $2", u.PasswordHash, u.Id); err != nil {
		return err
	}
	return nil
}

func (u *User) EditPhone(db *sql.DB) error {
	if _, err := db.Exec("UPDATE user_list SET phone = $1 WHERE id = $2", u.Phone, u.Id); err != nil {
		return err
	}
	return nil
}

func (u *User) EditEmail(db *sql.DB) error {
	if _, err := db.Exec("UPDATE user_list SET email = $1 WHERE id = $2", u.Email, u.Id); err != nil {
		return err
	}
	return nil
}

func (u *User) EditFio(db *sql.DB) error {
	if _, err := db.Exec("UPDATE user_list SET last_name = $1, first_name = $2 WHERE id = $3", u.LastName, u.FirstName, u.Id); err != nil {
		return err
	}
	return nil
}
