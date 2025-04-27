package student

import (
	"backend/internal/db/repositories"
	"backend/internal/models"
	response "backend/internal/utils/http"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	repo *repositories.StudentRepository
}

// NewStudentService creates a new instance of UserService.
func NewStudentService(repo *repositories.StudentRepository) *Service {
	return &Service{repo: repo}
}

func (service *Service) UpdateStudentInfo(studentInfo models.StudentInfoModel) *response.Response {
	err := service.repo.UpdateStudentInfo(studentInfo)
	if err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "")
}

func (service *Service) GetStudentInfo(userID string) *response.Response {

	if userID == "" {
		return response.ErrorResponse("User UUID must be provided")
	}

	userUUID, err := uuid.Parse(userID)

	if err != nil {
		return response.ErrorResponse("Unable to parse UUID")
	}

	return service.GetStudentInfoByUUID(userUUID)
}

func (service *Service) GetStudentInfoByUUID(userUUID uuid.UUID) *response.Response {
	info, err := service.repo.GetUserInfo(userUUID)
	if err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	if info == nil {
		return response.ErrorResponse("User doesn't exist")
	}

	return response.SuccessResponse(info, "")
}
