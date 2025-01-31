package auth

import(
	"fmt"
	"log"
	"time"
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDB() *Database {
	database := &Database{}
	database.connect()
	
	return database
}

func (d *Database) connect() {
	var err error
	creds := LoadCredentials()	
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", creds.User, creds.Pass, creds.Host, creds.Port, creds.Database)
	
	//Open connection
	d.db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err, " - AUTH OPEN")
	}
	
	//Must have
	d.db.SetConnMaxLifetime(time.Minute * 3)
	d.db.SetMaxOpenConns(10)
	d.db.SetMaxIdleConns(10)
	
	//Make sure connection is real
	err = d.db.Ping()
	if err != nil {
		log.Fatal(err, " - AUTH PING")
	}
}

func (d *Database) query(queryString string, args ...any) *sql.Rows {
	resp, err := d.db.Query(queryString, args)
	if err != nil {
		log.Println(err, " - AUTH Query")
	}
	
	return resp
}

func (d *Database) queryNoArgs(queryString string) *sql.Rows {
	resp, err := d.db.Query(queryString)
	if err != nil {
		log.Println(err, " - AUTH Query")
	}
	
	return resp
}