package db

import (
	"database/sql"
	"fmt"
	"log"
)

func Connectdb(USER, PASSWD, HOSTDB, DBNAME, PORT string) *sql.DB {
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
