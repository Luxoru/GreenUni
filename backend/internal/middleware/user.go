package middleware

import (
	"backend/internal/models"
	"backend/internal/security"
	"net/http"
)

func CheckIfAdminUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		userInfo, err := security.ExtractUserInfoFromJWT(request)
		if err != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if userInfo == nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if userInfo.Role != models.Admin {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(writer, request)
	})
}
