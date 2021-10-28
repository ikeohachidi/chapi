package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Endpoint struct {
	ID          uint      `json:"id" db:"id"`
	ProjectID   uint      `json:"projectId" db:"project_id"`
	UserID      uint      `json:"userId,omitempty" db:"user_id"`
	Method      string    `json:"method" db:"method"`
	Path        string    `json:"path" db:"path"`
	Destination string    `json:"destination" db:"destination"`
	Description string    `json:"description" db:"description"`
	Body        string    `json:"body" db:"body"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Queries     Queries   `json:"queries" db:"queries"`
	Headers     Headers   `json:"headers" db:"headers"`
}

type Headers []Header

func (he *Headers) Scan(src interface{}) (err error) {
	buf := bytes.NewBuffer(src.([]byte))

	err = json.Unmarshal(buf.Bytes(), &he)
	if err != nil {
		return
	}

	return nil
}

type Queries []Query

func (qu *Queries) Scan(src interface{}) (err error) {
	buf := bytes.NewBuffer(src.([]byte))

	err = json.Unmarshal(buf.Bytes(), &qu)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) GetRouteRequestData(projectName, routePath string) (endpoint Endpoint, err error) {
	stmt, err := c.db.Preparex(fmt.Sprintf(`
		WITH project_values(id) AS
		(
			SELECT id 
			FROM project 
			WHERE "name" = $1
		),

		route_id(id) AS 
		(
			SELECT id
			FROM route
			WHERE project_id = (
				SELECT id from project_values
			) AND "path" = $2 
		),

		query_values(id, route_id, "name", "value") AS 
		(
			SELECT id, route_id, pgp_sym_decrypt("name"::bytea, '%[1]v'), pgp_sym_decrypt("value"::bytea, '%[1]v') 
			FROM "query"
			WHERE route_id = (SELECT id from route_id)
		),

		header_values(id, route_id, "name", "value") AS 
		(
			SELECT id, route_id, pgp_sym_decrypt("name"::bytea, '%[1]v'), pgp_sym_decrypt("value"::bytea, '%[1]v') 
			FROM "header" 
			WHERE route_id = (SELECT id from route_id)
		)

		SELECT
			route.*,
			array_to_json(array_agg(distinct(header_values))) as "headers",
			array_to_json(array_agg(distinct(query_values))) as "queries"
		FROM route
		INNER JOIN header_values ON header_values.route_id = route.id
		INNER JOIN query_values ON query_values.route_id = route.id
		WHERE route.project_id = (
				SELECT id FROM project_values
			)
		AND route.id = (
			SELECT id from route_id
		) 
		GROUP BY route.id;
	`, PG_CRYPT_KEY))

	if err != nil {
		return
	}

	err = stmt.Get(&endpoint, projectName, routePath)
	if err != nil {
		return
	}
	log.Println(endpoint)

	return
}
