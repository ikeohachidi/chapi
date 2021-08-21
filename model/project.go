package model

type Project struct {
	Id        uint   `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	UserId    uint   `json:"userId" db:"user_id"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func (conn *Conn) CreateProject(name string, userId uint) (projectId uint, err error) {
	stmt, err := conn.db.Preparex(`INSERT INTO project ("name", user_id) values($1, $2) RETURNING id`)
	if err != nil {
		return
	}

	row := stmt.QueryRowx(name)

	row.Scan(&projectId)

	stmt.Close()

	return
}

func (conn *Conn) ListProjects() (projects []Project, err error) {
	err = conn.db.Select(&projects, "SELECT * FROM project")

	return
}

// TODO: can't figure out how to retrieve this properly
func (conn *Conn) GetProjects(userId uint) (projects []Project, err error) {
	query := `
		SELECT	project.*, 
				route.*, 
				"query".*
		FROM project
		LEFT OUTER JOIN route ON project.id = route.project_id
		LEFT OUTER JOIN "query" on route.id = "query".route_id
		WHERE project.user_id = $1 
	`
	err = conn.db.Select(&projects, query, userId)

	return
}

func (conn *Conn) GetUserProjects(userId uint) (projects []Project, err error) {
	query := `SELECT * FROM project WHERE user_id = $1`

	err = conn.db.Select(&projects, query, userId)

	if err != nil {
		return
	}

	return
}

func (conn *Conn) DeleteProject(projectId uint, userId uint) (err error) {
	_, err = conn.db.Exec("DELETE FROM project WHERE id=$1 AND user_id=$2", projectId, userId)

	return
}
