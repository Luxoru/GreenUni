package pathapi

import (
	"backend/internal/db/adapters/mysql"
	"github.com/go-chi/chi"
)

type PathComponent interface {
	SetupComponents(sqlRepository *mysql.Repository) chi.Router
}
