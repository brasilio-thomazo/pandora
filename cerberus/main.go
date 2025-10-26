package main

import (
	"log"

	"bitbucket.org/brasilio/pandora/cerberus/config"
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/server"
)

func main() {
	config := config.NewConfig()
	db, err := database.NewPool(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := server.Httpd(config, db); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
