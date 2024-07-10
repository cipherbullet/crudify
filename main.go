package main

import (
	"log"
)

func main() {
    pgStorage := &PostgresStorage{}
	connStr := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Tehran"
	if err := pgStorage.Connect(connStr); err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

	migrate()
	
	s := NewServer("0.0.0.0", "8080", pgStorage)
	s.Start()
}