package auth

import "backend/internal/models"

type LoginStatus struct {
	LoginAllowed bool        `json:"success"`
	Message      string      `json:"message"`
	Info         interface{} `json:"info"`
}

type SignupStatus struct {
	Model   *models.UserInfoModel `json:"user"`
	Message string                `json:"message"`
	Success bool                  `json:"success"`
}
