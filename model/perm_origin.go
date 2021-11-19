package model

import "github.com/jmoiron/sqlx"

type PermOrigin struct {
	ID      uint   `json:"id" db:"id"`
	URL     string `json:"url" db:"url"`
	RouteID uint   `json:"routeId" db:"route_id"`
}

func (p *PermOrigin) Create(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`INSERT INTO perm_origin(url, route_id) VALUES($1, $2) RETURNING id`)
	if err != nil {
		return
	}

	row := stmt.QueryRowx(p.URL, p.RouteID)

	err = row.Scan(p.ID)

	return
}

func (p *PermOrigin) FetchAll(db *sqlx.DB) ([]PermOrigin, error) {
	origins := []PermOrigin{}

	stmt, err := db.Preparex(`SELECT * FROM perm_origin WHERE route_id = $1`)
	if err != nil {
		return nil, err
	}

	err = stmt.Select(&origins, p.RouteID)
	if err != nil {
		return nil, err
	}

	return origins, nil
}

func (p *PermOrigin) Update(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`UPDATE perm_origin SET url = $1 WHERE id = $2 AND route_id = $3`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(p.URL, p.ID, p.RouteID)
	return
}

func (p *PermOrigin) Delete(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`DELETE FROM perm_origin WHERE id = $1 AND route_id = $2`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(p.ID, p.RouteID)

	return
}
