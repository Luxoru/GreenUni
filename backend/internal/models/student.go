package models

import "github.com/google/uuid"

type StudentInfoModel struct {
	StudentID    uuid.UUID `json:"studentID"`
	StudentName  string    `json:"studentName"`
	StudentEmail string    `json:"studentEmail,omitempty"`
	Description  string    `json:"description,omitempty"`
	ProfilePic   string    `json:"profilePic,omitempty"`
	TagsLiked    []string  `json:"tagsLiked,omitempty"`
	TagsDisliked []string  `json:"tagsDisliked,omitempty"`
}
