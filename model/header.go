package model

type Header struct {
	ID    uint   `json:"id,omitempty" db:"id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" value:"value"`
}

func (c *Conn) SaveHeader(header *Header, userID, routeID uint) (err error) {
	stmt := `
		INSERT INTO	header(user_id, route_id, name, value)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := c.db.QueryRow(stmt, userID, routeID, header.Name, header.Value)

	err = row.Scan(&header.ID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) UpdateHeader(header Header, userID, routeID uint) (err error) {
	stmt := `
		UPDATE header
		SET name = $1, value = $2
		WHERE id= $3 AND user_id = $4 AND route_id = $5
	`

	_, err = c.db.Exec(stmt, header.Name, header.Value, header.ID, userID, routeID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) DeleteHeader(headerName string, userID, routeID uint) (err error) {
	stmt, err := c.db.Prepare(`
		DELETE FROM header
		WHERE "name" = $1 AND user_id = $2 AND route_id = $3
	`)

	if err != nil {
		return
	}

	_, err = stmt.Exec(headerName, userID, routeID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) GetHeader(userID, routeID uint) (headers []Header, err error) {
	stmt := `
		SELECT id, name, value
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

		err := rows.Scan(&header.ID, &header.Name, &header.Value)
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
