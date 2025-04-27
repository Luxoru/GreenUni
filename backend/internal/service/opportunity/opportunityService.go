package opportunity

import (
	"backend/internal/db/repositories"
	"backend/internal/models"
	response "backend/internal/utils/http"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// OpportunityService provides methods for managing opportunities.
type OpportunityService struct {
	repo *repositories.OpportunityRepository
}

// NewOpportunityService creates a new instance of OpportunityService.
func NewOpportunityService(repo *repositories.OpportunityRepository) *OpportunityService {
	return &OpportunityService{repo: repo}
}

// CreateOpportunity creates a new opportunity with the given details.
func (service *OpportunityService) CreateOpportunity(request models.CreateOpportunityRequest) *models.CreateOpportunityStatus {

	if request.Title == "" || request.Type == "" || request.Description == "" || request.Location == "" || request.AuthorUUID == "" || request.Points == 0 {
		return writeStatus(nil, "Invalid request format", false)
	}

	opportunityUUID := uuid.New()

	postedByUUID, err := uuid.Parse(request.AuthorUUID) //TODO: maybe confirm this boi exists?

	if err != nil {
		return writeStatus(nil, "unable to parse postedByUUID", false)
	}

	opportunityModel := models.OpportunityModel{
		UUID:            opportunityUUID,
		Title:           request.Title,
		Description:     request.Description,
		Points:          request.Points,
		Location:        request.Location,
		OpportunityType: request.Type,
		PostedByUUID:    postedByUUID,
	}

	var modelTags []models.TagModel

	for _, tag := range request.Tags {
		model := models.TagModel{
			TagName: tag,
		}

		modelTags = append(modelTags, model)
	}
	var media []models.MediaModel

	types := request.MediaTypes
	urls := request.MediaURLs

	for i := range types {

		parseMediaType, err := models.ParseMediaType(types[i])
		if err != nil {
			log.Error(err)
			return writeStatus(nil, "Internal error occured whilst processing media type", false)
		}

		media = append(media, models.MediaModel{
			Type: parseMediaType,
			URL:  urls[i],
		})
	}

	opportunityModel.Tags = &modelTags
	opportunityModel.Media = &media

	err = service.repo.CreateOpportunity(&opportunityModel)
	if err != nil {
		log.Error(err)
		return writeStatus(nil, "Internal error occurred whilst creating opportunity", false)
	}

	return writeStatus(&opportunityModel, "", true)
}

func writeStatus(model *models.OpportunityModel, message string, success bool) *models.CreateOpportunityStatus {
	return &models.CreateOpportunityStatus{
		OpportunityModel: model,
		Success:          success,
		Message:          message,
	}
}

func (service *OpportunityService) UpdateOpportunity(request models.CreateOpportunityRequest) *response.Response {
	if request.Title == "" || request.Type == "" || request.Description == "" || request.Location == "" || request.AuthorUUID == "" || request.Points == 0 || request.UUID == "" {
		return response.ErrorResponse("Invalid request format")
	}

	opportunityUUID, err := uuid.Parse(request.UUID)
	if err != nil {
		return response.ErrorResponse("Unable to parse UUID")
	}

	model := &models.OpportunityModel{
		UUID:            opportunityUUID,
		Title:           request.Title,
		Description:     request.Description,
		Points:          request.Points,
		Location:        request.Location,
		OpportunityType: request.Type,
	}

	var modelTags []models.TagModel

	for _, tag := range request.Tags {
		model := models.TagModel{
			TagName: tag,
		}

		modelTags = append(modelTags, model)
	}
	var media []models.MediaModel

	types := request.MediaTypes
	urls := request.MediaURLs

	for i := range types {

		parseMediaType, err := models.ParseMediaType(types[i])
		if err != nil {
			log.Error(err)
			return response.ErrorResponse("Internal error occured whilst processing media type")
		}

		media = append(media, models.MediaModel{
			Type: parseMediaType,
			URL:  urls[i],
		})
	}

	model.Tags = &modelTags
	model.Media = &media

	err = service.repo.UpdateOpportunity(model)
	if err != nil {
		return response.ErrorResponse("Internal error occured")
	}

	return response.SuccessResponse(model, "")
}

func (service *OpportunityService) UpdateStatus(opportunityID string, status string) *response.Response {

	if opportunityID == "" {
		return response.ErrorResponse("Opportunity UUID not provided")
	}

	opportunityUUID, err := uuid.Parse(opportunityID)

	if err != nil {
		return response.ErrorResponse("Unable to parse UUID")
	}

	opportunityStatus, err := strconv.ParseBool(status)

	if err != nil {
		return response.ErrorResponse("Unable to parse status")
	}

	err = service.repo.UpdateOpportunityStatus(opportunityUUID, opportunityStatus)
	if err != nil {
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "")
}

// GetOpportunity retrieves an opportunity by its UUID.
func (service *OpportunityService) GetOpportunity(opportunityUUID uuid.UUID) (*models.OpportunityModel, error) {
	opportunity, err := service.repo.GetOpportunity(&opportunityUUID)
	if err != nil {
		return nil, err
	}
	if opportunity == nil {
		return nil, nil
	}
	model := *opportunity
	return &model[0], nil
}

func (service *OpportunityService) GetOpportunitiesByAuthor(authorID string) *response.Response {
	if authorID == "" {
		return response.ErrorResponse("Author uuid not provided")
	}

	authorUUID, err := uuid.Parse(authorID)

	if err != nil {
		return response.ErrorResponse("Unable to parse uuid")
	}

	oppportunities, err := service.repo.GetOpportunityByAuthor(&authorUUID)
	if err != nil {
		return response.ErrorResponse("Internal error occured")
	}

	return response.SuccessResponse(oppportunities, "")
}

// DeleteOpportunity deletes an opportunity by its UUID.
func (service *OpportunityService) DeleteOpportunity(opportunityUUID uuid.UUID) error {
	return service.repo.DeleteOpportunity(opportunityUUID)
}

func (service *OpportunityService) GetOpportunitiesByTag(tagName string) (*[]models.OpportunityModel, error) {
	return service.repo.GetOpportunitiesByTag(tagName)
}

func (service *OpportunityService) GetOpportunitiesFrom(from string, limit string, userUUID uuid.UUID) (*[]models.OpportunityModel, int64, error) {

	var fromInt int64
	var limitInt int64
	fromInt, err := strconv.ParseInt(from, 10, 64)

	if err != nil {
		return nil, 0, errors.New("unable limit parse 'from' as an integer")
	}

	limitInt, err = strconv.ParseInt(limit, 10, 64)

	if err != nil {
		return nil, 0, errors.New("unable limit parse 'limit' as an integer")
	}

	return service.repo.GetOpportunitiesFrom(fromInt, limitInt, userUUID)
}

func (service *OpportunityService) GetOpportunityByLikes(opportunityID string, from string, limit string) *response.Response {

	if opportunityID == "" || from == "" || limit == "" {
		return response.ErrorResponse("must be provided")
	}

	opportunityUUID, err := uuid.Parse(opportunityID)

	if err != nil {
		return response.ErrorResponse("Unable to parse Opportunity UUID")
	}

	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		return response.ErrorResponse("from is not an integer")
	}

	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return response.ErrorResponse("limit is not an integer")
	}

	likes, lastRow, err := service.repo.GetOpportunityByLikes(opportunityUUID, fromInt, limitInt)
	if err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occured")
	}

	if likes == nil {
		return response.ErrorResponse("No likes for given opportunity")
	}

	model := struct {
		Likes     *[]*models.StudentInfoModel `json:"likes"`
		LastIndex int64                       `json:"lastIndex"`
	}{likes, lastRow}

	return response.SuccessResponse(model, "")
}

func (service *OpportunityService) LikeOpportunity(userID, postID string) *response.Response {
	userUUID, postUUID, errResp := parseUUIDs(userID, postID)
	if errResp != nil {
		return errResp
	}

	if err := service.repo.LikeOpportunity(userUUID, postUUID); err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "Successfully liked opportunity")
}

func (service *OpportunityService) DislikeOpportunity(userID, postID string) *response.Response {
	userUUID, postUUID, errResp := parseUUIDs(userID, postID)
	if errResp != nil {
		return errResp
	}

	if err := service.repo.DislikeOpportunity(userUUID, postUUID); err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "Successfully disliked opportunity")
}

func (service *OpportunityService) DeleteLikeOpportunity(userID, postID string) *response.Response {
	userUUID, postUUID, errResp := parseUUIDs(userID, postID)
	if errResp != nil {
		return errResp
	}

	if err := service.repo.DeleteLikeOpportunity(userUUID, postUUID); err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "Successfully deleted liked opportunity")
}

func (service *OpportunityService) DeleteDislikeOpportunity(userID, postID string) *response.Response {
	userUUID, postUUID, errResp := parseUUIDs(userID, postID)
	if errResp != nil {
		return errResp
	}

	if err := service.repo.DeleteDislikeOpportunity(userUUID, postUUID); err != nil {
		log.Error(err)
		return response.ErrorResponse("Internal error occurred")
	}

	return response.SuccessResponse(nil, "Successfully deleted disliked opportunity")
}

func parseUUIDs(userID, postID string) (uuid.UUID, uuid.UUID, *response.Response) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, response.ErrorResponse("Unable to parse user uuid")
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return uuid.Nil, uuid.Nil, response.ErrorResponse("Unable to parse post uuid")
	}

	return userUUID, postUUID, nil
}
