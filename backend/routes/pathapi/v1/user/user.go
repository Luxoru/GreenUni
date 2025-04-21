package user

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/service/user"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"github.com/go-chi/chi"
	"net/http"
)

type Path struct {
	router  chi.Router
	service *user.Service
}

func (path *Path) SetupComponents(sqlRepository *mysql.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r
	r.Get("/", path.GetUser) // user
	//r.Delete("/", path.DeleteUser)
	//r.Put("/me", path.GetCurrentUser)

	//Photos upload stuff
	//r.Put("/me/photos", path.AddPhoto)       //Add photo to profile
	//r.Get("/photos", path.GetPhoto)          //Get photo from profile
	//r.Delete("/me/photos", path.DeletePhoto) //Delete photo from profile

	repository, err := repositories.NewUserRepository(sqlRepository)
	if err != nil {
		panic(err)
		return nil
	}
	path.service = user.NewUserService(repository)
	return r
}

func (path *Path) GetUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("userID")
	username := query.Get("username")

	model := models.GetUserRequest{
		UserUUID: userID,
		Username: username,
	}

	getUser := path.service.GetUser(&model)

	response.WriteJson(w, getUser)

}

func Route() pathapi.PathComponent {
	return &Path{}
}
