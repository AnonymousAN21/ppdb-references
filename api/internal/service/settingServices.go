package service

import (
	"spmb/back-end/models"
	"spmb/back-end/models/response"
	"spmb/back-end/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingService interface {
	GetAll() (int, response.SuccessResponse, response.ErrorResponse)
	Create(body models.Settings) (int, response.SuccessResponse, response.ErrorResponse)
	Activate(Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
	Deactivate(Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
}

type settingService struct {
	db *gorm.DB
}

func (s *settingService) Activate(Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.Settings
	if err := s.db.First(&data).Error; err != nil {
		return 404, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if data.IsActive {
		return 404, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data telah aktif",
			Error: map[string]string{
				"Query": "Data sudah aktif",
			},
		}
	}

	data.IsActive = true
	data.UpdatedBy = Updatedby

	if err := s.db.Save(&data).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal mengaktifasi data",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 201, response.SuccessResponse{
		Success: true,
		Message: "Berhasil mengaktifasi data",
		Data: gin.H{
			"is_active": true,
		},
	}, response.ErrorResponse{}
}

func (s *settingService) Deactivate(Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.Settings
	if err := s.db.First(&data).Error; err != nil {
		return 404, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if !data.IsActive {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data tidak aktif",
			Error: map[string]string{
				"Query": "Data sudah Tidak aktif",
			},
		}
	}

	data.IsActive = false
	data.UpdatedBy = Updatedby

	if err := s.db.Save(&data).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal menghapus data",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 201, response.SuccessResponse{
		Success: true,
		Message: "Berhasil menghapus data",
		Data: gin.H{
			"is_active": false,
		},
	}, response.ErrorResponse{}
}

func (s *settingService) Create(body models.Settings) (int, response.SuccessResponse, response.ErrorResponse) {
	var count int64
	if err := s.db.Model(&models.Settings{}).Count(&count).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": "Gagal Menghitung Jumlah Data",
			},
		}
	}
	if count == 0 {
		if err := s.db.Create(&body).Error; err != nil {
			return 400, response.SuccessResponse{}, response.ErrorResponse{
				Success: false,
				Message: "Kesalahan Server",
				Error: map[string]string{
					"Query": "Gagal Menambahkan data",
				},
			}
		}

		return 201, response.SuccessResponse{
			Success: true,
			Message: "Berhasil menambahkan data",
			Data:    body,
		}, response.ErrorResponse{}
	}

	if count == 1 {
		var existing models.Settings
		if err := s.db.First(&existing).Error; err != nil {
			return 400, response.SuccessResponse{}, response.ErrorResponse{
				Success: false,
				Message: "Kesalahan Server",
				Error: map[string]string{
					"Query": "Gagal mengambil data yang ada",
				},
			}
		}

		if err := s.db.Model(&existing).Updates(body).Error; err != nil {
			return 400, response.SuccessResponse{}, response.ErrorResponse{
				Success: false,
				Message: "Kesalahan Server",
				Error: map[string]string{
					"Query": "Gagal mengupdate data",
				},
			}
		}

		if err := s.db.First(&existing, existing.SettingID).Error; err != nil {
			return 400, response.SuccessResponse{}, response.ErrorResponse{
				Success: false,
				Message: "Kesalahan Server",
				Error: map[string]string{
					"Query": "Gagal mengambil data yang telah diupdate",
				},
			}
		}

		return 200, response.SuccessResponse{
			Success: true,
			Message: "Berhasil mengupdate data",
			Data:    existing,
		}, response.ErrorResponse{}
	}

	return 400, response.SuccessResponse{}, response.ErrorResponse{
		Success: false,
		Message: "Kesalahan Data",
		Error: map[string]string{
			"Data": "Terdapat lebih dari satu data pengaturan, tabel seharusnya hanya memiliki satu data.",
		},
	}
}

func (s *settingService) GetAll() (int, response.SuccessResponse, response.ErrorResponse) {
	var data []models.Settings

	if err := s.db.Where(&models.Settings{IsActive: true}).Find(&data).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}
	formattedData := utils.FormatSettings(data)
	if len(formattedData) == 0 {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Sedang Maintenance",
			Error: map[string]string{
				"Query": "Settings dalam status tidak diaktifkan",
			},
		}
	}

	return 200, response.SuccessResponse{
		Success: true,
		Message: "Data ditemukan",
		Data:    formattedData,
	}, response.ErrorResponse{}
}

func NewSettingService(db *gorm.DB) SettingService {
	return &settingService{db: db}
}
