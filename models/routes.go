package models

import "time"

type Route struct {
	Id          uint      `json:"id" db:"id"`
	Path        string    `json:"path" db:"path"`
	ProjectId   uint      `json:"project_id" db:"project_id"`
	RequestBody string    `json:"request_body" db:"request_body"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (c *Conn) CreateRoute(route Route) (routeId uint, err error) {
	stmt, err := c.db.Preparex(`
		INSERT INTO routes (path, project_id, request_body, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`)
	if err != nil {
		return
	}

	row, err := stmt.Queryx(route)

	if err != nil {
		return
	}

	row.Scan(&routeId)

	row.Close()

	return
}

func (c *Conn) DeleteRoute(routeId uint) (err error) {
	_, err = c.db.Exec("DELETE FROM routes WHERE id=$1")

	if err != nil {
		return
	}

	return
}
