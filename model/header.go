package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Header struct {
	ID      uint   `json:"id,omitempty" db:"id"`
	RouteID uint   `json:"routeId" db:"route_id"`
	UserID  uint   `json:"userId,omitempty" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" value:"value"`
}

func (h *Header) Create(db sqlx.DB) (err error) {
	stmt := fmt.Sprintf(`
		INSERT INTO	header(user_id, route_id, name, value)
		VALUES ($1, $2, pgp_sym_encrypt($3, '%[1]v'), pgp_sym_encrypt($4, '%[1]v'))
		RETURNING id
	`, PG_CRYPT_KEY)

	row := db.QueryRow(stmt, h.UserID, h.RouteID, h.Name, h.Value)

	err = row.Scan(h.ID)
	if err != nil {
		return
	}

	return nil
}

func (h *Header) Update(db sqlx.DB) (err error) {
	stmt := fmt.Sprintf(`
		UPDATE header
		SET name = pgp_sym_encrypt($1, '%[1]v'), value = pgp_sym_encrypt($2, '%[1]v')
		WHERE id= $3 AND user_id = $4 AND route_id = $5
	`, PG_CRYPT_KEY)

	_, err = db.Exec(stmt, h.Name, h.Value, h.ID, h.UserID, h.RouteID)
	if err != nil {
		return
	}

	return nil
}

func (h *Header) Delete(db sqlx.DB) (err error) {
	stmt, err := db.Prepare(`DELETE FROM header WHERE id = $1 AND user_id = $2 AND route_id = $3`)

	if err != nil {
		return
	}

	_, err = stmt.Exec(h.ID, h.UserID, h.RouteID)
	if err != nil {
		return
	}

	return nil
}

func (h *Header) FetchAll(db sqlx.DB) ([]Header, error) {
	headers := []Header{}

	stmt := fmt.Sprintf(`
		SELECT id, pgp_sym_decrypt(name::bytea, '%[1]v') as name, pgp_sym_decrypt(value::bytea, '%[1]v') as value
		FROM header
		WHERE user_id = $1 AND route_id = $2
	`, PG_CRYPT_KEY)

	rows, err := db.Query(stmt, h.UserID, h.RouteID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var header Header

		err := rows.Scan(&h.ID, &h.Name, &h.Value)
		if err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}
