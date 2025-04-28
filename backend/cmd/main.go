// Package main is the entry point for the backend API server.
package main

import (
	"backend/internal/db"
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/adapters/neo4j"
	"backend/internal/handlers"
	"fmt"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// main is the main entry point of the backend API.
func main() {
	fmt.Println("Hello world")

	// Setup databases and repositories

	// Mysql
	//TODO: maybe load from config
	container := mysql.Container{}
	config := mysql.Configurations{
		Authentication: &db.AuthenticationConfigurations{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "password",
		},
		DatabaseName: "mydatabase",
	}

	err := container.Connect(config)
	if err != nil {
		panic(err)
	}

	repo := mysql.Repository{
		Database: &container,
	}

	neoContainer := neo4j.Container{}
	neoConfig := neo4j.Configurations{
		Authentication: &db.URIConfigurations{
			AuthConfig: db.AuthenticationConfigurations{
				Username: "neo4j",
				Password: "password",
			},
			URI: "bolt://localhost:7687",
		},
	}

	err = neoContainer.Connect(neoConfig)
	if err != nil {
		panic(err)
	}

	neoRepo := neo4j.Repository{
		Database: &neoContainer,
	}

	// Setup router
	router := chi.NewRouter()
	log.SetReportCaller(true)

	handlers.Handler(router, &repo, &neoRepo)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
