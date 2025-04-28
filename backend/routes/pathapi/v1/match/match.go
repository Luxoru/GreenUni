package match

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/adapters/neo4j"
	"backend/internal/db/repositories"
	"backend/internal/service/match"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Path struct {
	router  chi.Router
	service *match.Service
}

func (path *Path) SetupComponents(sql *mysql.Repository, repo *neo4j.Repository) chi.Router {
	r := chi.NewRouter()

	repository, err := repositories.NewMatchesRepository(repo)
	if err != nil {
		log.Fatal("Failed to initialize Matchese Repo: ", err)
		return nil
	}
	userRepo, err := repositories.NewUserRepository(sql)

	if err != nil {
		log.Fatal("Failed to initialize User Repo: ", err)
		return nil
	}

	r.Post("/", path.CreateMatch)
	r.Get("/{userID}", path.GetMatches)

	path.service = match.NewMatchesService(repository, userRepo)

	path.router = r
	return r
}

func (path *Path) CreateMatch(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	uuid1 := query.Get("uuid1")
	uuid2 := query.Get("uuid2")

	res := path.service.Match(uuid1, uuid2)
	response.WriteJson(writer, res)

}

func (path *Path) GetMatches(writer http.ResponseWriter, request *http.Request) {
	uuidStr := chi.URLParam(request, "userID")

	res := path.service.GetMatches(uuidStr)

	response.WriteJson(writer, res)

}

func Route() pathapi.PathComponent {
	return &Path{}
}
