package pathapi

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/adapters/neo4j"
	"github.com/go-chi/chi"
)

type PathComponent interface {
	SetupComponents(sqlRepository *mysql.Repository, repository *neo4j.Repository) chi.Router
}
