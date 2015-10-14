package data
import (
	"database/sql"
)

type DataBase func()*sql.DB;

type Profile struct {
	Id int				`json:"id"`
	Email string        `json:"email"`
	Created int64       `json:"created"`
	County string       `json:"county"`
	Type string			`json:"type"`
	Store DataBase   	`json:"-"`
}

func NewProfile() * Profile{
	return &Profile{Store:DataBaseConnection}
}


type NoResult  struct{
	Err string
}


func (nf * NoResult)Code()int{
	return 404;
}

func (nf * NoResult)Error()string{
	return "no result for query" + nf.Err
}



const sql_TABLE_INSERT = "INSERT INTO profile(email,county,type) VALUES(?,?,?)";
const sql_SELECT_BY_ID = "SELECT id, email, UNIX_TIMESTAMP(created),county, type FROM profile WhERE id=?";
const sql_SELECT_BY_EMAIL = "SELECT id, email, UNIX_TIMESTAMP(created),county, type FROM profile WHERE email=?";
const sql_SELECT_BY_COUNTY = "SELECT id, email, UNIX_TIMESTAMP(created),county, type FROM profile WHERE county=?";
const sql_PROFILE_EXISTS  = "SELECT COUNT(email) as count FROM profile WHERE email=?";

func (p *Profile) Save()(*Profile,error){
	var(
		session * sql.DB
		err error
		stmt * sql.Stmt
	)
	session= p.Store();

	stmt,err = session.Prepare(sql_TABLE_INSERT)
	if nil != err {
		return p,err;
	}
	_,err= stmt.Exec(p.Email,p.County,p.Type);
	if nil != err {
		return p,err;
	}

	return p,nil;
}

func (p *Profile) FindById(id int64) (Profile, error){
	var(
		session * sql.DB
	)

	session = p.Store()
	row:= session.QueryRow(sql_SELECT_BY_ID,id);
	return single(row);
}

func (p *Profile) Exists(email string)bool{
	var(
		session * sql.DB
	)

	session = p.Store()
	row:= session.QueryRow(sql_PROFILE_EXISTS,email);
	var count int;
	row.Scan(&count);
	return count > 0
}

func (p * Profile) FindByEmail(email string)(Profile,error){
	var(
		session * sql.DB
	)

	session = p.Store()
	row := session.QueryRow(sql_SELECT_BY_EMAIL,email);

	return single(row);
}

func (p * Profile) FindByCounty(county string) ([]Profile,error){
	var(
		err error
		session * sql.DB
		profiles []Profile
	)

	session = p.Store()
	rows,err := session.Query(sql_SELECT_BY_COUNTY,county);
	if nil != err{
		return nil,err
	}
	profiles,err = list(rows);
	if nil != err{
		return nil,err
	}
	return profiles,err

}

func single(row * sql.Row)(Profile,error){
	var profile Profile;
	var err error
	err = row.Scan(&profile.Id,&profile.Email,&profile.Created, &profile.County, &profile.Type);
	if 0 == profile.Id && nil != err{
		return profile,&NoResult{Err:err.Error()}
	}
	return profile,nil;
}

func list(rows *sql.Rows)([]Profile,error){
	profiles := make([]Profile,0);
	for rows.Next(){
		var p Profile
		err := rows.Scan(&p.Id,&p.Email, &p.Created,&p.County,&p.Type)
		if err != nil {
			return nil,err;
		}
		profiles = append(profiles, p)
	}
	return profiles,nil
}