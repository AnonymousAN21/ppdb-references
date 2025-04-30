package models

import (
	"time"

	"github.com/google/uuid"
)

type StudentFileType struct {
	StudentFileID uint      `json:"student_file_type_id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" binding:"required"`
	Type          string    `json:"type" binding:"required,oneof=general domicile accomplishment job_transfer social_assistance"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime" `
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime" `
	CreatedBy     uuid.UUID `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy     uuid.UUID `json:"updated_by" gorm:"type:char(36);not null"`
}
