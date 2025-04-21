package auth

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/service/auth"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type SignupPath struct {
	router  chi.Router
	service *auth.Service
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (path *SignupPath) SetupComponents(sqlRepository *mysql.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r
	r.Post("/", path.Signup)
	repository, err := repositories.NewUserRepository(sqlRepository)
	if err != nil {
		panic(err)
		return nil
	}
	path.service = auth.NewAuthService(repository)
	return r
}

func (path *SignupPath) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		path.service.Signup("", "", "", "")
		return
	}

	status := path.service.Signup(req.Username, req.Password, req.Email, req.Role)
	response.WriteJson(w, status)

}

func SignupRoute() pathapi.PathComponent {
	return &SignupPath{}
}
