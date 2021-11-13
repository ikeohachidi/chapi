package model

type PermOrigin struct {
	ID      uint   `json:"id" db:"id"`
	URL     string `json:"url" db:"url"`
	RouteID uint   `json:"routeId" db:"route_id"`
}

func (c *Conn) CreatePermOrigin(permOrigin *PermOrigin) (err error) {
	stmt, err := c.db.Preparex(`INSERT INTO perm_origin(url, route_id) VALUES($1, $2) RETURNING id`)
	if err != nil {
		return
	}

	row := stmt.QueryRowx(permOrigin.URL, permOrigin.RouteID)

	err = row.Scan(&permOrigin.ID)

	return
}

func (c *Conn) GetPermOrigins(routeID uint) ([]PermOrigin, error) {
	origins := []PermOrigin{}

	stmt, err := c.db.Preparex(`SELECT * FROM perm_origin WHERE route_id = $1`)
	if err != nil {
		return nil, err
	}

	err = stmt.Select(&origins, routeID)
	if err != nil {
		return nil, err
	}

	return origins, nil
}

func (c *Conn) UpdatePermOrigin(permOrigin PermOrigin) (err error) {
	stmt, err := c.db.Preparex(`UPDATE perm_origin SET url = $1 WHERE id = $2 AND route_id = $3`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(permOrigin.URL, permOrigin.ID, permOrigin.RouteID)
	return
}

func (c *Conn) DeletePermOrigin(permOrigin PermOrigin) (err error) {
	stmt, err := c.db.Preparex(`DELETE FROM perm_origin WHERE id = $1 AND route_id = $2`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(permOrigin.ID, permOrigin.RouteID)

	return
}
