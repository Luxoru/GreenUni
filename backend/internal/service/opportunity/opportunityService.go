package opportunity

import (
	"backend/internal/db/repositories"
	"backend/internal/models"
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

// DeleteOpportunity deletes an opportunity by its UUID.
func (service *OpportunityService) DeleteOpportunity(opportunityUUID uuid.UUID) error {
	return service.repo.DeleteOpportunity(opportunityUUID)
}

func (service *OpportunityService) GetOpportunitiesByTag(tagName string) (*[]models.OpportunityModel, error) {
	return service.repo.GetOpportunitiesByTag(tagName)
}

func (service *OpportunityService) GetOpportunitiesFrom(from string, limit string) (*[]models.OpportunityModel, int64, error) {

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

	return service.repo.GetOpportunitiesFrom(fromInt, limitInt)
}
