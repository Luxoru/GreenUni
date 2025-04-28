package root

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/adapters/neo4j"
	"backend/routes/pathapi"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type Path struct {
	router chi.Router
}

type MessyData struct {
	ID   int
	Name string
}

func (path *Path) SetupComponents(_ *mysql.Repository, _ *neo4j.Repository) chi.Router {
	r := chi.NewRouter()
	r.Get("/", GetDefault)
	path.router = r
	return r
}

func GetDefault(w http.ResponseWriter, r *http.Request) {
	data := MessyData{
		ID:   1,
		Name: "Test",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func Route() pathapi.PathComponent {
	return &Path{}
}
