package service

import (
	"spmb/back-end/models"
	"spmb/back-end/models/response"
	"spmb/back-end/utils"

	"gorm.io/gorm"
)

type AccomplishmentService interface {
	GetRef(filter models.AccomplishmentRef) (int, response.SuccessResponseGet, response.ErrorResponse)
	GetOrganizationRef(filter models.AccomplishmentOrganizationRef) (int, response.SuccessResponseGet, response.ErrorResponse)
	CreateRef(body models.AccomplishmentRef) (int, response.SuccessResponse, response.ErrorResponse)
	CreateOrganizationRef(body models.AccomplishmentOrganizationRef) (int, response.SuccessResponse, response.ErrorResponse)
	UpdateRef(body models.UpdateAccomplishmentRef, id int) (int, response.SuccessResponse, response.ErrorResponse)
	UpdateRefOrg(body models.UpdateAccomplishmentOrganizationRef, id int) (int, response.SuccessResponse, response.ErrorResponse)
	ActivateOrgRef(id int) (int, response.SuccessResponse, response.ErrorResponse)
	DeactivateOrgRef(id int) (int, response.SuccessResponse, response.ErrorResponse)
	ActivateAccRef(id int) (int, response.SuccessResponse, response.ErrorResponse)
	DeactivateAccRef(id int) (int, response.SuccessResponse, response.ErrorResponse)
	// Calculate(accomplishment models.Accomplishment) (int, response.ErrorResponse)
}

type accomplishmentService struct {
	db *gorm.DB
}

// UpdateRef implements AccomplishmentService.
func (a *accomplishmentService) UpdateRef(body models.UpdateAccomplishmentRef, id int) (int, response.SuccessResponse, response.ErrorResponse) {
	var updateData = map[string]interface{}{}

	if body.BigGroupPoint > 0 {
		updateData["big_group_point"] = body.BigGroupPoint
	}
	if body.SmallGroupPoint > 0 {
		updateData["small_group_point"] = body.SmallGroupPoint
	}
	if body.SoloPoint > 0 {
		updateData["solo_point"] = body.SoloPoint
	}

	if len(updateData) == 0 {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Message: "Tidak ada data yang diubah",
		}
	}

	if err := a.db.Model(&models.AccomplishmentRef{}).
		Where("accomplishment_ref_id = ?", id).
		Updates(updateData).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal memperbarui data",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 200, response.SuccessResponse{
		Success: true,
		Message: "Berhasil memperbarui data referensi prestasi",
		Data:    updateData,
	}, response.ErrorResponse{}
}

// UpdateRefOrg implements AccomplishmentService.
func (a *accomplishmentService) UpdateRefOrg(body models.UpdateAccomplishmentOrganizationRef, id int) (int, response.SuccessResponse, response.ErrorResponse) {
	var updateData = map[string]interface{}{}

	if body.Point > 0 {
		updateData["point"] = body.Point // asumsi pakai BigGroupPoint untuk org
	}

	if len(updateData) == 0 {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Message: "Tidak ada data yang diubah",
		}
	}

	if err := a.db.Model(&models.AccomplishmentOrganizationRef{}).
		Where("accomplishment_organization_ref_id = ?", id).
		Updates(updateData).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal memperbarui data organisasi",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 200, response.SuccessResponse{
		Success: true,
		Message: "Berhasil memperbarui data referensi organisasi",
		Data:    updateData,
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) ActivateAccRef(id int) (int, response.SuccessResponse, response.ErrorResponse) {
	if err := a.db.Model(&models.AccomplishmentRef{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": true,
		}).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal mengaktifkan data",
			Error:   map[string]string{"query": err.Error()},
		}
	}
	return 200, response.SuccessResponse{
		Success: true,
		Message: "Data berhasil diaktifkan",
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) DeactivateAccRef(id int) (int, response.SuccessResponse, response.ErrorResponse) {
	if err := a.db.Model(&models.AccomplishmentRef{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": false,
		}).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal menonaktifkan data",
			Error:   map[string]string{"query": err.Error()},
		}
	}
	return 200, response.SuccessResponse{
		Success: true,
		Message: "Data berhasil dinonaktifkan",
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) ActivateOrgRef(id int) (int, response.SuccessResponse, response.ErrorResponse) {
	if err := a.db.Model(&models.AccomplishmentOrganizationRef{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": true,
		}).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal mengaktifkan data",
			Error:   map[string]string{"query": err.Error()},
		}
	}
	return 200, response.SuccessResponse{
		Success: true,
		Message: "Data berhasil diaktifkan",
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) DeactivateOrgRef(id int) (int, response.SuccessResponse, response.ErrorResponse) {
	if err := a.db.Model(&models.AccomplishmentOrganizationRef{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": false,
		}).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Gagal menonaktifkan data",
			Error:   map[string]string{"query": err.Error()},
		}
	}
	return 200, response.SuccessResponse{
		Success: true,
		Message: "Data berhasil dinonaktifkan",
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) CreateOrganizationRef(body models.AccomplishmentOrganizationRef) (int, response.SuccessResponse, response.ErrorResponse) {
	var data []models.AccomplishmentOrganizationRef
	if err := a.db.Model(&models.AccomplishmentOrganizationRef{}).Where(&models.AccomplishmentOrganizationRef{Role: body.Role}).Find(&data).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if len(data) > 0 {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": "Data sudah ada",
			},
		}
	}
	if err := a.db.Create(&body).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 201, response.SuccessResponse{
		Success: true,
		Message: "Berhasil menambahkan data",
		Data:    body,
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) CreateRef(body models.AccomplishmentRef) (int, response.SuccessResponse, response.ErrorResponse) {
	var data []models.AccomplishmentRef
	if err := a.db.Model(&models.AccomplishmentRef{}).Where(&models.AccomplishmentRef{Level: body.Level}).Find(&data).Error; err != nil {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	if len(data) > 0 {
		return 400, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": "Data sudah ada",
			},
		}
	}

	if err := a.db.Create(&body).Error; err != nil {
		return 500, response.SuccessResponse{}, response.ErrorResponse{
			Success: false,
			Message: "Data sudah ada",
			Error: map[string]string{
				"Query": err.Error(),
			},
		}
	}

	return 201, response.SuccessResponse{
		Success: true,
		Message: "Berhasil menambahkan data",
		Data:    body,
	}, response.ErrorResponse{}
}

func (a *accomplishmentService) GetOrganizationRef(filter models.AccomplishmentOrganizationRef) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.AccomplishmentOrganizationRef
	query := a.db.Model(&models.AccomplishmentOrganizationRef{})

	if filter.Role != "" {
		query = query.Where("role LIKE ?", filter.Role+"%")
	}
	query = query.Where("is_active = ?", filter.IsActive)

	if err := query.Find(&data).Error; err != nil {
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
	}
	return 200, res, response.ErrorResponse{}
}

func (a *accomplishmentService) GetRef(filter models.AccomplishmentRef) (int, response.SuccessResponseGet, response.ErrorResponse) {
	var data []models.AccomplishmentRef
	query := a.db.Model(&models.AccomplishmentRef{})

	if filter.Level != "" {
		query = query.Where("level = ?", filter.Level)
	}

	if filter.Ranking > 0 {
		query = query.Where("ranking = ?", filter.Ranking)
	}
	query = query.Where("is_active = ?", filter.IsActive)

	if err := query.Find(&data).Error; err != nil {
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
	}
	return 200, res, response.ErrorResponse{}
}

func NewAccomplishmentService(db *gorm.DB) AccomplishmentService {
	return &accomplishmentService{db: db}
}
