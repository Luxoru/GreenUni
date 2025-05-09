package repositories

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/models"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

// TODO: make async
// TODO: add profile updating, etc.

// UserRepository is used to manage user-related data in the database
// It stores and retrieves basic user information like username and password
type UserRepository struct {
	*BaseRepository
}

const AddUserToTableQuery = "INSERT INTO UserTable (uuid,username, email, hashed_pass, salt, role) VALUES (?, ?,?, ?, ?, ?);"

// When we fetch user from db, return everything (all tables here). And format to struct that we need //TODO: change this approach?
//Join on all acc type tables e.g Student and Recruiter. At least One table will always return null and deal with it like that

const GetUserByUUIDFromTableQuery = `
SELECT ut.uuid, ut.username, ut.email, ut.hashed_pass, ut.salt, ut.role, rt.organisationName, rt.applicationStatus, st.points FROM UserTable ut
LEFT JOIN RecruiterTable rt
    ON rt.uuid = ut.uuid
LEFT JOIN StudentTable st
    ON st.uuid = ut.uuid
WHERE ut.uuid = %s`

const GetUserByUsernameFromTableQuery = `
SELECT ut.uuid, ut.username, ut.email, ut.hashed_pass, ut.salt, ut.role, rt.organisationName, rt.applicationStatus, st.points FROM UserTable ut
LEFT JOIN RecruiterTable rt
    ON rt.uuid = ut.uuid
LEFT JOIN StudentTable st
    ON st.uuid = ut.uuid
WHERE ut.username = ?`

const GetUserByEmailFromTableQuery = `
SELECT ut.uuid, ut.username, ut.email, ut.hashed_pass, ut.salt, ut.role, rt.organisationName, rt.applicationStatus, st.points FROM UserTable ut
LEFT JOIN RecruiterTable rt
    ON rt.uuid = ut.uuid
LEFT JOIN StudentTable st
    ON st.uuid = ut.uuid
WHERE ut.email = ?`

//Student table -> See points.go

//Recruiter table -> On register account portal can register as student or recruiter. When registered initially status set as false
//Admin manually verifies recruiter and then can log in

const AddNewRecruiterQuery = `
INSERT INTO RecruiterTable (uuid, organisationName) VALUES (?,?)
`

const UpdateRecruiterStatus = `
UPDATE RecruiterTable
SET applicationStatus = ?
WHERE uuid = ?
`

// NewUserRepository initializes a new UserRepository instance
func NewUserRepository(db *mysql.Repository) (*UserRepository, error) {
	ur := &UserRepository{}
	baseRepo, err := InitRepository(ur, db)

	if err != nil {
		return nil, err
	}
	ur.BaseRepository = baseRepo
	return ur, nil
}

// CreateTablesQuery returns a list of SQL queries needed to create necessary tables for user management
func (_ *UserRepository) CreateTablesQuery() *[]string {
	return &[]string{} //Was for in code creation
}

// CreateIndexesQuery returns a list of SQL queries needed to create necessary indexes for user management
func (_ *UserRepository) CreateIndexesQuery() *[]string {
	return &[]string{}
}

// AddUser inserts a new user into the UserTable
func (repo *UserRepository) AddUser(userModel *models.UserModel, options mysql.InsertOptions) error {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", userModel.UUID),
		mysql.NewVarcharColumn("username", userModel.Username),
		mysql.NewVarcharColumn("email", userModel.Email),
		mysql.NewVarcharColumn("hashed_pass", userModel.HashedPassword),
		mysql.NewVarcharColumn("salt", userModel.Salt),
		mysql.NewVarcharColumn("role", userModel.Role.String()),
	}

	_, err := container.ExecuteInsert(AddUserToTableQuery, columns, options)
	if err != nil {
		return err
	}
	if userModel.Role != models.Student {
		return nil
	}

	columns = []mysql.Column{
		mysql.NewUUIDColumn("uuid", userModel.UUID),
	}
	_, err = container.ExecuteInsert(CreateStudentInfoQuery, columns, mysql.InsertOptions{})

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by UUID from the database
func (repo *UserRepository) GetUserByID(userUUID ...uuid.UUID) (*[]models.RawUserRow, error) {
	container := repo.Repository
	placeholders := strings.Repeat("?,", len(userUUID))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(GetUserByUUIDFromTableQuery, placeholders)
	var columns []mysql.Column
	for _, uID := range userUUID {
		columns = append(columns, mysql.NewUUIDColumn("uuid", uID))
	}

	rows, err := container.ExecuteQuery(query, columns, mysql.QueryOptions{})
	defer rows.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var userRows []models.RawUserRow

	for rows.Next() {
		var uid uuid.UUID
		var username string
		var email string
		var hashedPassword string
		var salt string
		var roleName string

		//Recruiter related things
		var organisationName sql.NullString
		var applicationStatus sql.NullBool
		var points sql.NullInt64

		err := rows.Scan(&uid, &username, &email, &hashedPassword, &salt, &roleName, &organisationName, &applicationStatus, &points)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		roleType, err := models.ParseRoleType(roleName)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		user := models.RawUserRow{
			UUID:              uid,
			Username:          username,
			Email:             email,
			HashedPassword:    hashedPassword,
			Salt:              salt,
			Role:              roleType,
			OrganisationName:  organisationName,
			ApplicationStatus: applicationStatus,
			Points:            points,
		}

		userRows = append(userRows, user)
	}

	if len(userRows) == 0 {
		return nil, nil
	}

	return &userRows, nil
}

// GetUserByName retrieves a user by username from the database
func (repo *UserRepository) GetUserByName(username string, options mysql.QueryOptions) (*models.RawUserRow, error) {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewVarcharColumn("username", username),
	}

	rows, err := container.ExecuteQuery(GetUserByUsernameFromTableQuery, columns, options)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return getUser(rows)
}

// GetUserByEmail retrieves a user by email from the database
func (repo *UserRepository) GetUserByEmail(email string, options mysql.QueryOptions) (*models.RawUserRow, error) {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewVarcharColumn("email", email),
	}

	rows, err := container.ExecuteQuery(GetUserByEmailFromTableQuery, columns, options)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return getUser(rows)
}

func getUser(rows *sql.Rows) (*models.RawUserRow, error) {
	for rows.Next() {
		var uid uuid.UUID
		var username string
		var email string
		var hashedPassword string
		var salt string
		var roleName string

		//Recruiter related things
		var organisationName sql.NullString
		var applicationStatus sql.NullBool
		var points sql.NullInt64

		err := rows.Scan(&uid, &username, &email, &hashedPassword, &salt, &roleName, &organisationName, &applicationStatus, &points)
		if err != nil {
			return nil, err
		}

		roleType, err := models.ParseRoleType(roleName)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		user := models.RawUserRow{
			UUID:              uid,
			Username:          username,
			Email:             email,
			HashedPassword:    hashedPassword,
			Salt:              salt,
			Role:              roleType,
			OrganisationName:  organisationName,
			ApplicationStatus: applicationStatus,
			Points:            points,
		}

		return &user, nil
	}
	return nil, nil
}
