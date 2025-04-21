package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type OpportunityModel struct {
	UUID            uuid.UUID `json:"uuid"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Points          int64     `json:"points"`
	Location        string    `json:"location"`
	OpportunityType string    `json:"opportunityType"`
	PostedByUUID    uuid.UUID `json:"postedByUUID"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	Tags  *[]TagModel   `json:"tags"`
	Media *[]MediaModel `json:"media"`
}

type TagModel struct {
	ID      int64  `json:"id"`
	TagName string `json:"tagName"`
}

type MediaModel struct {
	Type MediaType `json:"type"`
	URL  string    `json:"URL"`
}

type MediaType int

const (
	Image MediaType = iota
	Video
	Text
)

func (mediaType *MediaType) String() string {
	return [...]string{"Image", "Video", "Text"}[*mediaType]
}

func ParseMediaType(s string) (MediaType, error) {
	switch s {

	case "Image":
		return Image, nil
	case "Video":
		return Video, nil
	case "Text":
		return Text, nil
	default:
		return -1, fmt.Errorf("invalid media type: %s", s)

	}
}

func (mediaType *MediaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mediaType.String())
}

func (mediaType *MediaType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseMediaType(s)
	if err != nil {
		return err
	}

	*mediaType = parsed
	return nil
}

type CreateOpportunityStatus struct {
	*OpportunityModel
	Success bool
	Message string
}

type CreateOpportunityRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Type        string   `json:"type"`
	AuthorUUID  string   `json:"author"`
	Points      int64    `json:"points"`
	Tags        []string `json:"tags"`
	MediaTypes  []string `json:"mediaType"`
	MediaURLs   []string `json:"mediaURL"`
}
