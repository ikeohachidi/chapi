package model

type Header struct {
	ID      uint   `json:"id,omitempty" db:"id"`
	UserID  uint   `json:"user_id,omitempty" db:"user_id"`
	RouteID uint   `json:"route_id,omitempty" db:"route_id"`
	Name    string `json:"name" db:"name"`
	Value   string `json:"value" value:"value"`
}

func (c *Conn) SaveHeader(header Header, userID, routeID uint) (err error) {
	stmt := `
		INSERT INTO	header(user_id, route_id, name, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name)
		DO UPDATE
		SET value = EXCLUDED.value
	`

	_, err = c.db.Exec(stmt, header.UserID, header.RouteID, header.Name, header.Value)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) DeleteHeader(headerName string, userID, routeID uint) (err error) {
	stmt, err := c.db.Prepare(`
		DELETE FROM header
		WHERE name = $1 AND user_id = $2 AND route_id = $3
	`)

	if err != nil {
		return
	}

	_, err = stmt.Exec(stmt, headerName, userID, routeID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) GetHeader(userID, routeID uint) (headers []Header, err error) {
	stmt := `
		SELECT name, value
		FROM header
		WHERE user_id = $1 AND route_id = $2
	`

	rows, err := c.db.Query(stmt, userID, routeID)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var header Header

		err := rows.Scan(&header.Name, &header.Value)
		if err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
