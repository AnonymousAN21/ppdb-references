package controller

import (
	"net/http"
	"spmb/back-end/api/handlers"
	"spmb/back-end/api/internal/service"
	"spmb/back-end/models"
	"spmb/back-end/models/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AccomplishmentOrgRefController struct {
	s service.AccomplishmentService
}

func NewAccomplishmentOrgRefController(s service.AccomplishmentService) *AccomplishmentOrgRefController {
	return &AccomplishmentOrgRefController{s: s}
}

// GetAllOrganizationRef handles fetching all filtered organization refs
func (f *AccomplishmentOrgRefController) GetAllOrganizationRef(ctx *gin.Context) {
	var filter models.AccomplishmentOrganizationRef
	role := ctx.Query("role")
	isActive := ctx.Query("is_active")

	if isActive == "" {
		filter.IsActive = true
	} else {
		status, err := strconv.ParseBool(isActive)
		if err == nil {
			filter.IsActive = status
		}
	}

	if role != "" {
		filter.Role = role
	}

	statusCode, res, err := f.s.GetOrganizationRef(filter)
	if err.Message != "" {
		ctx.JSON(statusCode, err)
		return
	}

	ctx.JSON(statusCode, res)
}

// CreateOrganizationRef handles creating a new organization ref
func (f *AccomplishmentOrgRefController) CreateOrganizationRef(ctx *gin.Context) {
	var body models.AccomplishmentOrganizationRef

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak memiliki Akses",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}

	userUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak memiliki Akses",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}

	body.CreatedBy = userUUID
	body.UpdatedBy = userUUID
	body.IsActive = true

	if err := ctx.ShouldBindJSON(&body); err != nil {
		validationErrors := make(map[string]string)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range errs {
				jsonField := handlers.GetJSONFieldName(fieldErr, models.AccomplishmentOrganizationRef{})
				validationErrors[jsonField] = handlers.CustomErrorMessage(fieldErr, models.AccomplishmentOrganizationRef{})
			}
		}
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   validationErrors,
		})
		return
	}

	statusCode, res, err := f.s.CreateOrganizationRef(body)
	if err.Message != "" {
		ctx.JSON(statusCode, err)
		return
	}

	ctx.JSON(statusCode, res)
}

func (f *AccomplishmentOrgRefController) Activate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Error:   map[string]string{"id": err.Error()},
		})
		return
	}

	code, res, errRes := f.s.ActivateOrgRef(id)
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
}

func (f *AccomplishmentOrgRefController) Deactivate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Error:   map[string]string{"id": err.Error()},
		})
		return
	}

	code, res, errRes := f.s.DeactivateOrgRef(id)
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
}

func (f *AccomplishmentOrgRefController) UpdateRefOrg(ctx *gin.Context) {
	var body models.UpdateAccomplishmentOrganizationRef
	idparam := ctx.Param("id")

	id, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Error: map[string]string{
				"Query": "ID harus berupa angka",
			},
		})
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		validationErrors := make(map[string]string)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range errs {
				jsonField := handlers.GetJSONFieldName(fieldErr, models.StudentFileType{})
				validationErrors[jsonField] = handlers.CustomErrorMessage(fieldErr, models.StudentFileType{})
			}
		}
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   validationErrors,
		})
		return
	}

	code, res, errResp := f.s.UpdateRefOrg(body, id)

	if errResp.Message != "" {
		ctx.JSON(code, errResp)
		return
	}

	ctx.JSON(code, res)

}
