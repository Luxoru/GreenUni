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
// (18/4) pretty much done now
// (18/4) Need to speed up insertion. Takes too long

//Represents opportunities

//Use auto_increment id for easy lookups when doing index based searching. I.e look for posts between ID=12 and ID=20

const CreateOpportunityTable = `
CREATE TABLE IF NOT EXISTS OpportunitiesTable (
	id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
	points INT NOT NULL,
    location VARCHAR(100),
    opportunityType ENUM('event', 'volunteer', 'job', 'issue') NOT NULL,
    postedByUUID VARCHAR(36) NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (postedByUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE
);
`

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

// TODO Extract things i need don't want all
const GetOpportunityByIDQuery = `
SELECT * FROM OpportunitiesTable
LEFT JOIN OpportunityMediaTable 
  ON OpportunitiesTable.uuid = OpportunityMediaTable.opportunityUUID
LEFT JOIN OpportunityTagsTable 
  ON OpportunitiesTable.uuid = OpportunityTagsTable.opportunityUUID
LEFT JOIN TagsTable 
  ON OpportunityTagsTable.tagID = TagsTable.id
WHERE OpportunitiesTable.uuid IN (%s);
`

const DeleteOpportunityQuery = "DELETE FROM OpportunitiesTable WHERE uuid = ?"

// Tracks opportunities a user has liked
const CreateOpportunityLikesTable = `
CREATE TABLE IF NOT EXISTS OpportunityLikesTable (
    userUUID VARCHAR(36),
    opportunityUUID VARCHAR(36),
    likedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userUUID, opportunityUUID),
    FOREIGN KEY (userUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);
`

const AddOpportunityLikeQuery = `
INSERT INTO OpportunityLikesTable(userUUID, opportunityUUID) VALUES %s
`

const GetOpportunityLikesQuery = `
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

const RemoveOpportunityLikesQuery = `
DELETE FROM OpportunityLikesTable 
WHERE userUUID = ? AND opportunityUUID IN (%s)
`

const CreateOpportunityLikedIndex = "CREATE INDEX idx_like_user ON OpportunityLikesTable(userUUID);"

// Tracks opportunities a user has disliked
const CreateOpportunityDislikesTable = `
CREATE TABLE IF NOT EXISTS OpportunityDislikesTable (
    userUUID VARCHAR(36),
    opportunityUUID VARCHAR(36),
    dislikedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userUUID, opportunityUUID),
    FOREIGN KEY (userUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);
`

const AddOpportunityDisLikeQuery = `
INSERT INTO OpportunityDislikesTable(userUUID, opportunityUUID) VALUES %s
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

const RemoveOpportunityDisLikesQuery = `
DELETE FROM OpportunityDislikesTable 
WHERE userUUID = ? AND opportunityUUID IN (%s)
`

const CreateOpportunityDislikedIndex = "CREATE INDEX idx_dislike_user ON OpportunityDislikesTable(userUUID);"

const CreateTagsTable = `
CREATE TABLE IF NOT EXISTS TagsTable (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tagName VARCHAR(50) UNIQUE NOT NULL
);
`

const CreateTagIndex = "CREATE INDEX idx_tag_name ON TagsTable(tagName);"

const GetTagByName = "SELECT * FROM TagsTable WHERE tagName = ?"
const CreateTagQuery = `INSERT INTO TagsTable (tagName) VALUES (?)`

const GetTagByID = "SELECT tagName from TagsTable WHERE id = ?"
const DeleteTagQuery = "DELETE FROM TagsTable WHERE tagName = ?" //Inappropriate tag remove etc -> Not linked to posts!

// Opportunity tags
const CreateOpportunityTagsTable = `
CREATE TABLE IF NOT EXISTS OpportunityTagsTable (
    opportunityUUID VARCHAR(36),
    tagID INT,
    PRIMARY KEY (opportunityUUID, tagID),
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (tagID) REFERENCES TagsTable(id) ON DELETE CASCADE
);
`

const CreateOpportunityTagsQuery = "INSERT INTO OpportunityTagsTable (opportunityUUID, tagID) VALUES (?,?)"
const GetOpportunityTagsQuery = `
SELECT tagID FROM OpportunityTagsTable WHERE opportunityUUID = ?
`

// For tag filtering
const GetOpportunitiesByTagQuery = "SELECT opportunityUUID FROM OpportunityTagsTable WHERE tagID = ?"

// Attatching images to opportunity
const CreateOpportunityMediaTable = `
CREATE TABLE IF NOT EXISTS OpportunityMediaTable (
    id INT AUTO_INCREMENT PRIMARY KEY,
    opportunityUUID VARCHAR(36),
    mediaURL TEXT NOT NULL,
	mediaType VARCHAR(50) NOT NULL,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);
`

const InsertOpportunityMediaQuery = "INSERT IGNORE INTO OpportunityMediaTable (opportunityUUID, mediaURL, mediaType) VALUES (?,?,?)"
const GetOpportunityMediaQuery = "SELECT mediaURL, mediaType FROM OpportunityMediaTable WHERE opportunityUUID = ?"

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
	queries := []string{CreateOpportunityTable,
		CreateOpportunityLikesTable,
		CreateOpportunityDislikesTable,
		CreateTagsTable,
		CreateOpportunityTagsTable,
		CreateOpportunityMediaTable}
	return &queries
}

// CreateIndexesQuery returns a list of SQL queries needed to create necessary indexes for user management
func (_ *OpportunityRepository) CreateIndexesQuery() *[]string {
	return &[]string{CreateOpportunityPostedByIndex, CreateOpportunityLikedIndex, CreateOpportunityDislikedIndex, CreateTagIndex}
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

		tag, err := getTagByName(tagName, container, true)
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
	opportunity, err := getOpportunity(rows)
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

	tag, err := getTagByName(tagName, container, false)
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

func getOpportunity(rows *sql.Rows) (*[]models.OpportunityModel, error) {
	type seenData struct {
		media map[string]bool // use media URL as unique key
		tags  map[int64]bool  // use tag ID as unique key
	}

	opportunities := map[uuid.UUID]*models.OpportunityModel{}
	seen := map[uuid.UUID]*seenData{}

	for rows.Next() {
		// Main opportunity fields
		var id int64
		var opportunityUUID uuid.UUID
		var title, description, location, opportunityType string
		var points int64
		var postedByUUID uuid.UUID
		var createdAt, updatedAt time.Time

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
			&createdAt, &updatedAt,
			&mediaID, &mediaOpportunityUUID, &mediaURL, &mediaType,
			&tagOpportunityUUID, &tagID, &tagID, &tagName)

		if err != nil {
			return nil, err
		}

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
		return nil, nil
	}

	opportunitiesSlice := make([]models.OpportunityModel, 0, len(opportunities))
	for _, opp := range opportunities {
		opportunitiesSlice = append(opportunitiesSlice, *opp)
	}

	return &opportunitiesSlice, nil
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

func getTagByName(tagName string, container *mysql.Repository, createIfNotExist bool) (*models.TagModel, error) {
	columns := []mysql.Column{
		mysql.NewVarcharColumn("tagName", tagName),
	}

	rows, err := container.ExecuteQuery(GetTagByName, columns, mysql.QueryOptions{})

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
