package model

type MergeOptions struct {
	Id          uint `json:"id", db:"id"`
	RouteId     uint `json:"route_id" db:"route_id"`
	MergeHeader bool `json:"merge_header" db:"merge_header"`
	MergeBody   bool `json:"merge_body" db:"merge_body"`
	MergeQuery  bool `json:"merge_query" db:"merge_query"`
}

func (mo *MergeOptions) Scan(src interface{}) (err error) {
	err = JSONUnmarshaller(src, mo)
	return
}
