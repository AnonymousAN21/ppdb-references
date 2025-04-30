package controller

import (
	"net/http"
	"spmb/back-end/api/handlers"
	"spmb/back-end/api/internal/service"
	"spmb/back-end/models"
	"spmb/back-end/models/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type StfController struct {
	s service.StfService
}

func NewStfController(s service.StfService) *StfController {
	return &StfController{s: s}
}

func (f *StfController) GetType(ctx *gin.Context) {
	code, res, err := f.s.GetType()
	if err.Message != "" {
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(code, res)
}
func (f *StfController) Create(ctx *gin.Context) {
	var body models.StudentFileType
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

	newFaq := models.StudentFileType{
		Name:      body.Name,
		Type:      body.Type,
		IsActive:  true,
		CreatedAt: time.Now(),
		CreatedBy: userUUID,
		UpdatedAt: time.Now(),
		UpdatedBy: userUUID,
	}

	if err := ctx.ShouldBindJSON(&newFaq); err != nil {
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
	code, res, err := f.s.Create(newFaq)

	// Check if there is an error message
	if err.Message != "" {
		ctx.JSON(code, err)
		return
	}

	ctx.JSON(code, res)
}

func (f *StfController) GetAll(ctx *gin.Context) {
	var filter models.StudentFileType
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	keyword := ctx.Query("keyword")
	typeStr := ctx.Query("type")
	isActiveStr := ctx.Query("is_active")

	if isActiveStr == "" {
		filter.IsActive = true
	} else {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			filter.IsActive = isActive
		}
	}

	if keyword != "" {
		filter.Name = keyword
	}

	if typeStr != "" {
		filter.Type = typeStr
	}
	page, _ := strconv.Atoi(pageStr)

	limit, _ := strconv.Atoi(limitStr)

	statusCode, response, errorResponse := f.s.GetAll(page, limit, filter)
	if errorResponse.Message != "" {
		ctx.JSON(statusCode, errorResponse)
		return
	}
	ctx.JSON(statusCode, response)
}

func (f *StfController) Activate(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id64, err := strconv.ParseUint(idStr, 10, 32)
	id := uint(id64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"id": "Invalid ID, must be a valid integer"},
		})
		return
	}

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}
	userUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}

	code, successResponse, errorResponse := f.s.Activate(id, userUUID)

	if errorResponse.Message != "" {
		ctx.JSON(code, errorResponse)
		return
	}

	ctx.JSON(code, successResponse)
}

func (f *StfController) Deactivate(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id64, err := strconv.ParseUint(idStr, 10, 32)
	id := uint(id64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"id": "Invalid ID, must be a valid integer"},
		})
		return
	}

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}
	userUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"authorization": "Invalid User ID"},
		})
		return
	}

	code, successResponse, errorResponse := f.s.Deactivate(id, userUUID)

	if errorResponse.Message != "" {
		ctx.JSON(code, errorResponse)
		return
	}

	ctx.JSON(code, successResponse)
}

func (f *StfController) GetById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id64, err := strconv.ParseUint(idStr, 10, 32)
	id := uint(id64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   map[string]string{"id": "Invalid ID, must be a valid integer"},
		})
		return
	}

	code, successResponse, errorResponse := f.s.GetById(id)
	if errorResponse.Message != "" {
		ctx.JSON(code, errorResponse)
		return
	}

	ctx.JSON(code, successResponse)
}
