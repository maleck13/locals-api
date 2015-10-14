package data
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)


func connectionString()(string){
	user:= os.Getenv("MYSQL_USER");
	pass:= os.Getenv("MYSQL_PASS");
	host:= os.Getenv("MYSQL_HOST");
	dbName:= os.Getenv("MYSQL_DB_NAME");
	return user+":"+pass+"@tcp("+host+":3306)/"+dbName;
}

func handleFatal(err error){
	if err != nil {
		log.Fatal(err)
	}
}


/**
  this function returns a long lived object. The connection is closed when the program exits
 */

var db *sql.DB;
func DataBaseConnection()(*sql.DB){
	var err error;
	if nil != db{
		err = db.Ping();
		if nil != err{
			db=nil;
			return DataBaseConnection();
		}
		return db;
	}
	db,err = sql.Open("mysql", connectionString())
	handleFatal(err);
	err=db.Ping();
	handleFatal(err);
	log.Println("mysql connection ok")
	return db;
}
