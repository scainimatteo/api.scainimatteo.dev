package services

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("errore in sql.Open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("errore nel ping del database: %v", err)
	}

	log.Println("Connessione Postgres riuscita (via lib/pq)!")
	return db, nil

}
