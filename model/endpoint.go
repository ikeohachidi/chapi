package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Endpoint struct {
	ID            uint          `json:"id" db:"id"`
	ProjectID     uint          `json:"projectId" db:"project_id"`
	UserID        uint          `json:"userId,omitempty" db:"user_id"`
	Method        string        `json:"method" db:"method"`
	Path          string        `json:"path" db:"path"`
	Destination   string        `json:"destination" db:"destination"`
	Description   string        `json:"description" db:"description"`
	Body          string        `json:"body" db:"body"`
	CreatedAt     time.Time     `json:"createdAt" db:"created_at"`
	Queries       Queries       `json:"queries" db:"queries"`
	Headers       Headers       `json:"headers" db:"headers"`
	PermOrigins   PermOrigins   `json:"permOrigins" db:"perm_origins"`
	RequestConfig RequestConfig `json:"requestConfig" db:"request_config"`
}

func JSONUnmarshaller(src interface{}, dst interface{}) (err error) {
	buf := bytes.NewBuffer(src.([]byte))

	err = json.Unmarshal(buf.Bytes(), dst)
	if err != nil {
		return
	}

	return nil
}

type Headers []Header

func (he *Headers) Scan(src interface{}) (err error) {
	err = JSONUnmarshaller(src, he)
	return
}

type Queries []Query

func (qu *Queries) Scan(src interface{}) (err error) {
	err = JSONUnmarshaller(src, qu)
	return
}

type PermOrigins []PermOrigin

func (pe *PermOrigins) Scan(src interface{}) (err error) {
	err = JSONUnmarshaller(src, pe)
	return
}

func (c *Conn) GetRouteRequestData(projectName, routePath string) (endpoint Endpoint, err error) {
	stmt, err := c.Db.Preparex(fmt.Sprintf(`
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
		),

		perm_origin_values(id, route_id, url) AS
		(
			SELECT * FROM perm_origin
			WHERE route_id = (SELECT id from route_id)
		),

		request_config_values(id, route_id, merge_header, merge_body, merge_query) AS 
		(
			SELECT id, route_id, merge_header, merge_body, merge_query FROM request_config
			WHERE route_id = (SELECT id FROM route_id)
		)

		SELECT
			route.*,
			array_to_json(array_remove(array_agg(distinct(header_values)), NULL)) AS "headers",
			array_to_json(array_remove(array_agg(distinct(query_values)), NULL)) AS "queries",
			array_to_json(array_remove(array_agg(distinct(perm_origin_values)), NULL)) AS "perm_origins",
			row_to_json(request_config_values) AS "request_config"
		FROM route
		LEFT JOIN header_values ON header_values.route_id = route.id
		LEFT JOIN query_values ON query_values.route_id = route.id
		LEFT JOIN perm_origin_values ON perm_origin_values.route_id = route.id
		LEFT JOIN request_config_values ON request_config_values.route_id = route.id
		WHERE route.project_id = (
				SELECT id FROM project_values
			)
		AND route.id = (
			SELECT id from route_id
		)
		GROUP BY route.id, request_config_values.*;
	`, PG_CRYPT_KEY))

	if err != nil {
		return
	}

	err = stmt.Get(&endpoint, projectName, routePath)
	if err != nil {
		return
	}

	return
}
