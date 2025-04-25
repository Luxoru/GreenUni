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
	"strconv"
)

type Path struct {
	router  chi.Router
	service *opportunity.OpportunityService
}

func (path *Path) SetupComponents(repo *mysql.Repository) chi.Router {
	r := chi.NewRouter()

	repository, err := repositories.NewOpportunityRepository(repo)
	if err != nil {
		log.Fatal("Failed to initialize OpportunityRepository: ", err)
	}

	path.service = opportunity.NewOpportunityService(repository)

	r.Get("/", path.GetOpportunities)
	r.Post("/", path.CreateOpportunity)
	r.Delete("/{uuid}", path.DeleteOpportunity)
	r.Get("/{uuid}", path.GetByUUID)

	path.router = r
	return r
}

func (path *Path) CreateOpportunity(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJson(w, response.ErrorResponse("Invalid request body"))
		return
	}

	status := path.service.CreateOpportunity(req)
	response.WriteJson(w, status)
}

func (path *Path) GetOpportunities(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	switch {
	case query.Has("tag"):
		path.getByTag(w, query.Get("tag"))
	case query.Has("uuid"):
		path.getByUUID(w, query.Get("uuid"))
	case query.Has("from") && query.Has("limit"):
		path.getPaginated(w, query.Get("from"), query.Get("limit"))
	default:
		response.WriteJson(w, response.ErrorResponse("Missing query: expected ?tag=, ?uuid=, or ?from=&limit="))
	}
}

func (path *Path) getByTag(w http.ResponseWriter, tag string) {
	opportunities, err := path.service.GetOpportunitiesByTag(tag)
	if err != nil {
		log.Error("GetOpportunitiesByTag error: ", err)
		response.WriteJson(w, response.ErrorResponse("Internal error occurred"))
		return
	}
	if opportunities == nil {
		response.WriteJson(w, response.ErrorResponse("No opportunities found for given tag"))
		return
	}

	response.WriteJson(w, response.SuccessResponse(opportunities, ""))
}

func (path *Path) GetByUUID(writer http.ResponseWriter, request *http.Request) {
	uuidStr := chi.URLParam(request, "uuid")
	path.getByUUID(writer, uuidStr)
}

func (path *Path) getByUUID(w http.ResponseWriter, uuidStr string) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.WriteJson(w, response.ErrorResponse("Invalid UUID format"))
		return
	}

	opp, err := path.service.GetOpportunity(id)
	if err != nil {
		log.Error("GetOpportunity error: ", err)
		response.WriteJson(w, response.ErrorResponse("Internal error occurred"))
		return
	}
	if opp == nil {
		response.WriteJson(w, response.ErrorResponse("Opportunity not found"))
		return
	}

	response.WriteJson(w, response.SuccessResponse([]models.OpportunityModel{*opp}, ""))
}

func (path *Path) getPaginated(w http.ResponseWriter, fromStr, limitStr string) {
	from := fromStr
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.WriteJson(w, response.ErrorResponse("Invalid 'limit' value"))
		return
	}

	opportunities, lastIndex, err := path.service.GetOpportunitiesFrom(from, limitStr)
	if err != nil {
		log.Error("Pagination error: ", err)
		response.WriteJson(w, response.ErrorResponse("Failed to retrieve opportunities"))
		return
	}

	payload := map[string]interface{}{
		"success":   true,
		"data":      opportunities,
		"lastIndex": lastIndex,
	}
	response.WriteJson(w, payload)
}

func (path *Path) DeleteOpportunity(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		response.WriteJson(w, response.ErrorResponse("UUID is required"))
		return
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.WriteJson(w, response.ErrorResponse("Invalid UUID format"))
		return
	}

	if err := path.service.DeleteOpportunity(id); err != nil {
		log.Error("DeleteOpportunity error: ", err)
		response.WriteJson(w, response.ErrorResponse("Internal error occurred"))
		return
	}

	response.WriteJson(w, response.SuccessResponse(nil, "Opportunity deleted"))
}

func OpportunityRoute() pathapi.PathComponent {
	return &Path{}
}
