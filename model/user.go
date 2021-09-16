package model

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (c *Conn) CreateUser(user User) (userID uint, err error) {
	stmt, err := c.db.Preparex(`
		WITH e AS(
			INSERT INTO "user" (email) 
			VALUES ($1)
			ON CONFLICT (email) DO NOTHING
			RETURNING id
		)

		SELECT * FROM e
		UNION 
		sELECT id FROM "user" WHERE email=$1;
	`)

	if err != nil {
		return
	}

	row := stmt.QueryRowx(user.Email)

	err = row.Scan(&userID)
	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteUser(userID uint) (err error) {
	_, err = c.db.Exec("DELETE FROM user WHERE id=$1", userID)

	return
}
