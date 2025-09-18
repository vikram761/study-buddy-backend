package main

import (
	"log"
	"os"
	"study-buddy-backend/routes"
	"study-buddy-backend/services/db"

	"github.com/joho/godotenv"
)


func main() {
	// LOAD ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	username := os.Getenv("USERNAME")
	passwd := os.Getenv("PASSWD")
	hostDB := os.Getenv("HOSTDB")
	dbName := os.Getenv("DBNAME")
    port := os.Getenv("DBPORT")
	log.Println("not connected")

	// DB INIT
	database := db.Connectdb(username, passwd, hostDB, dbName, port)
	log.Println("connected")
	db.CreateDB(database)

	appPort := os.Getenv("PORT")
    if appPort == "" {
      appPort = "8080"
    }
	
	// ROUTES INIT AND RUNNING THE SERVER
	routes.InitRoutes(appPort,database)
}
