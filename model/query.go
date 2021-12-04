package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Query struct {
	ID      uint   `json:"id" db:"id"`
	RouteID uint   `json:"routeId" db:"route_id"`
	UserID  uint   `json:"userId,omitempty" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" db:"value"`
}

// Create will either create a new query or update(really upsert) an existing one
func (q *Query) Create(db *sqlx.DB) (err error) {
	queryStmt := fmt.Sprintf(`
		INSERT INTO "query" (route_id, user_id, "name", "value")
		VALUES ($1, $2, pgp_sym_encrypt($3, '%[1]v'), pgp_sym_encrypt($4, '%[1]v'))
		RETURNING id
	`, PG_CRYPT_KEY)

	stmt, err := db.Preparex(queryStmt)

	if err != nil {
		return
	}

	row := stmt.QueryRowx(q.RouteID, q.UserID, q.Name, q.Value)

	if err != nil {
		return
	}

	err = row.Scan(&q.ID)

	return
}

func (q *Query) Update(db *sqlx.DB) (err error) {
	queryStmt := fmt.Sprintf(`
		UPDATE "query" 
		SET "name" = pgp_sym_encrypt($1, '%[1]v'), "value" = pgp_sym_encrypt($2, '%[1]v')
		WHERE id = $3 AND route_id = $4 AND user_id = $5
	`, PG_CRYPT_KEY)

	_, err = db.Exec(queryStmt, q.Name, q.Value, q.ID, q.RouteID, q.UserID)

	return
}

func (q *Query) GetRouteQueries(db *sqlx.DB) ([]Query, error) {
	queries := []Query{}

	stmt, err := db.Preparex(
		fmt.Sprintf(`
			SELECT id, route_id, pgp_sym_decrypt(name::bytea, '%[1]v') as name, pgp_sym_decrypt(value::bytea, '%[1]v') as value FROM query WHERE route_id = $1 AND user_id = $2
		`, PG_CRYPT_KEY),
	)

	if err != nil {
		return nil, err
	}

	err = stmt.Select(&queries, q.RouteID, q.UserID)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

func (q *Query) Delete(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex("DELETE FROM query WHERE id = $1 AND route_id = $2 AND user_id = $3")
	if err != nil {
		return
	}

	_, err = stmt.Exec(q.ID, q.RouteID, q.UserID)
	if err != nil {
		return
	}

	return
}
