package model

type Query struct {
	ID      uint   `json:"id" db:"id"`
	RouteID uint   `json:"routeId" db:"route_id"`
	UserID  uint   `json:"userId" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" db:"value"`
}

// SaveQuery will either create a new query or update(really upsert) an existing one
func (c *Conn) SaveQuery(query *Query) (err error) {
	queryStmt := `
		INSERT INTO "query" (route_id, user_id, "name", "value")
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
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
	queryStmt := `
		UPDATE "query" 
		SET "name" = $1, "value" = $2
		WHERE id = $3 AND route_id = $4 AND user_id = $5
	`

	_, err = c.db.Exec(queryStmt, query.Name, query.Value, query.ID, query.RouteID, query.UserID)

	return
}

func (c *Conn) GetRouteQueries(routeID uint, userID uint) (queries []Query, err error) {
	stmt, err := c.db.Preparex(`
		SELECT id, route_id, name, value FROM query WHERE route_id = $1 AND user_id = $2
	`)

	if err != nil {
		return
	}

	err = stmt.Select(&queries, routeID, userID)
	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteQuery(query Query) (err error) {
	_, err = c.db.Exec(
		"DELETE FROM query WHERE id = $1 AND route_id = $2 AND user_id = $3",
		query.ID,
		query.RouteID,
		query.UserID,
	)

	return
}
