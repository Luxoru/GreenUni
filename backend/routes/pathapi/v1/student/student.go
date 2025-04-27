package student

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/security"
	"backend/internal/service/student"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type Path struct {
	router  chi.Router
	service *student.Service
}

func (path *Path) SetupComponents(sqlRepository *mysql.Repository) chi.Router {
	r := chi.NewRouter()
	path.router = r
	r.Get("/me", path.GetCurrentUserInfo)
	r.Put("/me", path.UpdateCurrentUserInfo)

	repository, err := repositories.NewStudentRepository(sqlRepository)
	if err != nil {
		return nil
	}
	path.service = student.NewStudentService(repository)
	return r
}

func (path *Path) GetCurrentUserInfo(writer http.ResponseWriter, request *http.Request) {
	userInfo, err := security.ExtractUserInfoFromJWT(request)
	if err != nil {
		response.WriteJson(writer, response.ErrorResponse("Unauthorized"))
		return
	}

	returnMessage := path.service.GetStudentInfoByUUID(userInfo.UUID)

	response.WriteJson(writer, returnMessage)
}

func (path *Path) UpdateCurrentUserInfo(writer http.ResponseWriter, request *http.Request) {
	var studentInfo models.StudentInfoModel
	if err := json.NewDecoder(request.Body).Decode(&studentInfo); err != nil {
		response.WriteJson(writer, response.ErrorResponse("Invalid request body"))
		return
	}

	returnMessage := path.service.UpdateStudentInfo(studentInfo)

	response.WriteJson(writer, returnMessage)
}

func Route() pathapi.PathComponent {
	return &Path{}
}
