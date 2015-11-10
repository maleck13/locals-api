package data

import (
	"database/sql"
	"log"
)

type Project struct {
	Id          int      `json:"id"`
	UserId      int      `json:"userid,string"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Store       DataBase `json:"-"`
}

const (
	_project_fields          = "id, title, userid, description"
	sql_PROJECT_TABLE_INSERT = "INSERT INTO `locals`.`project`(`title`,`userid`,`description`)VALUES(?,?,?);"
	sql_PROJECT_SELECT_BY_ID = "SELECT " + _project_fields + " FROM project WHERE id=?"
	sql_PROJECT_LIST         = "SELECT " + _project_fields + " FROM project LIMIT ?,?"
)

func NewProject() *Project {
	return &Project{Store: DataBaseConnection}
}

func (p *Project) Save() (*Project, error) {
	var (
		session *sql.DB
		err     error
		stmt    *sql.Stmt
	)
	session = p.Store()
	stmt, err = session.Prepare(sql_PROJECT_TABLE_INSERT)
	if nil != err {
		return p, err
	}
	_, err = stmt.Exec(p.Title, p.UserId)
	if nil != err {
		return p, err
	}
	return p, nil
}

func (p *Project) FindById(id int64) (*Project, error) {
	var (
		session *sql.DB
	)
	log.Printf("finding project in db with id %d", id)
	session = p.Store()
	row := session.QueryRow(sql_PROJECT_SELECT_BY_ID, id)
	return singleProject(row)
}

func (p *Project) ListProjects(from, to int) ([]*Project, error) {
	var (
		err      error
		session  *sql.DB
		projects []*Project
	)

	session = p.Store()
	rows, err := session.Query(sql_PROJECT_LIST, from, to)
	if nil != err {
		return nil, err
	}
	projects, err = listProject(rows)
	if nil != err {
		return nil, err
	}
	return projects, err
}

func singleProject(row *sql.Row) (*Project, error) {
	var project *Project
	var err error
	project = NewProject()
	var pProjectText []byte

	err = row.Scan(&project.Id, &project.Title, &project.UserId, &pProjectText)
	if 0 == project.Id && nil != err {
		return project, &NoResult{Err: err.Error()}
	}
	if 0 == project.Id {
		return nil, &NoResult{Err: "no row for id"}
	}
	project.Description = string(pProjectText)
	return project, nil
}

func listProject(rows *sql.Rows) ([]*Project, error) {
	projects := make([]*Project, 0)
	for rows.Next() {
		project := NewProject()
		var pProjectText []byte
		err := rows.Scan(&project.Id, &project.Title, &project.UserId, &pProjectText)

		project.Description = string(pProjectText)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
