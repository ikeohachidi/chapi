package model

import (
	"fmt"
)

type Query struct {
	ID      uint   `json:"id" db:"id"`
	RouteID uint   `json:"routeId" db:"route_id"`
	UserID  uint   `json:"userId,omitempty" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" db:"value"`
}

// SaveQuery will either create a new query or update(really upsert) an existing one
func (c *Conn) SaveQuery(query *Query) (err error) {
	queryStmt := fmt.Sprintf(`
		INSERT INTO "query" (route_id, user_id, "name", "value")
		VALUES ($1, $2, pgp_sym_encrypt($3, '%[1]v'), pgp_sym_encrypt($4, '%[1]v'))
		RETURNING id
	`, PG_CRYPT_KEY)

	stmt, err := c.db.Preparex(queryStmt)

	if err != nil {
		return
	}

	row := stmt.QueryRowx(query.RouteID, query.UserID, query.Name, query.Value)

	if err != nil {
		return
	}

	err = row.Scan(&query.ID)

	return
}

func (c *Conn) UpdateQuery(query Query) (err error) {
	queryStmt := fmt.Sprintf(`
		UPDATE "query" 
		SET "name" = pgp_sym_encrypt($1, '%[1]v'), "value" = pgp_sym_encrypt($2, '%[1]v')
		WHERE id = $3 AND route_id = $4 AND user_id = $5
	`, PG_CRYPT_KEY)

	_, err = c.db.Exec(queryStmt, query.Name, query.Value, query.ID, query.RouteID, query.UserID)

	return
}

func (c *Conn) GetRouteQueries(routeID uint, userID uint) ([]Query, error) {
	queries := []Query{}

	stmt, err := c.db.Preparex(
		fmt.Sprintf(`
			SELECT id, route_id, pgp_sym_decrypt(name::bytea, '%[1]v') as name, pgp_sym_decrypt(value::bytea, '%[1]v') as value FROM query WHERE route_id = $1 AND user_id = $2
		`, PG_CRYPT_KEY),
	)

	if err != nil {
		return nil, err
	}

	err = stmt.Select(&queries, routeID, userID)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

func (c *Conn) DeleteQuery(query Query) (err error) {
	stmt, err := c.db.Preparex("DELETE FROM query WHERE id = $1 AND route_id = $2 AND user_id = $3")
	if err != nil {
		return
	}

	_, err = stmt.Exec(query.ID, query.RouteID, query.UserID)
	if err != nil {
		return
	}

	return
}
