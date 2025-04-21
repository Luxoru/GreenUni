package user

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"github.com/google/uuid"
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

func (service *Service) GetUser() {

}

func (service *Service) GetRawUserByID(userID uuid.UUID) (*models.RawUserRow, error) {
	repository := service.repo
	return repository.GetUserByID(userID, mysql.QueryOptions{})
}

// GetUserByID retrieves a user by their UUID.
func (service *Service) GetUserByID(userID uuid.UUID) (interface{}, error) {

	repository := service.repo

	rawUser, err := repository.GetUserByID(userID, mysql.QueryOptions{})
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
