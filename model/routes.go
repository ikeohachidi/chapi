package model

import "time"

type Route struct {
	ID          uint      `json:"id" db:"id"`
	ProjectID   uint      `json:"projectId" db:"project_id"`
	UserID      uint      `json:"userId" db:"user_id"`
	Type        string    `json:"type" db:"type"`
	Path        string    `json:"path" db:"path"`
	Destination string    `json:"destination" db:"destination"`
	Body        string    `json:"body" db:"body"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Queries     []Query   `json:"queries" db:"queries"`
}

// SaveRoute will either create a new Route or update and existing one
func (c *Conn) SaveRoute(route Route) (routeID uint, err error) {
	queryStmt := `
		INSERT INTO route (project_id, user_id, type, path, destination, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	if route.ID > 0 {
		queryStmt = `
			INSERT INTO route (id, project_id, user_id, type, path, destination, body)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id)
			DO UPDATE SET 
			type = EXCLUDED.type, path = EXCLUDED.path, destination = EXCLUDED.destination, body = EXCLUDED.body
			RETURNING id
		`
	}

	stmt, err := c.db.Preparex(queryStmt)
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
