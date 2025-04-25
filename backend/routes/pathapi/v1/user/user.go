package user

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/security"
	"backend/internal/service/user"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

type Path struct {
	router  chi.Router
	service *user.Service
}

func (path *Path) SetupComponents(sqlRepository *mysql.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r
	r.Get("/me", path.GetCurrentUser)
	r.Get("/{id}", path.GetUserByID)
	r.Get("/username/{username}", path.GetUserByUsername)

	//r.Delete("/", path.DeleteUser)

	//Photos upload stuff
	//r.Put("/me/photos", path.AddPhoto)       //Add photo to profile
	//r.Get("/photos", path.GetPhoto)          //Get photo from profile
	//r.Delete("/me/photos", path.DeletePhoto) //Delete photo from profile

	repository, err := repositories.NewUserRepository(sqlRepository)
	if err != nil {
		return nil
	}
	path.service = user.NewUserService(repository)
	return r
}
func (path *Path) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		response.WriteJson(w, response.ErrorResponse("User ID is required"))
		return
	}

	if _, err := uuid.Parse(idStr); err != nil {
		response.WriteJson(w, response.ErrorResponse("Invalid UUID format"))
		return
	}

	req := &models.GetUserRequest{UserUUID: idStr}
	result := path.service.GetUser(req)
	response.WriteJson(w, result)
}

func (path *Path) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		response.WriteJson(w, response.ErrorResponse("Username is required"))
		return
	}

	req := &models.GetUserRequest{Username: username}
	result := path.service.GetUser(req)
	response.WriteJson(w, result)
}

func (path *Path) GetCurrentUser(writer http.ResponseWriter, request *http.Request) {
	userInfo, err := security.ExtractUserInfoFromJWT(request)
	if err != nil {
		response.WriteJson(writer, response.ErrorResponse("Unauthorized"))
		return
	}

	req := &models.GetUserRequest{Username: userInfo.Username}
	res := path.service.GetUser(req)
	response.WriteJson(writer, res)
}

func Route() pathapi.PathComponent {
	return &Path{}
}
