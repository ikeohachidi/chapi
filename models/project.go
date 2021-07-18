package models

type Project struct {
	Id     uint   `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	UserId string `json:"-" db:"user_id"`
}

func (conn *Conn) CreateProject(name string) (projectId uint, err error) {
	stmt, err := conn.db.Preparex("INSERT INTO project (name) values($1) RETURNING id")
	if err != nil {
		return
	}

	row, _ := stmt.Query(name)
	if err != nil {
		return
	}

	row.Scan(&projectId)

	stmt.Close()

	return
}

func (conn *Conn) ListProjects() (projects []Project, err error) {
	err = conn.db.Select(&projects, "SELECT * FROM projects")

	return
}

func (conn *Conn) GetUserProjects(userId string) (projects []Project, err error) {
	err = conn.db.Select(&projects, "SELECT * FROM projects WHERE user_id=$1", userId)

	return
}

func (conn *Conn) DeleteProject(projectId uint) (err error) {
	_, err = conn.db.Exec("DELETE FROM projects WHERE id=$1", projectId)

	return
}
