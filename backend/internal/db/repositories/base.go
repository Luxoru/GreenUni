package repositories

import (
	"backend/internal/db/adapters/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

//TODO: implement caching layer
//TODO: make async

// SQLRepository Interface for repositories that require SQL table creation queries
type SQLRepository interface {
	CreateTablesQuery() *[]string
	CreateIndexesQuery() *[]string
}

// BaseRepository Struct that holds a reference to a MySQL Repository
type BaseRepository struct {
	Repository *mysql.Repository
}

// InitRepository Initializes the repository, runs SQL queries to create tables
func InitRepository(repo SQLRepository, container *mysql.Repository) (*BaseRepository, error) {

	repository := container

	base := &BaseRepository{
		Repository: repository,
	}

	createTableQueries := *repo.CreateTablesQuery()

	// Execute SQL query for each create table query
	for _, query := range createTableQueries {
		rows, err := repository.ExecuteQuery(query, make([]mysql.Column, 0), mysql.QueryOptions{})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		rows.Close()
	}

	//Manage indexes

	createTableIndexes := *repo.CreateIndexesQuery()

	for _, query := range createTableIndexes {
		_, err := repository.ExecuteQuery(query, nil, mysql.QueryOptions{})
		if err != nil {

			if strings.Contains(err.Error(), "Duplicate key name") {
				// The index already exists — can ignore error
				continue
			}
			return nil, err

		}
	}

	return base, nil
}
