package handlers

import (
	"backend/internal/db/adapters/mysql"
	"backend/routes/pathapi"
	"backend/routes/pathapi/v1/auth"
	"backend/routes/pathapi/v1/opportunities"
	"backend/routes/pathapi/v1/root"
	"backend/routes/pathapi/v1/student"
	"backend/routes/pathapi/v1/user"
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// routeRegistry maps API versions to their respective routes and handlers
var routeRegistry = map[string]map[string]func() pathapi.PathComponent{
	"v1": {
		"/":              root.Route,
		"/user":          user.Route,
		"/auth/signup":   auth.SignupRoute,
		"/auth/login":    auth.LoginRoute,
		"/opportunities": opportunities.OpportunityRoute,
		"/student":       student.Route,
	},
}

// Handler sets up the routing for the application using the Chi router
func Handler(r *chi.Mux, sqlRepository *mysql.Repository) {
	r.Use(chimiddle.StripSlashes)
	for version, routes := range routeRegistry {
		r.Route("/api/"+version, func(v chi.Router) {
			for path, routeFunc := range routes {
				components := routeFunc().SetupComponents(sqlRepository)
				if components == nil {
					log.Warnf("Unable to setup path for %s", path)
				}
				v.Mount(path, components)
			}
		})
	}
}
