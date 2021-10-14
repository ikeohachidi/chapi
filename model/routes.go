package model

import (
	"time"
)

type Route struct {
	ID          uint      `json:"id" db:"id"`
	ProjectID   uint      `json:"projectId" db:"project_id"`
	UserID      uint      `json:"userId,omitempty" db:"user_id"`
	Method      string    `json:"method" db:"method"`
	Path        string    `json:"path" db:"path"`
	Destination string    `json:"destination" db:"destination"`
	Description string    `json:"description" db:"description"`
	Body        string    `json:"body" db:"body"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// SaveRoute will either create a new Route or update and existing one
func (c *Conn) SaveRoute(route *Route) (err error) {
	queryStmt := `
		INSERT INTO route (project_id, user_id, method, path, destination, body, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	stmt, err := c.db.Preparex(queryStmt)
	if err != nil {
		return
	}

	row := stmt.QueryRow(route.ProjectID, route.UserID, route.Method, route.Path, route.Destination, route.Body, route.Description)

	err = row.Scan(&route.ID)

	return
}

func (c *Conn) UpdateRoute(route Route) (err error) {
	queryStmt := `
		UPDATE route
		SET method = $1, path = $2, destination = $3, body = $4, description = $5
		WHERE id = $6 AND user_id = $7
	`

	_, err = c.db.Exec(queryStmt, route.Method, route.Path, route.Destination, route.Body, route.Description, route.ID, route.UserID)
	if err != nil {
		return
	}

	return
}

func (c *Conn) GetRoutesByProjectId(projectID uint, userID uint) (routes []Route, err error) {
	query := `
		SELECT * FROM route
		WHERE project_id = $1 AND user_id = $2
	`

	err = c.db.Select(&routes, query, projectID, userID)

	if err != nil {
		return
	}

	return
}

func (c *Conn) DeleteRoute(routeID uint, userID uint) (err error) {
	_, err = c.db.Exec("DELETE FROM routes WHERE id=$1 AND user_id=$2", routeID, userID)

	return
}
