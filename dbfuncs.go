package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "idea-evolver-instance.csffmzjjjhky.us-east-1.rds.amazonaws.com"
	port     = 54323
	user     = "applicant_john"
	dbname   = "idea_evolver"
	password = "ideaevolverpass"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err 
}
