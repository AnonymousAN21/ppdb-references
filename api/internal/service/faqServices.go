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

type FaqServices interface {
	Create(body models.FAQ) (int, response.SuccessResponse, response.ErrorResponse)
	GetAll(page, limit int, filter models.FAQ) (int, response.SuccessResponseGet, response.ErrorResponse)
	GetById(id uint) (int, response.SuccessResponseGet, response.ErrorResponse)
	Activate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
	Deactivate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse)
}

type faqServices struct {
	db *gorm.DB
}

// GetById implements FaqServices.
func (f *faqServices) GetById(id uint) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.FAQ

	if err := f.db.Where(&models.FAQ{FaqID: id}).First(&data).Error; err != nil {
		return 404, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	formattedData := utils.FormatFAQ(data)

	res := response.SuccessResponseGet{
		Success: true,
		Message: "Data Ditemukan",
		Data:    formattedData,
	}

	return 200, res, response.ErrorResponse{}
}

// Activate FaqServices.
func (f *faqServices) Activate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.FAQ
	var count int64
	if err := f.db.Where(&models.FAQ{FaqID: id}).First(&data).Error; err != nil {
		return 404, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data Tidak ditemukan",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if err := f.db.Model(&models.FAQ{}).Where(models.FAQ{IsActive: true}).Count(&count).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": "Gagal Menghitung Jumlah Data",
			},
		}
	}

	if count >= 10 {
		return 422, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Maximal 10 data aktif",
			Error: map[string]string{
				"Limit": "Limit hanya bisa 10 data aktif",
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

	if err := f.db.Save(&data).Error; err != nil {
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
			"faq_id":    id,
			"is_active": true,
		},
	}, response.ErrorResponse{}
}

// Deactivate implements FaqServices.
func (f *faqServices) Deactivate(id uint, Updatedby uuid.UUID) (int, response.SuccessResponse, response.ErrorResponse) {
	var data models.FAQ
	if err := f.db.Where(&models.FAQ{FaqID: id}).First(&data).Error; err != nil {
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

	if err := f.db.Save(&data).Error; err != nil {
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
			"faq_id":    id,
			"is_active": false,
		},
	}, response.ErrorResponse{}
}

// Create FaqServices.
func (f *faqServices) Create(body models.FAQ) (int, response.SuccessResponse, response.ErrorResponse) {
	var count int64
	if err := f.db.Model(&models.FAQ{}).Where(&models.FAQ{IsActive: true}).Count(&count).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": "Gagal Menghitung Jumlah Data",
			},
		}
	}

	if count >= 10 {
		return 422, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Maximal 10 data aktif",
			Error: map[string]string{
				"Limit": "Limit hanya bisa 10 data aktif",
			},
		}
	}

	if err := f.db.Create(&body).Error; err != nil {
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

// GetAll FaqServices.
func (f *faqServices) GetAll(page int, limit int, filter models.FAQ) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.FAQ
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

	query := f.db.Model(&models.FAQ{})

	if filter.Question != "" {
		query = query.Where("question LIKE ?", filter.Question+"%")
	}

	if filter.Answer != "" {
		query = query.Where("answer = ?", filter.Answer)
	}

	query = query.Where("is_active = ?", filter.IsActive)

	if err := query.Count(&total).Error; err != nil {
		return 400, response.SuccessResponseGet{}, response.ErrorResponse{
			Success: false,
			Message: "Kesalahan Server",
			Error: map[string]string{
				"Query": "Gagal Menghitung Jumlah Data",
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

func NewFaqService(db *gorm.DB) FaqServices {
	return &faqServices{db: db}
}
