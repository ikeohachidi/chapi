package model

func (c *Conn) GetRouteRequestData(projectName, routePath string) (route Route, err error) {
	stmt := `
		SELECT
			route.*,
			array_agg(json_build_object('id', "query".id, 'name', "query"."name", 'value', "query"."value")) as queries,
			array_agg(json_build_object('id', header.id, 'name', header."name", 'value', header."value")) as headers
		FROM route
		INNER JOIN "query" ON "query".route_id = route.id
		INNER JOIN "header" ON "header".route_id = "query".id
		WHERE route.project_id = (
				SELECT id FROM project WHERE "name" = $1 
			)
		AND path = $2 
		GROUP BY route.id
	`

	row := c.db.QueryRow(stmt, projectName, routePath)

	err = row.Scan(&route)
	if err != nil {
		return
	}

	return
}
