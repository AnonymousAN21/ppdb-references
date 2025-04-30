package service

import (
	"fmt"
	"spmb/back-end/models"
	"spmb/back-end/models/response"
	"spmb/back-end/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StfService interface {
	Create(body models.StudentFileType) (int, response.SuccessResponse, response.ErrorResponse)
	GetAll(page, limit int, filter models.StudentFileType) (int, response.SuccessResponseGet, response.ErrorResponse)
	GetType() (int, response.SuccessResponseGet, response.ErrorResponse)
	Activate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
	Deactivate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
	GetById(id uint) (int, response.SuccessResponseGet, response.ErrorResponse)
}

type stfService struct {
	db *gorm.DB
}

// GetType implements StfService.
func (s *stfService) GetType() (int, response.SuccessResponseGet, response.ErrorResponse) {
	return 200, response.SuccessResponseGet{
		Success: true,
		Message: "Data Ditemukan",
		Data:    []string{"domicile", "general", "accomplishment", "job_transfer", "social_assistance"},
	}, response.ErrorResponse{}
}

// Activate implements StfService.
func (s *stfService) Activate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.StudentFileType
	if err := s.db.Where(&models.StudentFileType{StudentFileID: id}).First(&data).Error; err != nil {
		return 404, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if data.IsActive {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
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
			Message: "Gagal menghapus data",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 201, response.SuccessResponse{
		Success: true,
		Message: "Berhasil mengaktifkan data",
		Data: gin.H{
			"student_file_type_id": id,
			"is_active":            true,
		},
	}, response.ErrorResponse{}
}

// Create implements StfService.
func (s *stfService) Create(body models.StudentFileType) (int, response.SuccessResponse, response.ErrorResponse) {
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

// Deactivate implements StfService.
func (s *stfService) Deactivate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.StudentFileType
	if err := s.db.Where(&models.StudentFileType{StudentFileID: id}).First(&data).Error; err != nil {
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
			"student_file_type_id": id,
			"is_active":            false,
		},
	}, response.ErrorResponse{}
}

// GetAll implements StfService.
func (s *stfService) GetAll(page int, limit int, filter models.StudentFileType) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.StudentFileType
	var total int64

	if page < 1 {
		return 404, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Halaman Tidak ditemukan",
			Error: map[string]string{
				"Page": "Page minimal harus 1",
			},
		}
	}
	if limit < 1 {
		return 404, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Halaman Tidak ditemukan",
			Error: map[string]string{
				"Page": "limit minimal harus 1",
			},
		}
	}
	if limit > 20 {
		return 400, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Limit": "Limit hanya bisa 20",
			},
		}
	}

	query := s.db.Model(&models.StudentFileType{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", filter.Name+"%")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	query = query.Where("is_active = ?", filter.IsActive)

	if err := query.Count(&total).Error; err != nil {
		return 400, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if err := query.Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return 500, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}
	FormattedData, err := utils.FormatStructSlice(data, "2006-01-02 15:04:05")

	if err != nil {
		return 500, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}
	res := response.SuccessResponseGet{
		Success: true,
		Message: "Data Ditemukan",
		Data:    FormattedData,
		Pagination: response.Pagination{
			Page:         page,
			Limit:        limit,
			Has_next:     (offset + limit) < int(total),
			Has_previous: offset > 0,
			Next_url:     fmt.Sprintf("/faq?page=%d&limit=%d", page+1, limit),
			Prev_url:     fmt.Sprintf("/faq?page=%d&limit=%d", page-1, limit),
		},
	}

	return 200, res, response.ErrorResponse{}
}

// GetById implements StfService.
func (s *stfService) GetById(id uint) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.StudentFileType

	if err := s.db.Where(&models.StudentFileType{StudentFileID: id}).First(&data).Error; err != nil {
		return 404, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	formattedData, err := utils.FormatStructSlice(data, "2006-01-02 15:04:05")
	if err != nil {
		return 500, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}
	res := response.SuccessResponseGet{
		Success: true,
		Message: "Data Ditemukan",
		Data:    formattedData,
	}

	return 200, res, response.ErrorResponse{}
}

func NewStfService(db *gorm.DB) StfService {
	return &stfService{db: db}
}
