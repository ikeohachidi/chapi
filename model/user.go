package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (u *User) Create(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`
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

	row := stmt.QueryRowx(u.Email)

	err = row.Scan(u.ID)
	if err != nil {
		return
	}

	return
}

func (u *User) Delete(db *sqlx.DB) (err error) {
	_, err = db.Exec("DELETE FROM user WHERE id=$1", u.ID)

	return
}
