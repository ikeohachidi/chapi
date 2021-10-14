package model

import (
	"bytes"
	"encoding/json"
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
	stmt, err := c.db.Preparex(`
		SELECT
			route.*,
			array_to_json(array_distinct(array_agg("query"))) as queries,
			array_to_json(array_distinct(array_agg("header"))) as headers
		FROM route
		INNER JOIN "query" ON "query".route_id = route.id
		INNER JOIN "header" ON "header".route_id = route.id
		WHERE route.project_id = (
				SELECT id FROM project WHERE "name" = $1 
			)
		AND path = $2 
		GROUP BY route.id;
	`)

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
