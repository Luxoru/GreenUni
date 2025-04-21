package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type RawUserRow struct {
	UUID              uuid.UUID
	Username          string
	Email             string
	HashedPassword    string
	Salt              string
	Role              RoleType
	OrganisationName  sql.NullString
	ApplicationStatus sql.NullBool
	Points            sql.NullInt64
}

type UserModel struct {
	UUID           uuid.UUID `json:"uuid"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedPassword"`
	Salt           string    `json:"salt"`
	Role           RoleType  `json:"role"`
}

type UserCarrier interface {
	GetInfo() *UserInfoModel
}

type UserInfoModel struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Role     RoleType  `json:"role"`
}

type RecruiterModel struct {
	*UserInfoModel
	OrganisationName  string
	ApplicationStatus bool
}

func (model *RecruiterModel) GetInfo() *UserInfoModel {
	return model.UserInfoModel
}

type StudentModel struct {
	*UserInfoModel
	Points int64
}

func (model *StudentModel) GetInfo() *UserInfoModel {
	return model.UserInfoModel
}

type RoleType int

const (
	Student RoleType = iota
	Recruiter
	Admin
)

func (roleType *RoleType) String() string {
	return [...]string{"Student", "Recruiter", "Admin"}[*roleType]
}

func ParseRoleType(s string) (RoleType, error) {
	switch s {

	case "Student":
		return Student, nil
	case "Recruiter":
		return Recruiter, nil
	case "Admin":
		return Admin, nil
	default:
		return -1, fmt.Errorf("invalid role type: %s", s)

	}
}

func (roleType *RoleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(roleType.String())
}

func (roleType *RoleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseRoleType(s)
	if err != nil {
		return err
	}

	*roleType = parsed
	return nil
}

type GetUserRequest struct {
	UserUUID string
	Username string
}

type GetUserResponse struct {
	Success bool
	Message string
	Data    *UserInfoModel
}
