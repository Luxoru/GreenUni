package auth

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/service/auth"
	"backend/internal/utils/http"
	"backend/routes/pathapi"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type LoginPath struct {
	router  chi.Router
	service *auth.Service
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (path *LoginPath) SetupComponents(sqlRepository *mysql.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r
	r.Post("/", path.Login)
	repository, err := repositories.NewUserRepository(sqlRepository)
	if err != nil {
		panic(err)
		return nil
	}
	path.service = auth.NewAuthService(repository)
	return r
}

func (path *LoginPath) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		path.service.Login("", "")
		return
	}

	loginStatus := path.service.Login(req.Username, req.Password)

	response.WriteJson(w, loginStatus)
}

func LoginRoute() pathapi.PathComponent {
	return &LoginPath{}
}
