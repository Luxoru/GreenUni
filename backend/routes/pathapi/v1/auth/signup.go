package auth

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/adapters/neo4j"
	"backend/internal/db/repositories"
	"backend/internal/service/auth"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
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

func (path *SignupPath) SetupComponents(sqlRepository *mysql.Repository, _ *neo4j.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r

	userRepo, err := repositories.NewUserRepository(sqlRepository)
	if err != nil {
		return nil
	}

	path.service = auth.NewAuthService(userRepo)

	r.Post("/", path.Signup)
	return r
}

func (path *SignupPath) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJson(w, response.ErrorResponse("Invalid request body"))
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" || req.Role == "" {
		response.WriteJson(w, response.ErrorResponse("All fields are required"))
		return
	}

	status := path.service.Signup(req.Username, req.Password, req.Email, req.Role)
	response.WriteJson(w, status)
}

func SignupRoute() pathapi.PathComponent {
	return &SignupPath{}
}
