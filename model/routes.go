package model

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

// Create will either create a new Route or update and existing one
func (r *Route) Create(db *sqlx.DB) (err error) {

	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return
	}
	defer tx.Rollback()

	var routeInsertId int
	if err = tx.QueryRowx(`
		INSERT INTO route (project_id, user_id, method, path, destination, body, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, r.ProjectID, r.UserID, r.Method, strings.ToLower(r.Path), r.Destination, r.Body, r.Description).Scan(&routeInsertId); err != nil {
		return
	}

	_, err = tx.Exec(`INSERT INTO request_config(route_id) VALUES($1)`, routeInsertId)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	r.ID = uint(routeInsertId)

	return
}

func (r *Route) Update(db *sqlx.DB) (err error) {
	queryStmt := `
		UPDATE route
		SET method = $1, path = $2, destination = $3, body = $4, description = $5
		WHERE id = $6 AND user_id = $7
	`

	_, err = db.Exec(queryStmt, r.Method, strings.ToLower(r.Path), r.Destination, r.Body, r.Description, r.ID, r.UserID)
	if err != nil {
		return
	}

	return
}

func (r *Route) GetRoutesByProjectId(db *sqlx.DB) ([]Route, error) {
	routes := []Route{}

	query := `
		SELECT * FROM route
		WHERE project_id = $1 AND user_id = $2
	`

	err := db.Select(&routes, query, r.ProjectID, r.UserID)

	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *Route) Delete(db *sqlx.DB) (err error) {
	_, err = db.Exec("DELETE FROM route WHERE id=$1 AND user_id=$2", r.ID, r.UserID)

	return
}
