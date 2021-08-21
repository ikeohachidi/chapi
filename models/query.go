package models

type Query struct {
	Id    uint   `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" db:"value"`
}

func (c *Conn) SetQuery(query Query) (queryId uint, err error) {
	stmt, err := c.db.Preparex(`
		INSERT INTO query (id, name, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET 
		name = EXCLUDED.name, value = EXCLUDED.value,
		RETURNING id
	`)
	if err != nil {
		return
	}

	row, err := stmt.Queryx(query)

	if err != nil {
		return
	}

	row.Scan(&queryId)

	row.Close()

	return
}

func (c *Conn) DeleteQuery(queryId uint) (err error) {
	_, err = c.db.Exec("DELETE FROM query WHERE id=$1", queryId)

	if err != nil {
		return
	}

	return
}
