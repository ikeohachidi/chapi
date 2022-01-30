package model

import (
	"github.com/jmoiron/sqlx"
)

type MergeOptions struct {
	ID          uint `json:"id" db:"id"`
	RouteId     uint `json:"routeId" db:"route_id"`
	MergeHeader bool `json:"mergeHeader" db:"merge_header"`
	MergeBody   bool `json:"mergeBody" db:"merge_body"`
	MergeQuery  bool `json:"mergeQuery" db:"merge_query"`
}

func (mo *MergeOptions) Scan(src interface{}) (err error) {
	err = JSONUnmarshaller(src, mo)
	return
}

func (mo *MergeOptions) GetRouteMergeOptions(db *sqlx.DB) (err error) {
	row := db.QueryRow(`SELECT merge_header, merge_query, merge_body FROM merge_options WHERE route_id = $1`, mo.RouteId)

	if err = row.Scan(&mo.MergeHeader, &mo.MergeQuery, &mo.MergeBody); err != nil {
		return
	}

	return
}

func (mo *MergeOptions) SaveMergeOptions(db *sqlx.DB) (err error) {
	_, err = db.NamedExec(`
		UPDATE merge_options
		SET merge_body = :merge_body,
			merge_header = :merge_header,
			merge_query = :merge_query
		WHERE route_id = :route_id
	`, mo)

	if err != nil {
		return
	}

	return
}
