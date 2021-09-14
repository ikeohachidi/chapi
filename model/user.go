package model

import "time"

type User struct {
	ID        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (c *Conn) CreateUser(user User) (userID uint, err error) {
	stmt, err := c.db.Preparex(`
		INSERT INTO "user" (email)
		VALUES ($1)
		ON CONFLICT (email)
		DO NOTHING
	`)

	if err != nil {
		return
	}

	row, err := stmt.Queryx(user.Email)

	if err != nil {
		return
	}

	row.Scan(&userID)

	row.Close()

	return
}

func (c *Conn) DeleteUser(userID uint) (err error) {
	_, err = c.db.Exec("DELETE FROM user WHERE id=$1", userID)

	return
}
