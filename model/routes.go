package model

import "time"

type Route struct {
	Id          uint      `json:"id" db:"id"`
	ProjectId   uint      `json:"project_id" db:"project_id"`
	Type        string    `json:"type" db:"type"`
	Path        string    `json:"path" db:"path"`
	Destination string    `json:"destination" db:"destination"`
	Body        string    `json:"body" db:"body"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Queries     []Query   `json:"queries" db:"queries"`
}

func (c *Conn) SaveRoute(route Route) (routeId uint, err error) {
	stmt, err := c.db.Preparex(`
		INSERT INTO routes (id, project_id, type, path, destination, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO UPDATE SET 
		type = EXCLUDED.type, path = EXCLUDED.path, destination = EXCLUDED.destination, body = EXCLUDED.body
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

func (c *Conn) GetRoutesByProjectId(projectId uint) (routes []Route, err error) {
	query := `SELECT * FROM route WHERE project_id=$1`

	err = c.db.Select(&routes, query, projectId)

	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteRoute(routeId uint, userId uint) (err error) {
	_, err = c.db.Exec("DELETE FROM routes WHERE id=$1 AND user_id=$2", routeId, userId)

	return
}
