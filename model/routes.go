package model

import "time"

type Route struct {
	ID          uint      `json:"id" db:"id"`
	ProjectID   uint      `json:"project_id" db:"project_id"`
	Type        string    `json:"type" db:"type"`
	Path        string    `json:"path" db:"path"`
	Destination string    `json:"destination" db:"destination"`
	Body        string    `json:"body" db:"body"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Queries     []Query   `json:"queries" db:"queries"`
}

// TODO: query should be conditional on id existing or not
func (c *Conn) SaveRoute(route Route) (routeID uint, err error) {
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

	row.Scan(&routeID)

	row.Close()

	return
}

func (c *Conn) GetRoutesByProjectId(projectID uint) (routes []Route, err error) {
	query := `SELECT * FROM route WHERE project_id=$1`

	err = c.db.Select(&routes, query, projectID)

	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteRoute(routeID uint, userID uint) (err error) {
	_, err = c.db.Exec("DELETE FROM routes WHERE id=$1 AND user_id=$2", routeID, userID)

	return
}
