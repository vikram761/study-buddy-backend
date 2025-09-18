package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func Connectdb(USER, PASSWD, HOSTDB, DBNAME, PORT string) *sql.DB {
	log.Println(USER,PASSWD, HOSTDB, DBNAME, PORT)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", USER, PASSWD, HOSTDB, PORT, DBNAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("connection is successful")
	return db
}
