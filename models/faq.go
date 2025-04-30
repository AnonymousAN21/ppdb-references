package models

import (
	"time"

	"github.com/google/uuid"
)

type FAQ struct {
	FaqID     uint      `json:"faq_id" gorm:"primaryKey;autoIncrement"`
	Question  string    `json:"question" gorm:"type:varchar(50)" binding:"required,min=10"`
	Answer    string    `json:"answer" gorm:"type:varchar(50)" binding:"required,min=10"`
	IsActive  bool      `json:"is_active" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:char(36);not null"`
}
type FAQBody struct {
	Question string `json:"question" binding:"required,min=10"`
	Answer   string `json:"answer" binding:"required,min=10"`
	IsActive *bool  `json:"is_active"`
}

type FAQResponse struct {
	FaqID     uint      `json:"faq_id" gorm:"primaryKey;autoIncrement;type:int(6)"`
	Question  string    `json:"question" gorm:"type:varchar(50)" binding:"required,min=10"`
	Answer    string    `json:"answer" gorm:"type:varchar(50)" binding:"required,min=10"`
	IsActive  bool      `json:"is_active" binding:"required"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at" `
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}
