package user

import (
	"backend/api"
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/service/user"
	"backend/routes/pathapi"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

	path.service.GetUser()

	var user *models.RawUserRow

	if userID != "" {
		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			err = fmt.Errorf("uuid couldn't be parsed")
			api.RequestErrorhandler(w, err)
			return
		}
		user, err = path.service.GetRawUserByID(parsedUUID)
		if err != nil {
			log.Errorf("error occured: %s", err)
			err = fmt.Errorf("error occured fetching user")
			api.RequestErrorhandler(w, err)
			return
		}
		if user == nil {
			err = fmt.Errorf("user with id %s doesnt exist", parsedUUID)
			api.RequestErrorhandler(w, err)
			return
		}

	} else {
		username := query.Get("username")
		var err error
		if username == "" {
			err = fmt.Errorf("invalid request. Username/UserID not defined")
			api.RequestErrorhandler(w, err)
			return
		}

		user, err = path.service.GetRawUserByName(username)
		if err != nil {
			log.Errorf("error occured: %s", err)
			err = fmt.Errorf("error occured fetching user")
			api.RequestErrorhandler(w, err)
			return
		}
		if user == nil {
			err = fmt.Errorf("user with name %s doesnt exist", username)
			api.RequestErrorhandler(w, err)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")

	userInfo := &models.UserInfoModel{
		UUID:     user.UUID,
		Username: user.Username,
	}

	if userInfo == nil {
		err := fmt.Errorf("ended up nill")
		api.RequestErrorhandler(w, err)
		return
	}

	json.NewEncoder(w).Encode(userInfo)
}

func Route() pathapi.PathComponent {
	return &Path{}
}
