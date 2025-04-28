package user

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Service provides methods for managing users.
type Service struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo *repositories.UserRepository) *Service {
	return &Service{repo: repo}
}

func GetUserInfo(rawModel *models.RawUserRow) *models.UserInfoModel {
	if rawModel == nil {
		return nil
	}

	return &models.UserInfoModel{
		UUID:     rawModel.UUID,
		Username: rawModel.Username,
		Role:     rawModel.Role,
	}
}

func GetStudentModel(rawModel *models.RawUserRow) *models.StudentModel {
	if rawModel == nil {
		return nil
	}

	if !rawModel.Points.Valid {
		return nil
	}

	return &models.StudentModel{
		UserInfoModel: GetUserInfo(rawModel),
		Points:        rawModel.Points.Int64,
	}
}

func GetRecruiterModel(rawModel *models.RawUserRow) *models.RecruiterModel {

	if rawModel == nil {
		return nil
	}

	if !rawModel.OrganisationName.Valid || !rawModel.ApplicationStatus.Valid {
		return nil
	}

	return &models.RecruiterModel{
		UserInfoModel:     GetUserInfo(rawModel),
		OrganisationName:  rawModel.OrganisationName.String,
		ApplicationStatus: rawModel.ApplicationStatus.Bool,
	}
}

func (service *Service) GetRawUserByName(username string) (*models.RawUserRow, error) {
	repository := service.repo
	return repository.GetUserByName(username, mysql.QueryOptions{})
}

// GetUserByName retrieves a user by their username.
func (service *Service) GetUserByName(username string) (interface{}, error) {
	repository := service.repo

	rawUser, err := repository.GetUserByName(username, mysql.QueryOptions{})
	if err != nil {
		return nil, err
	}

	info := &models.UserInfoModel{
		UUID:     rawUser.UUID,
		Username: rawUser.Username,
		Role:     rawUser.Role,
	}

	switch rawUser.Role {
	case models.Student:
		return &models.StudentModel{
			UserInfoModel: info,
			Points:        rawUser.Points.Int64,
		}, nil

	case models.Recruiter:
		return &models.RecruiterModel{
			UserInfoModel:     info,
			OrganisationName:  rawUser.OrganisationName.String,
			ApplicationStatus: rawUser.ApplicationStatus.Bool,
		}, nil

	default:
		return info, nil
	}
}

func (service *Service) GetUser(request *models.GetUserRequest) *models.GetUserResponse {

	if request.UserUUID == "" && request.Username == "" {
		return writeStatus(nil, "username/userID not provided", false)
	}

	if request.UserUUID != "" {
		parsedUUID, err := uuid.Parse(request.UserUUID)
		if err != nil {
			return writeStatus(nil, "uuid couldn't be parsed", false)
		}
		user, err := service.GetRawUserByID(parsedUUID)
		if err != nil {
			log.Errorf("error occured: %s", err)
			return writeStatus(nil, "error occured fetching user", false)
		}
		if user == nil {
			return writeStatus(nil, "user with id "+request.UserUUID+" doesnt exist", false)
		}

		return writeStatus(GetUserInfo(user), "", true)
	}

	if request.Username == "" {
		return writeStatus(nil, "invalid request. Username/UserID not defined", false)
	}

	user, err := service.GetRawUserByName(request.Username)
	if err != nil {
		log.Errorf("error occured: %s", err)
		return writeStatus(nil, "error occured fetching user", false)
	}
	if user == nil {
		return writeStatus(nil, "user with name "+request.Username+" doesnt exist", false)
	}

	return writeStatus(GetUserInfo(user), "", true)

}

func writeStatus(model *models.UserInfoModel, message string, success bool) *models.GetUserResponse {
	return &models.GetUserResponse{
		Success: success,
		Message: message,
		Data:    model,
	}
}

func (service *Service) GetRawUserByID(userID uuid.UUID) (*models.RawUserRow, error) {
	repository := service.repo
	id, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	row := *id

	return &row[0], err
}

// GetUserByID retrieves a user by their UUID.
func (service *Service) GetUserByID(userID uuid.UUID) (interface{}, error) {

	repository := service.repo

	rawUser, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	user := *rawUser

	info := &models.UserInfoModel{
		UUID:     user[0].UUID,
		Username: user[0].Username,
		Role:     user[0].Role,
	}

	switch user[0].Role {
	case models.Student:
		return &models.StudentModel{
			UserInfoModel: info,
			Points:        user[0].Points.Int64,
		}, nil

	case models.Recruiter:
		return &models.RecruiterModel{
			UserInfoModel:     info,
			OrganisationName:  user[0].OrganisationName.String,
			ApplicationStatus: user[0].ApplicationStatus.Bool,
		}, nil

	default:
		return info, nil
	}
}
