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

const CreateStudentInfoTableQuery = `
CREATE TABLE IF NOT EXISTS StudentInfoTable(
uuid VARCHAR(36) PRIMARY KEY,
description TEXT,
profile TEXT,
FOREIGN KEY (uuid) REFERENCES UserTable(uuid) ON DELETE CASCADE)`

const CreateStudentInfoQuery = `
INSERT INTO StudentInfoTable(uuid, description, profile) VALUES (?,NULL,NULL)`

const UpdateStudentInfoQuery = `
UPDATE StudentInfoTable
SET description = ?, profile = ?
WHERE uuid = ?
`

const GetStudentInfoQuery = `
SELECT 
    ut.uuid, 
    ut.username, 
    ut.email, 
    st.description, 
    st.profile, 
    utl.tagID,
	ttl.tagName,
    utd.tagID,
	ttd.tagName
	
FROM UserTable ut
INNER JOIN StudentInfoTable st
    ON st.uuid = ut.uuid
LEFT JOIN UserTagsLiked utl
    ON utl.uuid = ut.uuid
LEFT JOIN TagsTable ttl
    ON ttl.id = utl.tagID
LEFT JOIN UserTagsDisLiked utd
    ON utd.uuid = ut.uuid
LEFT JOIN TagsTable ttd
    ON ttd.id = utd.tagID
WHERE ut.uuid = ?
`

//Student tag preferences things

const InsertStudentTagsLikedQuery = `
INSERT INTO UserTagsLiked(uuid, tagID) VALUES %s
`

const GetStudentTagsLikedQuery = `
SELECT ul.uuid, ut.username, tg.tagID, tg.tagName FROM UserTagsLiked ul
INNER JOIN TagsTable tg
	ON tg.id = ul.tagID
INNER JOIN UserTable ut
	ON ut.uuid = ul.uuid
WHERE ul.uuid = ?
`

const RemoveStudentTagsLikedQuery = `
DELETE FROM UserTagsLiked
WHERE uuid = ?
`

const CreateStudentTagsDisLikedTableQuery = `
CREATE TABLE IF NOT EXISTS UserTagsDisLiked(
	uuid VARCHAR(36), 
	tagID int,
	PRIMARY KEY (uuid, tagID),
	FOREIGN KEY (uuid) REFERENCES UserTable(uuid) ON DELETE CASCADE,
	FOREIGN KEY (tagID) REFERENCES TagsTable(id) ON DELETE CASCADE)
`

const InsertStudentTagsDisLikedQuery = `
INSERT INTO UserTagsDisLiked(uuid, tagID) VALUES %s
`

const GetUserStudentDisLikedQuery = `
SELECT ul.uuid, ut.username, tg.tagID, tg.tagName FROM UserTagsDisLiked ul
INNER JOIN TagsTable tg
	ON tg.id = ul.tagID
INNER JOIN UserTable ut
	ON ut.uuid = ul.uuid
WHERE ul.uuid = ?
`

const RemoveStudentTagsDisLikedQuery = `
DELETE FROM UserTagsDisLiked
WHERE uuid = ?
`

type StudentRepository struct {
	*BaseRepository
}

func NewStudentRepository(db *mysql.Repository) (*StudentRepository, error) {
	ur := &StudentRepository{}
	baseRepo, err := InitRepository(ur, db)

	if err != nil {
		return nil, err
	}
	ur.BaseRepository = baseRepo
	return ur, nil
}

func (_ *StudentRepository) CreateTablesQuery() *[]string {
	return &[]string{CreateStudentInfoTableQuery}
}

// CreateIndexesQuery returns a list of SQL queries needed to create necessary indexes for user management
func (_ *StudentRepository) CreateIndexesQuery() *[]string {
	return &[]string{}
}

func (repo *StudentRepository) UpdateStudentInfo(studentInfo models.StudentInfoModel) error {
	container := repo.Repository

	transaction, err := container.StartTransaction()
	if err != nil {
		return err
	}

	columns := []mysql.Column{
		mysql.NewTextColumn("description", studentInfo.Description),
		mysql.NewTextColumn("profile", studentInfo.ProfilePic),
		mysql.NewUUIDColumn("uuid", studentInfo.StudentID),
	}

	_, err = container.AddExecuteTransaction(transaction, UpdateStudentInfoQuery, columns)
	if err != nil {
		return err
	}

	uuidColumn := []mysql.Column{
		mysql.NewUUIDColumn("uuid", studentInfo.StudentID),
	}

	_, err = container.AddExecuteTransaction(transaction, RemoveStudentTagsLikedQuery, uuidColumn)
	if err != nil {
		return err
	}

	_, err = container.AddExecuteTransaction(transaction, RemoveStudentTagsDisLikedQuery, uuidColumn)
	if err != nil {
		return err
	}

	if len(studentInfo.TagsLiked) > 0 {
		err = updateStudentTagOpinion(container, transaction, studentInfo.StudentID, studentInfo.TagsLiked, InsertStudentTagsLikedQuery)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if len(studentInfo.TagsDisliked) > 0 {
		err = updateStudentTagOpinion(container, transaction, studentInfo.StudentID, studentInfo.TagsDisliked, InsertStudentTagsDisLikedQuery)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	err = container.CommitTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func updateStudentTagOpinion(container *mysql.Repository, transaction *sql.Tx, userID uuid.UUID, tags []string, query string) error {
	if len(tags) == 0 {
		return nil
	}

	uniqueTags := make(map[string]bool)
	for _, tag := range tags {
		uniqueTags[tag] = true
	}

	var allTagIDs []int64
	var allTagNames []string

	// check which tags already exist
	for tag := range uniqueTags {
		var tagColumns []mysql.Column
		tagColumns = append(tagColumns, mysql.NewVarcharColumn("tagName", tag))

		rows, err := container.ExecuteQuery("SELECT * FROM TagsTable WHERE tagName = ?", tagColumns, mysql.QueryOptions{})
		if err != nil {
			return err
		}

		found := false
		var id int64
		var name string

		for rows.Next() {
			found = true
			err = rows.Scan(&id, &name)
			if err != nil {
				rows.Close()
				return err
			}
			allTagIDs = append(allTagIDs, id)
			allTagNames = append(allTagNames, name)
		}
		rows.Close()

		// Create tag if it don't exist
		if !found {
			tagColumn := []mysql.Column{
				mysql.NewVarcharColumn("tagName", tag),
			}

			result, err := container.AddExecuteTransaction(transaction, CreateTagQuery, tagColumn)
			if err != nil {
				return err
			}

			newTagID, err := result.LastInsertId()
			if err != nil {
				return err
			}

			allTagIDs = append(allTagIDs, newTagID)
			allTagNames = append(allTagNames, tag)
		}
	}

	if len(allTagIDs) > 0 {
		var tagColumns []mysql.Column
		placeholders := make([]string, 0, len(allTagIDs))

		for _, tagID := range allTagIDs {
			placeholders = append(placeholders, "(?, ?)")
			tagColumns = append(tagColumns,
				mysql.NewUUIDColumn("uuid", userID),
				mysql.NewIntegerColumn("tagID", tagID),
			)
		}

		finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ", "))
		_, err := container.AddExecuteTransaction(transaction, finalQuery, tagColumns)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *StudentRepository) GetUserInfo(userID uuid.UUID) (*models.StudentInfoModel, error) {
	container := repo.Repository
	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", userID),
	}

	rows, err := container.ExecuteQuery(GetStudentInfoQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userInfo *models.StudentInfoModel
	var tagsLiked []string
	var tagsDisliked []string

	for rows.Next() {
		var username string
		var email string
		var description sql.NullString
		var profilePic sql.NullString
		var likeTagID sql.NullInt64
		var likeTagName sql.NullString
		var dislikeTagID sql.NullInt64
		var dislikeTagName sql.NullString

		err = rows.Scan(&userID, &username, &email, &description, &profilePic, &likeTagID, &likeTagName, &dislikeTagID, &dislikeTagName)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if userInfo == nil {
			userInfo = &models.StudentInfoModel{
				StudentID:    userID,
				StudentName:  username,
				StudentEmail: email,
				Description:  description.String,
				ProfilePic:   profilePic.String,
				TagsLiked:    []string{},
				TagsDisliked: []string{},
			}
		}

		if likeTagID.Valid && likeTagName.Valid {
			tagsLiked = append(tagsLiked, likeTagName.String)
		}

		if dislikeTagID.Valid && dislikeTagName.Valid {
			tagsDisliked = append(tagsDisliked, dislikeTagName.String)
		}
	}

	if userInfo == nil {
		return nil, nil
	}

	userInfo.TagsLiked = tagsLiked
	userInfo.TagsDisliked = tagsDisliked

	return userInfo, nil
}
