package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"
)

type Route struct {
	ID          uint           `json:"id" db:"id"`
	ProjectID   uint           `json:"projectId" db:"project_id"`
	UserID      uint           `json:"userId" db:"user_id"`
	Type        string         `json:"type" db:"type"`
	Path        string         `json:"path" db:"path"`
	Destination string         `json:"destination" db:"destination"`
	Body        sql.NullString `json:"body" db:"body"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	Queries     Queries        `json:"queries" db:"queries"`
}

type Queries []Query

func (qu *Queries) Scan(src interface{}) (err error) {
	buf := bytes.NewBuffer(src.([]uint8))

	trimmed := bytes.TrimPrefix(buf.Bytes(), []byte("{\""))
	trimmed = bytes.TrimSuffix(trimmed, []byte("\"}"))

	queries := bytes.Split(trimmed, []byte("\",\""))

	for _, query := range queries {
		var q Query

		cleanedJSON := bytes.ReplaceAll(query, []byte("\\"), []byte(""))

		err = json.Unmarshal(cleanedJSON, &q)
		if err != nil {
			return
		}

		*qu = append(*qu, q)
	}

	return nil
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

func (c *Conn) GetRoutesByProjectId(projectID uint, userID uint) (routes []Route, err error) {
	query := `SELECT * FROM route WHERE project_id=$1 AND user_id=$2`

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

func (c *Conn) GetRouteFromNameAndPath(name string, path string) (route Route, err error) {
	queryStmt := `
		SELECT
			route.*,
			array_agg(json_build_object('id', "query".id, 'name', "query"."name", 'value', "query"."value")) as queries
		FROM route
		INNER JOIN "query" ON route.id = "query".route_id
		WHERE route.project_id = (
				SELECT id FROM project WHERE "name" = $1
			)
		AND path = $2
		GROUP BY route.id
	`

	row := c.db.QueryRowx(queryStmt, name, path)

	err = row.StructScan(&route)
	if err != nil {
		return
	}

	return
}
