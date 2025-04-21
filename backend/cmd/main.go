// Package main is the entry point for the backend API server.
package main

import (
	"backend/internal/db"
	"backend/internal/db/adapters/mysql"
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

	// Setup router
	router := chi.NewRouter()
	log.SetReportCaller(true)

	handlers.Handler(router, &repo)

	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		panic(err)
	}
}
