package model

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Project struct {
	ID        uint   `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	UserID    uint   `json:"userId" db:"user_id"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func (p *Project) Create(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`
		WITH e AS (
			INSERT INTO project ("name", user_id) 
			VALUES($1, $2) 
			ON CONFLICT("name") DO NOTHING 
			RETURNING id
		)

		SELECT * FROM e
		UNION
		SELECT id FROM project WHERE "name"=$1
	`)
	if err != nil {
		return
	}

	row := stmt.QueryRowx(strings.ToLower(p.Name), p.UserID)

	err = row.Scan(&p.ID)
	if err != nil {
		return
	}

	stmt.Close()

	return
}

func (p *Project) Update(db *sqlx.DB) (err error) {
	stmt, err := db.Preparex(`UPDATE project SET "name"=$1 WHERE id=$2 AND user_id=$3`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(strings.ToLower(p.Name), p.ID, p.UserID)

	return
}

func (p *Project) ProjectExists(db *sqlx.DB) (exists bool, err error) {
	stmt, err := db.Preparex(`
		SELECT EXISTS (
			SELECT *
			FROM project
			WHERE "name" ILIKE '%' || $1 || '%'
		)
	`)

	if err != nil {
		return
	}

	row := stmt.QueryRowx(p.Name)

	err = row.Scan(&exists)

	stmt.Close()

	return
}

func ListProjects(db *sqlx.DB) ([]Project, error) {
	projects := []Project{}

	err := db.Select(&projects, "SELECT * FROM project")
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Project) GetProjectByName(db *sqlx.DB) (project Project, err error) {
	err = db.Select(&project, fmt.Sprintf("SELECT * FROM project WHERE name = %v", p.Name))

	return
}

// TODO: can't figure out how to retrieve this properly
func (p *Project) GetProjects(db *sqlx.DB) ([]Project, error) {
	projects := []Project{}

	query := `
		SELECT	project.*, 
				route.*, 
				"query".*
		FROM project
		LEFT OUTER JOIN route ON project.id = route.project_id
		LEFT OUTER JOIN "query" on route.id = "query".route_id
		WHERE project.user_id = $1 
	`
	err := db.Select(&projects, query, p.UserID)
	if err != nil {
		return nil, err
	}

	return projects, err
}

func (p *Project) GetUserProjects(db *sqlx.DB) ([]Project, error) {
	projects := []Project{}

	query := `SELECT * FROM project WHERE user_id = $1`

	err := db.Select(&projects, query, p.UserID)

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Project) DeleteProject(db *sqlx.DB) (err error) {
	_, err = db.Exec("DELETE FROM project WHERE id=$1 AND user_id=$2", p.ID, p.UserID)

	return
}
