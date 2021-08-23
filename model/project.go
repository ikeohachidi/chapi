package model

type Project struct {
	ID        uint   `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	UserID    uint   `json:"userId" db:"user_id"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func (conn *Conn) CreateProject(name string, userID uint) (projectID uint, err error) {
	stmt, err := conn.db.Preparex(`INSERT INTO project ("name", user_id) values($1, $2) RETURNING id`)
	if err != nil {
		return
	}

	row := stmt.QueryRowx(name)

	row.Scan(&projectID)

	stmt.Close()

	return
}

func (conn *Conn) ProjectExists(name string) (exists bool, err error) {
	// SELECT id FROM project WHERE name = $1
	stmt, err := conn.db.Preparex(`
		SELECT
			CASE 
				WHEN "name" IS NOT NULL THEN 1
				ELSE 0
			END "exists"
		FROM project WHERE "name" =$1 
	`)

	if err != nil {
		return
	}

	row := stmt.QueryRowx(name)

	row.Scan(&exists)

	stmt.Close()

	return
}

func (conn *Conn) ListProjects() (projects []Project, err error) {
	err = conn.db.Select(&projects, "SELECT * FROM project")

	return
}

// TODO: can't figure out how to retrieve this properly
func (conn *Conn) GetProjects(userID uint) (projects []Project, err error) {
	query := `
		SELECT	project.*, 
				route.*, 
				"query".*
		FROM project
		LEFT OUTER JOIN route ON project.id = route.project_id
		LEFT OUTER JOIN "query" on route.id = "query".route_id
		WHERE project.user_id = $1 
	`
	err = conn.db.Select(&projects, query, userID)

	return
}

func (conn *Conn) GetUserProjects(userID uint) (projects []Project, err error) {
	query := `SELECT * FROM project WHERE user_id = $1`

	err = conn.db.Select(&projects, query, userID)

	if err != nil {
		return
	}

	return
}

func (conn *Conn) DeleteProject(projectID uint, userID uint) (err error) {
	_, err = conn.db.Exec("DELETE FROM project WHERE id=$1 AND user_id=$2", projectID, userID)

	return
}
