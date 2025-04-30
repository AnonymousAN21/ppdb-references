package models

import (
	"time"

	"github.com/google/uuid"
)

type Settings struct {
	SettingID      uint      `json:"setting_id" gorm:"primaryKey;autoIncrement;type:int(6)"`
	AppName        string    `json:"app_name" binding:"required"`
	Modal          bool      `json:"modal" binding:"required"`
	ModalHeader    string    `json:"modal_header" binding:"required"`
	ModalBody      string    `json:"modal_body" binding:"required,min=20"`
	SchedulePic    string    `json:"schedule_picture" binding:"required"`
	RequirementPic string    `json:"requirement_picture" binding:"required"`
	HowToRegPic    string    `json:"how_to_register_picture" binding:"required"`
	HowToRegVid    string    `json:"how_to_register_video" binding:"required"`
	IsActive       bool      `json:"is_active" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime" `
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:datetime" `
	CreatedBy      uuid.UUID `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy      uuid.UUID `json:"updated_by" gorm:"type:char(36);not null"`
}
type SettingResponse struct {
	SettingID      uint   `json:"setting_id" gorm:"primaryKey;autoIncrement;type:int(6)"`
	AppName        string `json:"app_name"`
	Modal          bool   `json:"modal"`
	ModalHeader    string `json:"modal_header"`
	ModalBody      string `json:"modal_body" binding:"required,min=20"`
	SchedulePic    string `json:"schedule_picture"`
	RequirementPic string `json:"requirement_picture"`
	HowToRegPic    string `json:"how_to_register_picture"`
	HowToRegVid    string `json:"how_to_register_video"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at" gorm:"type:datetime" `
	UpdatedAt      string `json:"updated_at" gorm:"type:datetime" `
	CreatedBy      string `gorm:"type:uuid"`
	UpdatedBy      string `gorm:"type:uuid"`
}
