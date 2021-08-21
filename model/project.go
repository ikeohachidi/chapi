package model

type Project struct {
	Id        uint    `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	UserId    uint    `json:"userId" db:"user_id"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
	Routes    []Route `json:"routes" db:"routes"`
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
	err = conn.db.Select(&projects, "SELECT * FROM projects")

	return
}

func (conn *Conn) GetUserProjects(userId uint) (projects []Project, err error) {
	query := `
		SELECT p.*, r.*, q.*
		FROM project as p WHERE user_id = $1
		INNER JOIN route AS r ON p.id = r.project_id
		INNER JOIN query AS q ON r.id = q.route_id
	`
	err = conn.db.Select(&projects, query, userId)

	return
}

func (conn *Conn) DeleteProject(projectId uint, userId uint) (err error) {
	_, err = conn.db.Exec("DELETE FROM projects WHERE id=$1 AND user_id=$2", projectId, userId)

	return
}
