package repositories

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

//Repository used for stored points

//From what im aware points are given by employer -> set between 5 - 20 bonus
//Default amount of points given by system depending on opportunity type 10,15,20,25 etc
//Points can be redeemed for prizes or smt?

//Literally only storing user points on here

const IncrementPointsQuery = `
INSERT INTO UserPointsTable (uuid, points)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE points = points + VALUES(points);

`

const DecrementPointsQuery = `
INSERT INTO UserPointsTable (uuid, points)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE points = GREATEST(points - VALUES(points), 0);
`

const GetPointsQuery = `
SELECT ut.uuid, ut.username, ut.role, upt.points FROM UserPointsTable upt
	INNER JOIN UserTable ut
	ON ut.uuid = upt.uuid
WHERE upt.uuid IN (%s)
`

type PointsRepository struct {
	*BaseRepository
}

func NewPointsRepository(db *mysql.Repository) (*PointsRepository, error) {
	ur := &PointsRepository{}
	baseRepo, err := InitRepository(ur, db)

	if err != nil {
		return nil, err
	}
	ur.BaseRepository = baseRepo
	return ur, nil
}

func (_ *PointsRepository) CreateTablesQuery() *[]string {
	queries := []string{}
	return &queries
}

func (_ *PointsRepository) CreateIndexesQuery() *[]string {
	return &[]string{}
}

func (repo *PointsRepository) IncrementPoints(userUUID uuid.UUID, points int64) error {

	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", userUUID),
		mysql.NewIntegerColumn("points", points),
	}

	_, err := container.ExecuteQuery(IncrementPointsQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return err
	}

	return nil

}

func (repo *PointsRepository) DecrementPoints(userUUID uuid.UUID, points int64) error {

	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", userUUID),
		mysql.NewIntegerColumn("points", points),
	}

	_, err := container.ExecuteQuery(IncrementPointsQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return err
	}

	return nil

}

func (repo *PointsRepository) GetPoints(userUUIDs ...uuid.UUID) (*[]models.StudentModel, error) {

	container := repo.Repository

	placeholders := strings.Repeat("?,", len(userUUIDs))
	placeholders = placeholders[:len(placeholders)-1]

	var columns []mysql.Column

	for _, uID := range userUUIDs {

		uuidColumn := mysql.NewUUIDColumn("uuid", uID)
		columns = append(columns, uuidColumn)

	}

	rows, err := container.ExecuteQuery(GetPointsQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return nil, err
	}

	var pointModels []models.StudentModel

	for rows.Next() {
		var userUUID uuid.UUID
		var username string
		var roleName string
		var userPoints int64

		err := rows.Scan(&userUUID, &username, &roleName, &userPoints)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		roleType, err := models.ParseRoleType(roleName)
		if err != nil {
			return nil, err
		}

		model := models.StudentModel{
			UserInfoModel: &models.UserInfoModel{
				UUID:     userUUID,
				Username: username,
				Role:     roleType,
			},
			Points: userPoints,
		}

		pointModels = append(pointModels, model)

	}

	return &pointModels, nil
}
