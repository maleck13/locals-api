package data_test

import (
	"database/sql"
	"errors"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/maleck13/locals-api/data"
	"testing"
)

var columns []string = []string{"id", "email", "created", "county", "type", "profilePic", "bio", "interests", "phone"}

func TestProfileShouldSave(t *testing.T) {

	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()
	profile.Email = "test@test.com"
	profile.County = "waterford"
	profile.Type = "local"
	var prep *sqlmock.ExpectedPrepare = mock.ExpectPrepare("INSERT INTO profile")
	prep.ExpectExec().WithArgs(profile.Email, profile.County, profile.Type).
		WillReturnResult(sqlmock.NewResult(1, 1))
	profile, err = profile.Save()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if nil == profile {
		t.Fatal("no profile")
	}

}

func TestFindById(t *testing.T) {

	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()

	mock.ExpectQuery("^SELECT (.+) FROM profile WHERE id=?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test@test.com,13232323232,waterford,local,somepic,test,test, 342343242"))
	prof, err := profile.FindById(1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when doing find", err)
	}

	t.Log(prof)
	if prof.Id != 1 {
		t.Fatal("wrong id")
	}
	if prof.Email != "test@test.com" {
		t.Fatal("wrong email")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

func TestExists(t *testing.T) {
	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()
	mock.ExpectQuery("SELECT (.+) as count FROM profile WHERE email=?").
		WithArgs("test@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	exists := profile.Exists("test@test.com")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if !exists {
		t.Fail()
		t.Fatalf("an error '%s' was not expected when checking exists", exists)
	}
}

func TestExistsFails(t *testing.T) {
	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()
	mock.ExpectQuery("SELECT (.+) as count FROM profile WHERE email=?").
		WithArgs("test@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	exists := profile.Exists("test@test.com")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if exists {
		t.Fail()
		t.Fatalf("an error '%s' was not expected when checking exists", exists)
	}
}

func TestFindByEmail(t *testing.T) {
	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()
	mock.ExpectQuery("^SELECT (.+) FROM profile WHERE email=?").
		WithArgs("test@test.com").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test@test.com,13232323232,waterford,local,somepic,test,test, 342343242"))
	prof, err := profile.FindByEmail("test@test.com")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if "test@test.com" != prof.Email {
		t.Fail()
		t.Fatal("expected email to be test@test.com")
	}

}

func TestFindByEmailError(t *testing.T) {
	mock, profile, err := getMockProfile()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer profile.Store().Close()
	mock.ExpectQuery("^SELECT (.+) FROM profile WHERE email=?").
		WithArgs("test@test.com").
		WillReturnError(errors.New("error with query"))
	_, err = profile.FindByEmail("test@test.com")
	if nil == err {
		t.Fail()
		t.Fatal("expected an error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

func getMockProfile() (sqlmock.Sqlmock, *data.Profile, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	var getDb data.DataBase = func() *sql.DB {
		return db
	}
	profile := data.NewProfile()
	profile.Store = getDb
	return mock, profile, err
}
