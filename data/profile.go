package data
import (
	"database/sql"
	"errors"
	"log"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

type DataBase func()*sql.DB;


type Profile struct {
	Id int						`json:"id"`
	Email string        		`json:"email"`
	Created int64       		`json:"created"`
	County string       		`json:"county"`
	Type string					`json:"type"`
	ProfilePic string			`json:"profilePic"`
	Interests string            `json:"interests"`
	Bio string                   `json:"bio"`
	Phone string                `json:"phone"`
	RegisterToken string        `json:"-"`
	Store DataBase   			`json:"-"`
}

func NewProfile() * Profile{
	return &Profile{Store:DataBaseConnection}
}

func FindProfileById( id int64) (* Profile,error){
	var(
		session * sql.DB
	)
	p := NewProfile()
	session = p.Store()
	row:= session.QueryRow(sql_SELECT_BY_ID,id);
	return single(row);
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


const _fields = "id, email, UNIX_TIMESTAMP(created),county, type, profilePic, bio, interests, phone, registerToken";
const sql_TABLE_INSERT = "INSERT INTO profile(email,county,type,registerToken) VALUES(?,?,?,?)";
const sql_SELECT_BY_ID = "SELECT "+_fields+" FROM profile WHERE id=?";
const sql_SELECT_BY_EMAIL = "SELECT "+_fields+" FROM profile WHERE email=?";
const sql_SELECT_BY_COUNTY = "SELECT "+_fields+" FROM profile WHERE county=?";
const sql_PROFILE_EXISTS  = "SELECT COUNT(email) as count FROM profile WHERE email=?";
const sql_UPDATE_PROFILE_PIC = "UPDATE profile set profilePic=? WHERE id=?";
const sql_UPDATE_PROFILE = "UPDATE profile SET email=?,county=?,type=?,bio=?,interests=?, phone=? WHERE id=?"
const sql_SELECT_BY_REGISTER_TOKEN = "SELECT " + _fields + " FROM profile WHERE registerToken=?"


func (p * Profile) getRegisterToken()(string , error){
	var(
		err error
		hexString string
	)
	if 0 == p.Id || 0 == p.Created{
		err = errors.New("no id or created field")
	}
	hash:= md5.New();
	fields:= strconv.Itoa(p.Id) + p.Email
	_,err = hash.Write([]byte(fields))
	if nil != err{
		return "",err
	}
	hexString =  hex.EncodeToString(hash.Sum(nil))
	return hexString,nil;
}


func (p *Profile) findByRegisterToken(registerToken string)(*Profile,error){
	var(
		session * sql.DB
	)
	session= p.Store();
	row:= session.QueryRow(sql_SELECT_BY_REGISTER_TOKEN,registerToken);
	return single(row)
}

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
	token,err := p.getRegisterToken()
	if nil != err {
		return p,err;
	}
	_,err= stmt.Exec(p.Email,p.County,p.Type,token);
	if nil != err {
		return p,err;
	}
	return p,nil;
}

func (p *Profile) FindById(id int64) (* Profile, error){
	var(
		session * sql.DB
	)
	log.Printf("finding profile in db with id %i", id)
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

func (p * Profile) FindByEmail(email string)(* Profile,error){
	var(
		session * sql.DB
	)

	session = p.Store()
	row := session.QueryRow(sql_SELECT_BY_EMAIL,email);

	return single(row);
}

func (p * Profile) FindByCounty(county string) ([]*Profile,error){
	var(
		err error
		session * sql.DB
		profiles []*Profile
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

func (p*Profile) UpdateProfilePic(location string)(error){
	var(
		session *sql.DB
		error error
		stmt * sql.Stmt
	)
	session = p.Store();
	stmt,error = session.Prepare(sql_UPDATE_PROFILE_PIC)
	if nil != error {
		return error;
	}

	if 0 == p.Id{
		return errors.New("no id for profile during update")
	}
	_,error= stmt.Exec(location,p.Id);
	return error

}

func (p*Profile) Update()(error){
	var(
		session *sql.DB
		error error
		stmt * sql.Stmt
	)
	session = p.Store();
	stmt,error = session.Prepare(sql_UPDATE_PROFILE)
	if nil != error {
		return error;
	}
	_,error= stmt.Exec(p.Email,p.County,p.Type,p.Bio,p.Interests,p.Phone,p.Id);
	return error
}

func single(row * sql.Row)(*Profile,error){
	var profile * Profile;
	var err error
	profile = NewProfile()
	var pProfile []byte;
	var pInterests []byte;
	var pPhone []byte;

	err = row.Scan(&profile.Id,&profile.Email,&profile.Created, &profile.County, &profile.Type, &pProfile,&profile.Bio,&pInterests,&pPhone);
	if 0 == profile.Id && nil != err{
		return profile,&NoResult{Err:err.Error()}
	}
	if 0 == profile.Id{
		return nil, &NoResult{Err:"no row for id"};
	}
	profile.ProfilePic = string(pProfile)
	profile.Interests = string(pInterests)
	profile.Phone = string(pPhone)
	return profile,nil;
}

func list(rows *sql.Rows)([]*Profile,error){
	profiles := make([]*Profile,0);
	for rows.Next(){
		p := NewProfile()
		var pPic []byte;
		err := rows.Scan(&p.Id,&p.Email,&p.Created,&p.County,&p.Type,&pPic)

		p.ProfilePic = string(pPic);

		if err != nil {
			return nil,err;
		}
		profiles = append(profiles, p)
	}
	return profiles,nil
}

