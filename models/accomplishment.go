package models

import (
	"time"

	"github.com/google/uuid"
)

type AccomplishmentRef struct {
	AccomplishmentRefId uint      `json:"accomplishment_ref_id" gorm:"primaryKey:autoIncrement;type:int(6)"`
	Level               string    `json:"level" binding:"oneof=kabupaten/kota provinsi nasionl internasional"`
	Ranking             int       `json:"ranking"`
	SoloPoint           int       `json:"solo_point"`
	SmallGroupPoint     int       `json:"small_group_point"`
	BigGroupPoint       int       `json:"big_group_point"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:datetime" `
	UpdatedAt           time.Time `json:"updated_at" gorm:"type:datetime" `
	CreatedBy           uuid.UUID `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy           uuid.UUID `json:"updated_by" gorm:"type:char(36);not null"`
}

type UpdateAccomplishmentRef struct {
	SoloPoint       int `json:"solo_point"`
	SmallGroupPoint int `json:"small_group_point"`
	BigGroupPoint   int `json:"big_group_point"`
}
type AccomplishmentOrganizationRef struct {
	AccomplishmentOrganizationRefId uint      `json:"accomplishment_organization_ref_id" gorm:"primaryKey:autoIncrement;type:int(6)"`
	Role                            string    `json:"role"`
	Point                           int       `json:"point"`
	IsActive                        bool      `json:"is_active"`
	CreatedAt                       time.Time `json:"created_at" gorm:"type:datetime" `
	UpdatedAt                       time.Time `json:"updated_at" gorm:"type:datetime" `
	CreatedBy                       uuid.UUID `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy                       uuid.UUID `json:"updated_by" gorm:"type:char(36);not null"`
}
type UpdateAccomplishmentOrganizationRef struct {
	Point int `json:"point"`
}
