package model

type Header map[string]interface{}

func (h Header) Scan(src interface{}) (err error) {
	return nil
}

func (c *Conn) SaveHeader(header Header, userID, routeID uint) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return
	}

	for k := range header {
		_, err = tx.Exec(`
			INSERT INTO	header(user_id, route_id, name, value)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (name)
			DO UPDATE
			SET value = EXCLUDED.value
		`, userID, routeID, k, header[k])

		if err != nil {
			err = tx.Rollback()

			if err != nil {
				return
			}

			return
		}
	}

	err = tx.Commit()
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

func (c *Conn) GetHeader(userID, routeID uint) (header Header, err error) {
	stmt := `
		SELECT name, value
		FROM header
		WHERE user_id = $1 AND route_id = $2
	`

	rows, err := c.db.Query(stmt)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&header)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
