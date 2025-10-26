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
	defer db.Close()

	dbc := database.NewCache(config)
	defer dbc.Close()

	if err := server.Httpd(config, db, dbc); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
