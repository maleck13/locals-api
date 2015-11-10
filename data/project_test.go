package data_test

import (
	"database/sql"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/maleck13/locals-api/data"
	"testing"
)

func TestNewProject(t *testing.T) {
	project := data.NewProject()
	assert.NotNil(t, project, "project should not be nil")
}

func TestProjectShouldSave(t *testing.T) {

	mock, project, err := getMockProject()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer project.Store().Close()
	project.Title = "test project"
	project.UserId = 1
	var prep *sqlmock.ExpectedPrepare = mock.ExpectPrepare("INSERT INTO project")
	prep.ExpectExec().WithArgs(project.Title, project.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	project, err = project.Save()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if nil == project {
		t.Fatal("no profile")
	}

	assert.Equal(t, "test project", project.Title, "expected title to be the same")

}

func getMockProject() (sqlmock.Sqlmock, *data.Project, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	var getDb data.DataBase = func() *sql.DB {
		return db
	}
	project := data.NewProject()
	project.Store = getDb
	return mock, project, err
}
