package model

type Query struct {
	Id      uint   `json:"id" db:"id"`
	RouteId uint   `json:"routeId" db:"route_id"`
	UserId  uint   `json:"userId" db:"user_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" db:"value"`
}

// SaveQuery will either create a new query or update(really upsert) an existing one
func (c *Conn) SaveQuery(query *Query) (err error) {
	queryStmt := `
		INSERT INTO "query"(route_id, user_id, name, value) RETURNING id
	`

	if query.Id > 0 {
		queryStmt = `
			INSERT INTO "query" (id, route_id, user_id, name, value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id)
			DO UPDATE SET 
			name = EXCLUDED.name, value = EXCLUDED.value,
			RETURNING id
		`
	}

	stmt, err := c.db.Preparex(queryStmt)

	if err != nil {
		return
	}

	row, err := stmt.Queryx(query)

	if err != nil {
		return
	}

	row.Scan(query.Id)

	row.Close()

	return
}

func (c *Conn) GetRouteQueries(routeId uint, userId uint) (queries []Query, err error) {
	stmt, err := c.db.Preparex(`
		SELECT id, route_id, name, value FROM query WHERE route_id = $1 AND user_id = $2
	`)

	if err != nil {
		return
	}

	err = stmt.Select(&queries, routeId, userId)
	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteQuery(queryId uint, routeId uint, userId uint) (err error) {
	_, err = c.db.Exec("DELETE FROM query WHERE id = $1 AND route_id = $2 AND user_id = $3", queryId)

	return
}
