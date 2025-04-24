package auth

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/db/repositories"
	"backend/internal/models"
	"backend/internal/models/auth"
	"backend/internal/security"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	repo *repositories.UserRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo *repositories.UserRepository) *Service {
	return &Service{repo: repo}
}

// Signup creates a new user with the given username and password.
func (service *Service) Signup(username string, password string, email string, role string) *auth.SignupStatus {

	if username == "" || password == "" || email == "" || role == "" {
		return setSignupStatus(nil, "Invalid request format", false)
	}

	name, err := service.repo.GetUserByName(username, mysql.QueryOptions{})
	if err != nil {
		return setSignupStatus(nil, "error creating account", false)
	}

	if name != nil {
		return setSignupStatus(nil, "An account with this name already exists", false)
	}

	name, err = service.repo.GetUserByEmail(email, mysql.QueryOptions{})
	if err != nil {
		return setSignupStatus(nil, "error creating account", false)
	}

	if name != nil {
		return setSignupStatus(nil, "An account with this email already exists", false)
	}

	salt, err := security.GenerateSalt(16)
	if err != nil {
		return setSignupStatus(nil, "error creating account", false)
	}

	hashedPassword, err := security.HashPassword(password, salt)
	if err != nil {
		return setSignupStatus(nil, "error creating account", false)
	}

	userUUID := uuid.New()

	parsedRole, err := models.ParseRoleType(role)

	if err != nil {
		return setSignupStatus(nil, "role parsed doesn't exist", false)
	}

	user := &models.UserModel{
		UUID:           userUUID,
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		Salt:           salt,
		Role:           parsedRole,
	}

	err = service.repo.AddUser(user, mysql.InsertOptions{})
	if err != nil {
		return setSignupStatus(nil, "error creating account", false)
	}

	info := &models.UserInfoModel{
		UUID:     userUUID,
		Username: username,
		Role:     parsedRole,
	}

	return setSignupStatus(info, "success", true)
}

func setSignupStatus(model *models.UserInfoModel, message string, loginSuccess bool) *auth.SignupStatus {
	signupStatus := auth.SignupStatus{
		Model:   model,
		Message: message,
		Success: loginSuccess,
	}
	return &signupStatus
}

func (service *Service) Login(username string, password string) *auth.LoginStatus {

	if username == "" || password == "" {
		return setLoginStatus("invalid request format", false)
	}

	user, err := service.repo.GetUserByName(username, mysql.QueryOptions{})
	if err != nil {
		return setLoginStatus("internal error occurred whilst fetching user", false)
	}

	if user == nil {
		return setLoginStatus("user doesn't exist", false)
	}
	actualPass := user.HashedPassword
	salt := user.Salt

	isCorrectPassword := security.VerifyPassword(password, salt, actualPass)

	if !isCorrectPassword {
		return setLoginStatus("invalid username or password", false)
	}

	//Create JWT token
	userInfo := &models.UserInfoModel{
		Username: user.Username,
		UUID:     user.UUID,
		Role:     user.Role,
	}

	jwt, err := security.GenerateJWT(userInfo)
	if err != nil {
		log.Error(err)
		return setLoginStatus("error occured creating login token", false)
	}

	s := struct {
		Token string                `json:"token"`
		User  *models.UserInfoModel `json:"user"`
	}{jwt, userInfo}

	status := auth.LoginStatus{
		LoginAllowed: true,
		Message:      "",
		Info:         s,
	}

	return &status
}

func setLoginStatus(message string, loginSuccess bool) *auth.LoginStatus {
	return &auth.LoginStatus{
		Message:      message,
		LoginAllowed: loginSuccess,
	}
}
