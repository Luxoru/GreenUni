package opportunities

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/service/opportunity"
	response "backend/internal/utils/http"
	"backend/routes/pathapi"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Path struct {
	router  chi.Router
	service *opportunity.OpportunityService
}

func (path *Path) SetupComponents(repo *mysql.Repository) chi.Router {
	r := chi.NewRouter()

	repository, err := repositories.NewOpportunityRepository(repo)
	if err != nil {
		log.Error(err)
		return nil
	}

	path.service = opportunity.NewOpportunityService(repository)
	r.Get("/", path.GetOpportunity)
	r.Post("/", path.CreateOpportunity)
	r.Delete("/", path.DeleteOpportunity)
	//r.Get("/feed", path.GetOpportunityFeed) //Gets feed based on liked tags etc. -> Most likely handled by tiktok api?

	path.router = r
	return r
}

func (path *Path) CreateOpportunity(writer http.ResponseWriter, request *http.Request) {

	var req models.CreateOpportunityRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		response.WriteJson(writer, response.ErrorResponse("invalid JSON body"))
		return
	}

	status := path.service.CreateOpportunity(req)

	response.WriteJson(writer, status)

}

type GetOpportunityResponseModel struct {
	Models  *[]models.OpportunityModel `json:"models"`
	Success bool                       `json:"Success"`
	Message string                     `json:"Message"`
}

func (path *Path) GetOpportunity(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postUUID := query.Get("uuid")
	tagName := query.Get("tag")

	switch {
	case tagName != "":
		path.getByTag(w, tagName)
	case postUUID != "":
		path.getByUUID(w, postUUID)
	default:
		response.WriteJson(w, response.ErrorResponse("no postUUID or tagName provided"))
	}
}

func (path *Path) getByTag(w http.ResponseWriter, tagName string) {

	model := GetOpportunityResponseModel{}
	opportunities, err := path.service.GetOpportunitiesByTag(tagName)
	if err != nil {
		log.Error(err)
		model.Message = "Internal error occurred"
		response.WriteJson(w, model)
		return
	}
	if opportunities == nil {
		model.Message = "no opportunities with given tags"
		response.WriteJson(w, model)
		return
	}

	model.Success = true
	model.Models = opportunities

	response.WriteJson(w, model)
}

func (path *Path) getByUUID(w http.ResponseWriter, uuidStr string) {
	model := GetOpportunityResponseModel{}
	postID, err := uuid.Parse(uuidStr)
	if err != nil {
		model.Message = "unable to parse uuid"
		response.WriteJson(w, model)
		return
	}

	post, err := path.service.GetOpportunity(postID)
	if err != nil {
		log.Error(err)
		model.Message = "Internal error occurred"
		response.WriteJson(w, model)
		return
	}
	if post == nil {
		model.Message = "post doesn't exist"
		response.WriteJson(w, model)
		return
	}

	model.Success = true
	opportunityModels := []models.OpportunityModel{*post}
	model.Models = &opportunityModels

	response.WriteJson(w, model)
}

func (path *Path) DeleteOpportunity(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")

	model := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{}

	if uuidStr == "" {
		model.Message = "post uuid not provided"
		response.WriteJson(w, model)
		return
	}

	postUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		model.Message = "unable to parse uuid"
		response.WriteJson(w, model)
		return
	}

	if err := path.service.DeleteOpportunity(postUUID); err != nil {
		log.Error(err)
		model.Message = "internal error occurred"
		response.WriteJson(w, model)
		return
	}

	model.Success = true

	response.WriteJson(w, model)
}

func OpportunityRoute() pathapi.PathComponent {
	return &Path{}
}
