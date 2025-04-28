package match

import (
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/service/user"
	response "backend/internal/utils/http"
)

type Service struct {
	repo     *repositories.MatchesRepository
	userRepo *repositories.UserRepository
}

// NewMatchesService creates a new instance of MatchesService.
func NewMatchesService(repo *repositories.MatchesRepository, repository *repositories.UserRepository) *Service {
	return &Service{repo: repo, userRepo: repository}
}

func (service *Service) Match(uuid1 string, uuid2 string) *response.Response {

	err := service.repo.CreateMatch(uuid1, uuid2)
	if err != nil {
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "")
}

func (service *Service) GetMatches(userID string) *response.Response {
	allIds, err := service.repo.GetMatches(userID)

	if err != nil {
		return response.ErrorResponse("Internal error occurred")
	}

	id, err := service.userRepo.GetUserByID(*allIds...)
	if err != nil {
		return response.ErrorResponse("Interal error occurred")
	}

	var users []models.UserInfoModel

	for _, rawUser := range *id {
		info := user.GetUserInfo(&rawUser)
		info.Email = rawUser.Email
		users = append(users, *info)
	}

	return response.SuccessResponse(users, "")

}
