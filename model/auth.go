package model

import "time"

type User struct {
	Id        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (c *Conn) CreateUser(user User) (userId uint, err error) {
	stmt, err := c.db.Preparex(`
		INSERT INTO user (id, email)
		VALUES ($1, $2)
	`)

	if err != nil {
		return
	}

	row, err := stmt.Queryx(user)

	if err != nil {
		return
	}

	row.Scan(&userId)

	row.Close()

	return
}

func (c *Conn) DeleteUser(userId uint) (err error) {
	_, err = c.db.Exec("DELETE FROM user WHERE id=$1", userId)

	return
}
