package repositories

import (
	"backend/internal/db/adapters/mysql"
	"backend/internal/models"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

//TODO: optimise fetching process. Batch calls needed
// (20/4) pretty much done now
// (21/4) Need to speed up insertion. Takes too long

//Represents opportunities

//Use auto_increment id for easy lookups when doing index based searching. I.e look for posts between ID=12 and ID=20

const CreateOpportunityPostedByIndex = "CREATE INDEX idx_opportunity_posted_by ON OpportunitiesTable(postedByUUID);"

const InsertOpportunityQuery = `
INSERT INTO OpportunitiesTable (
    uuid,
    title,
    description,
	points,
    location,
    opportunityType,
    postedByUUID
) VALUES (?, ?, ?, ?, ?, ?, ?);
`

const UpdateOpportunityQuery = `
UPDATE OpportunitiesTable
SET title = ?, description = ?, points = ?, location = ?, opportunityType = ?
WHERE uuid = ?`

const UpdateApproveOpporunityQuery = `
UPDATE OpportunitiesTable 
SET approved = ? 
WHERE uuid = ?
`

// TODO Extract things i need don't want all
const GetOpportunityByIDQuery = `
SELECT * FROM OpportunitiesTable
LEFT JOIN OpportunityMediaTable 
  ON OpportunitiesTable.uuid = OpportunityMediaTable.opportunityUUID
LEFT JOIN OpportunityTagsTable 
  ON OpportunitiesTable.uuid = OpportunityTagsTable.opportunityUUID
LEFT JOIN TagsTable 
  ON OpportunityTagsTable.tagID = TagsTable.id
WHERE OpportunitiesTable.uuid IN (%s) AND OpportunitiesTable.approved = TRUE;
`

const GetOpportunityByAuthorIDQuery = `
SELECT * FROM OpportunitiesTable
LEFT JOIN OpportunityMediaTable 
  ON OpportunitiesTable.uuid = OpportunityMediaTable.opportunityUUID
LEFT JOIN OpportunityTagsTable 
  ON OpportunitiesTable.uuid = OpportunityTagsTable.opportunityUUID
LEFT JOIN TagsTable 
  ON OpportunityTagsTable.tagID = TagsTable.id
WHERE OpportunitiesTable.postedByUUID IN (%s);
`

const GetOpportunityByFromQuery = `
WITH limited_opportunities AS (
  SELECT * FROM OpportunitiesTable
  WHERE id > ? -- First parameter: from_id
  AND uuid NOT IN (
    -- Filter out opportunities the user has liked
    SELECT opportunityUUID FROM OpportunityLikesTable WHERE userUUID = ?
    UNION
    -- Filter out opportunities the user has disliked
    SELECT opportunityUUID FROM OpportunityDislikesTable WHERE userUUID = ?
  )
  AND approved = TRUE
  ORDER BY id ASC
  LIMIT ? -- Second parameter: limit
)
SELECT lo.*, 
       omt.*, 
       ott.*, 
       tt.*
FROM limited_opportunities lo
LEFT JOIN OpportunityMediaTable omt 
  ON lo.uuid = omt.opportunityUUID
LEFT JOIN OpportunityTagsTable ott 
  ON lo.uuid = ott.opportunityUUID
LEFT JOIN TagsTable tt 
  ON ott.tagID = tt.id;
`

const DeleteOpportunityQuery = "DELETE FROM OpportunitiesTable WHERE uuid = ?"

// Tracks opportunities a user has liked

const AddOpportunityLikeQuery = `
INSERT INTO OpportunityLikesTable(userUUID, opportunityUUID) VALUES (?,?)
`

const RemoveOpportunityLikesQuery = `
DELETE FROM OpportunityLikesTable 
WHERE userUUID = ? AND opportunityUUID IN (?)
`

const GetOpportunityLikesByUserIDQuery = `
SELECT 
	olt.opportunityUUID,ot.title, ot.description, ot.location, 
	ot.opportunityType, ot.postedByUUID, omt.mediaURL, omt.mediaType,
	ott.tagID, tt.tagName
FROM OpportunityLikesTable olt
INNER JOIN OpportunitiesTable ot 
 ON ot.uuid = olt.opportunityUUID
LEFT JOIN OpportunityMediaTable omt
  ON ot.uuid = omt.opportunityUUID
LEFT JOIN OpportunityTagsTable ott
  ON ot.uuid = ott.opportunityUUID
LEFT JOIN TagsTable tt
  ON ott.tagID = tt.id
WHERE OpportunityLikesTable.userUUID = ?
`

const GetOpportunityLikesByOpportunityIDQuery = `
SELECT 
    olt.userUUID, ut.username, ut.email, st.description, st.profile, ot.id
FROM OpportunityLikesTable olt
INNER JOIN UserTable ut
    ON ut.uuid = olt.userUUID
INNER JOIN StudentInfoTable st
    ON st.uuid = ut.uuid
INNER JOIN OpportunitiesTable ot
    ON ot.uuid = olt.opportunityUUID
WHERE ot.id > ?
  AND olt.opportunityUUID = ?
LIMIT ?;
`

const CreateOpportunityLikedIndex = "CREATE INDEX idx_like_user ON OpportunityLikesTable(userUUID);"

// Tracks opportunities a user has disliked

const AddOpportunityDisLikeQuery = `
INSERT INTO OpportunityDislikesTable(userUUID, opportunityUUID) VALUES (?,?)
`

const RemoveOpportunityDisLikesQuery = `
DELETE FROM OpportunityDislikesTable 
WHERE userUUID = ? AND opportunityUUID IN (?)
`

const GetOpportunityDisLikesQuery = `
SELECT 
	olt.opportunityUUID,ot.title, ot.description, ot.location, 
	ot.opportunityType, ot.postedByUUID, omt.mediaURL, omt.mediaType,
	ott.tagID, tt.tagName
FROM OpportunityDislikesTable olt
INNER JOIN OpportunitiesTable ot 
 ON ot.uuid = olt.opportunityUUID
LEFT JOIN OpportunityMediaTable omt
  ON ot.uuid = omt.opportunityUUID
LEFT JOIN OpportunityTagsTable ott
  ON ot.uuid = ott.opportunityUUID
LEFT JOIN TagsTable tt
  ON ott.tagID = tt.id
WHERE OpportunityDislikesTable.userUUID = ?
`

const CreateOpportunityDislikedIndex = "CREATE INDEX idx_dislike_user ON OpportunityDislikesTable(userUUID);"

const CreateTagIndex = "CREATE INDEX idx_tag_name ON TagsTable(tagName);"

const GetTagByName = "SELECT * FROM TagsTable WHERE tagName = %s"
const CreateTagQuery = `INSERT INTO TagsTable (tagName) VALUES (?)`

const GetTagByID = "SELECT tagName from TagsTable WHERE id = ?"
const DeleteTagQuery = "DELETE FROM TagsTable WHERE tagName = ?" //Inappropriate tag remove etc -> Not linked to posts!

// Opportunity tags

const CreateOpportunityTagsQuery = "INSERT INTO OpportunityTagsTable (opportunityUUID, tagID) VALUES (?,?)"
const GetOpportunityTagsQuery = `
SELECT tagID FROM OpportunityTagsTable WHERE opportunityUUID = ?
`
const DeleteOpportunityTagsQuery = `
DELETE FROM OpportunityTagsTable WHERE opportunityUUID = ?
`

// For tag filtering
const GetOpportunitiesByTagQuery = "SELECT opportunityUUID FROM OpportunityTagsTable WHERE tagID = ?"

// Attatching images to opportunity

const InsertOpportunityMediaQuery = "INSERT IGNORE INTO OpportunityMediaTable (opportunityUUID, mediaURL, mediaType) VALUES (?,?,?)"
const GetOpportunityMediaQuery = "SELECT mediaURL, mediaType FROM OpportunityMediaTable WHERE opportunityUUID = ?"
const DeleteOpportunityMediaQuery = `
DELETE FROM OpportunityMediaTable WHERE opportunityUUID = ?
`

type OpportunityRepository struct {
	*BaseRepository
}

func NewOpportunityRepository(db *mysql.Repository) (*OpportunityRepository, error) {
	ur := &OpportunityRepository{}
	baseRepo, err := InitRepository(ur, db)

	if err != nil {
		return nil, err
	}
	ur.BaseRepository = baseRepo
	return ur, nil
}

func (_ *OpportunityRepository) CreateTablesQuery() *[]string {
	var queries []string
	return &queries //Was for in code creation
}

// CreateIndexesQuery returns a list of SQL queries needed to create necessary indexes for user management
func (_ *OpportunityRepository) CreateIndexesQuery() *[]string {
	return &[]string{CreateOpportunityPostedByIndex, CreateOpportunityLikedIndex, CreateOpportunityDislikedIndex, CreateTagIndex}
}

func (repo *OpportunityRepository) UpdateOpportunityStatus(opportunityUUID uuid.UUID, status bool) error {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewBoolColumn("approved", status),
		mysql.NewUUIDColumn("uuid", opportunityUUID),
	}

	_, err := container.ExecuteQuery(UpdateApproveOpporunityQuery, columns, mysql.QueryOptions{})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil

}

func (repo *OpportunityRepository) CreateOpportunity(model *models.OpportunityModel) error {

	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", model.UUID),
		mysql.NewVarcharColumn("title", model.Title),
		mysql.NewVarcharColumn("description", model.Description),
		mysql.NewIntegerColumn("points", model.Points),
		mysql.NewVarcharColumn("location", model.Location),
		mysql.NewVarcharColumn("opportunityType", model.OpportunityType),
		mysql.NewUUIDColumn("postedByUUID", model.PostedByUUID),
	}

	transaction, err := container.StartTransaction()
	if err != nil {
		log.Error(err)
		return err
	}

	defer transaction.Rollback()

	_, err = container.AddExecuteTransaction(transaction, InsertOpportunityQuery, columns)
	if err != nil {
		log.Error(err)
		return err
	}

	err = handlePostTags(transaction, model, container)
	if err != nil {
		return err
	}

	err = handlePostMedia(transaction, model, container)
	if err != nil {
		return err
	}

	return container.CommitTransaction(transaction)

}

func (repo *OpportunityRepository) UpdateOpportunity(model *models.OpportunityModel) error {
	container := repo.Repository

	uuidColumn := mysql.NewUUIDColumn("uuid", model.UUID)

	columns := []mysql.Column{
		mysql.NewVarcharColumn("title", model.Title),
		mysql.NewVarcharColumn("description", model.Description),
		mysql.NewIntegerColumn("points", model.Points),
		mysql.NewVarcharColumn("location", model.Location),
		mysql.NewVarcharColumn("opportunityType", model.OpportunityType),
		uuidColumn,
	}

	transaction, err := container.StartTransaction()
	if err != nil {
		log.Error(err)
		return err
	}

	defer transaction.Rollback()
	_, err = container.AddExecuteTransaction(transaction, UpdateOpportunityQuery, columns)
	if err != nil {
		log.Error(err)
		return err
	}

	//Delete old tags

	_, err = container.AddExecuteTransaction(transaction, DeleteOpportunityTagsQuery, []mysql.Column{uuidColumn})
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = container.AddExecuteTransaction(transaction, DeleteOpportunityMediaQuery, []mysql.Column{uuidColumn})
	if err != nil {
		log.Error(err)
		return err
	}

	//Handle new shit
	err = handlePostTags(transaction, model, container)
	if err != nil {
		return err
	}

	err = handlePostMedia(transaction, model, container)
	if err != nil {
		return err
	}

	return container.CommitTransaction(transaction)
}

func handlePostMedia(transaction *sql.Tx, model *models.OpportunityModel, container *mysql.Repository) error {
	images := *model.Media

	if len(images) < 1 {
		return nil
	}

	for _, image := range images {

		columns := []mysql.Column{
			mysql.NewUUIDColumn("opportunityUUID", model.UUID),
			mysql.NewTextColumn("mediaURL", image.URL),
			mysql.NewVarcharColumn("mediaType", image.Type.String()),
		}

		_, err := container.AddExecuteTransaction(transaction, InsertOpportunityMediaQuery, columns)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil

}

func handlePostTags(transaction *sql.Tx, postModel *models.OpportunityModel, container *mysql.Repository) error {
	tags := *postModel.Tags
	var columns []mysql.Column
	if len(tags) < 1 {
		return nil
	}

	//Insert tags
	for _, tag := range tags {

		tagName := tag.TagName

		tag, err := GetTagModelByName(tagName, container, true)
		if err != nil {
			return err
		}
		//Insert tagID and postID into table

		columns = []mysql.Column{
			mysql.NewUUIDColumn("postUUID", postModel.UUID),
			mysql.NewIntegerColumn("tagID", tag.ID),
		}

		_, err = container.AddExecuteTransaction(transaction, CreateOpportunityTagsQuery, columns)
		if err != nil {
			log.Error(err)
			return err
		}

	}
	return nil
}

func (repo *OpportunityRepository) GetOpportunity(opportunityUUIDs ...*uuid.UUID) (*[]models.OpportunityModel, error) {

	if len(opportunityUUIDs) == 0 {
		return nil, nil
	}

	container := repo.Repository

	var columns []mysql.Column

	placeholders := strings.Repeat("?,", len(opportunityUUIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(GetOpportunityByIDQuery, placeholders)

	for _, uID := range opportunityUUIDs {
		columns = append(columns, mysql.NewUUIDColumn("uuid", *uID))
	}

	rows, err := container.ExecuteQuery(query, columns, mysql.QueryOptions{})
	if err != nil {
		return nil, err
	}
	opportunity, _, err := getOpportunity(rows)
	if err != nil {
		return nil, err
	}
	if opportunity == nil {
		return nil, nil
	}

	return opportunity, nil
}

func (repo *OpportunityRepository) GetOpportunityByAuthor(authorUUIDs ...*uuid.UUID) (*[]models.OpportunityModel, error) {

	if len(authorUUIDs) == 0 {
		return nil, nil
	}

	container := repo.Repository

	var columns []mysql.Column

	placeholders := strings.Repeat("?,", len(authorUUIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(GetOpportunityByAuthorIDQuery, placeholders)

	for _, uID := range authorUUIDs {
		columns = append(columns, mysql.NewUUIDColumn("uuid", *uID))
	}

	rows, err := container.ExecuteQuery(query, columns, mysql.QueryOptions{})
	if err != nil {
		return nil, err
	}
	opportunity, _, err := getOpportunity(rows)
	if err != nil {
		return nil, err
	}
	if opportunity == nil {
		return nil, nil
	}

	return opportunity, nil
}

func (repo *OpportunityRepository) GetOpportunitiesByTag(tagName string) (*[]models.OpportunityModel, error) {

	container := repo.Repository

	tag, err := GetTagModelByName(tagName, container, false)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, nil
	}

	columns := []mysql.Column{
		mysql.NewIntegerColumn("tagID", tag.ID),
	}

	rows, err := container.ExecuteQuery(GetOpportunitiesByTagQuery, columns, mysql.QueryOptions{})

	var opportunities []*uuid.UUID

	for rows.Next() {

		var opportunityUUID uuid.UUID

		err := rows.Scan(&opportunityUUID)
		if err != nil {
			return nil, err
		}

		opportunities = append(opportunities, &opportunityUUID)

	}

	opportunity, err := repo.GetOpportunity(opportunities...)
	if err != nil {
		return nil, err
	}

	return opportunity, nil

}

func (repo *OpportunityRepository) GetOpportunitiesFrom(from int64, limit int64, userUUID uuid.UUID) (*[]models.OpportunityModel, int64, error) {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewIntegerColumn("from", from),
		mysql.NewUUIDColumn("uuid", userUUID),
		mysql.NewUUIDColumn("uuid", userUUID),
		mysql.NewIntegerColumn("limit", limit),
	}
	rows, err := container.ExecuteQuery(GetOpportunityByFromQuery, columns, mysql.QueryOptions{})

	if err != nil {
		return nil, 0, err
	}

	return getOpportunity(rows)
}

func (repo *OpportunityRepository) LikeOpportunity(userUUID, opportunityUUID uuid.UUID) error {
	return repo.insertAction(userUUID, opportunityUUID, AddOpportunityLikeQuery)
}

func (repo *OpportunityRepository) DeleteLikeOpportunity(userUUID, opportunityUUID uuid.UUID) error {
	return repo.deleteAction(userUUID, opportunityUUID, RemoveOpportunityLikesQuery)
}

func (repo *OpportunityRepository) DislikeOpportunity(userUUID, opportunityUUID uuid.UUID) error {
	return repo.insertAction(userUUID, opportunityUUID, AddOpportunityDisLikeQuery)
}

func (repo *OpportunityRepository) DeleteDislikeOpportunity(userUUID, opportunityUUID uuid.UUID) error {
	return repo.deleteAction(userUUID, opportunityUUID, RemoveOpportunityDisLikesQuery)
}
func (repo *OpportunityRepository) insertAction(userUUID, opportunityUUID uuid.UUID, query string) error {
	columns := buildUUIDColumnsAction(userUUID, opportunityUUID)
	_, err := repo.Repository.ExecuteInsert(query, columns, mysql.InsertOptions{})
	return err
}

func (repo *OpportunityRepository) deleteAction(userUUID, opportunityUUID uuid.UUID, query string) error {
	columns := buildUUIDColumnsAction(userUUID, opportunityUUID)
	_, err := repo.Repository.ExecuteQuery(query, columns, mysql.QueryOptions{})
	return err
}

func buildUUIDColumnsAction(userUUID, opportunityUUID uuid.UUID) []mysql.Column {
	return []mysql.Column{
		mysql.NewUUIDColumn("userUUID", userUUID),
		mysql.NewUUIDColumn("opportunityUUID", opportunityUUID),
	}
}

func getOpportunity(rows *sql.Rows) (*[]models.OpportunityModel, int64, error) {
	type seenData struct {
		media map[string]bool // use media URL as unique key
		tags  map[int64]bool  // use tag ID as unique key
	}

	opportunities := map[uuid.UUID]*models.OpportunityModel{}
	seen := map[uuid.UUID]*seenData{}

	var lastIDSeen int64

	for rows.Next() {
		// Main opportunity fields
		var id int64
		var opportunityUUID uuid.UUID
		var title, description, location, opportunityType string
		var points int64
		var postedByUUID uuid.UUID
		var createdAt, updatedAt time.Time
		var approved bool

		// Media
		var mediaID sql.Null[int64]
		var mediaOpportunityUUID uuid.UUID
		var mediaURL sql.NullString
		var mediaType sql.NullString

		// Tag
		var tagID sql.NullInt64
		var tagOpportunityUUID uuid.UUID
		var tagName sql.NullString

		err := rows.Scan(&id, &opportunityUUID, &title, &description, &points,
			&location, &opportunityType, &postedByUUID,
			&createdAt, &updatedAt, &approved,
			&mediaID, &mediaOpportunityUUID, &mediaURL, &mediaType,
			&tagOpportunityUUID, &tagID, &tagID, &tagName)

		if err != nil {
			return nil, 0, err
		}
		lastIDSeen = id

		// Initialize opportunity and tracking
		if _, exists := opportunities[opportunityUUID]; !exists {
			opportunities[opportunityUUID] = &models.OpportunityModel{
				UUID:            opportunityUUID,
				Title:           title,
				Description:     description,
				Points:          points,
				Location:        location,
				OpportunityType: opportunityType,
				PostedByUUID:    postedByUUID,
				CreatedAt:       createdAt,
				UpdatedAt:       updatedAt,
				Approved:        approved,
				Tags:            &[]models.TagModel{},
				Media:           &[]models.MediaModel{},
			}
			seen[opportunityUUID] = &seenData{
				media: map[string]bool{},
				tags:  map[int64]bool{},
			}
		}

		opportunity := opportunities[opportunityUUID]
		tracker := seen[opportunityUUID]

		// Deduplicate and add media
		if mediaURL.Valid && mediaType.Valid {
			url := mediaURL.String
			if !tracker.media[url] {
				mType, err := models.ParseMediaType(mediaType.String)
				if err == nil {
					*opportunity.Media = append(*opportunity.Media, models.MediaModel{
						Type: mType,
						URL:  url,
					})
					tracker.media[url] = true
				}
			}
		}

		// Deduplicate and add tags
		if tagName.Valid && tagID.Valid {
			id := tagID.Int64
			if !tracker.tags[id] {
				*opportunity.Tags = append(*opportunity.Tags, models.TagModel{
					ID:      id,
					TagName: tagName.String,
				})
				tracker.tags[id] = true
			}
		}
	}

	if len(opportunities) == 0 {
		return nil, 0, nil
	}

	opportunitiesSlice := make([]models.OpportunityModel, 0, len(opportunities))
	for _, opp := range opportunities {
		opportunitiesSlice = append(opportunitiesSlice, *opp)
	}

	return &opportunitiesSlice, lastIDSeen, nil
}

func (repo *OpportunityRepository) GetOpportunityTags(opportunityUUID uuid.UUID) *[]models.TagModel {

	var tags []models.TagModel

	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("opportunityUUID", opportunityUUID),
	}

	rows, err := container.ExecuteQuery(GetOpportunityTagsQuery, columns, mysql.QueryOptions{})
	if err != nil {
		log.Error(err)
		return nil
	}

	for rows.Next() {

		var tagID int64

		err := rows.Scan(&tagID)
		if err != nil {
			log.Error(err)
			return nil
		}

		//Fetch tagName

		columns := []mysql.Column{
			mysql.NewIntegerColumn("id", tagID),
		}

		rows, err := container.ExecuteQuery(GetTagByID, columns, mysql.QueryOptions{})
		if err != nil {
			log.Error(err)
			return nil
		}
		var tagName string
		if rows.Next() {
			err := rows.Scan(&tagName)
			if err != nil {
				log.Error(err)
				return nil
			}
		}

		tag := models.TagModel{
			ID:      tagID,
			TagName: tagName,
		}

		tags = append(tags, tag)
	}

	if len(tags) == 0 {
		return nil
	}

	return &tags
}

func (repo *OpportunityRepository) DeleteOpportunity(opportunityUUID uuid.UUID) error {

	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewUUIDColumn("uuid", opportunityUUID),
	}

	_, err := container.ExecuteQuery(DeleteOpportunityQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return err
	}

	return nil

}

func (repo *OpportunityRepository) GetOpportunityMedia(opportunityUUID uuid.UUID) *[]models.MediaModel {

	container := repo.Repository

	var media []models.MediaModel

	columns := []mysql.Column{
		mysql.NewUUIDColumn("opportunityUUID", opportunityUUID),
	}

	rows, err := container.ExecuteQuery(GetOpportunityMediaQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return nil
	}

	for rows.Next() {

		var mediaURL string
		var mediaType string

		err := rows.Scan(&mediaURL, &mediaType)
		if err != nil {
			log.Error(err)
			return nil
		}

		parsedMediaType, err := models.ParseMediaType(mediaType)
		if err != nil {
			log.Error(err)
			return nil
		}

		model := models.MediaModel{
			Type: parsedMediaType,
			URL:  mediaURL,
		}

		media = append(media, model)

	}

	if len(media) == 0 {
		return nil
	}

	return &media

}

func GetTagModelByName(tagName string, container *mysql.Repository, createIfNotExist bool) (*models.TagModel, error) {
	columns := []mysql.Column{
		mysql.NewVarcharColumn("tagName", tagName),
	}
	query := fmt.Sprintf(GetTagByName, "?")
	rows, err := container.ExecuteQuery(query, columns, mysql.QueryOptions{})

	if err != nil {
		return nil, err
	}

	if !rows.Next() {

		//Tag doesn't exist

		if !createIfNotExist {
			return nil, nil
		}

		columns = []mysql.Column{
			mysql.NewVarcharColumn("tagName", tagName),
		}
		var tag models.TagModel
		_, err := container.ExecuteInsert(CreateTagQuery, columns, mysql.InsertOptions{
			OnComplete: func(result sql.Result) {
				tagID, err := result.LastInsertId()
				if err != nil {
					log.Error(err)
				}

				tag = models.TagModel{
					ID:      tagID,
					TagName: tagName,
				}

			},
		})
		if err != nil {
			return nil, err
		}

		return &tag, nil
	}

	var tagID int64

	err = rows.Scan(&tagID, &tagName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	model := models.TagModel{
		ID:      tagID,
		TagName: tagName,
	}

	return &model, nil
}

func (repo *OpportunityRepository) GetOpportunityByLikes(opportunityUUID uuid.UUID, from int64, limit int64) (*[]*models.StudentInfoModel, int64, error) {
	container := repo.Repository

	columns := []mysql.Column{
		mysql.NewIntegerColumn("id", from),
		mysql.NewUUIDColumn("uuid", opportunityUUID),
		mysql.NewIntegerColumn("limit", limit),
	}

	rows, err := container.ExecuteQuery(GetOpportunityLikesByOpportunityIDQuery, columns, mysql.QueryOptions{})
	if err != nil {
		return nil, 0, err
	}

	var students []*models.StudentInfoModel
	var lastRow int64
	for rows.Next() {
		var userUUID uuid.UUID
		var username string
		var email string
		var description sql.NullString
		var profile sql.NullString
		var rowID int64

		err = rows.Scan(&userUUID, &username, &email, &description, &profile, &rowID)
		if err != nil {
			return nil, 0, err
		}

		students = append(students, &models.StudentInfoModel{
			StudentID:    userUUID,
			StudentName:  username,
			StudentEmail: email,
			Description:  description.String,
			ProfilePic:   profile.String,
		})

		lastRow = rowID

	}

	if len(students) == 0 {
		return nil, 0, nil
	}

	return &students, lastRow, nil

}
