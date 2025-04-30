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

type FaqController struct {
	s service.FaqServices
}

func NewFaqController(s service.FaqServices) *FaqController {
	return &FaqController{s: s}
}

func (f *FaqController) Create(ctx *gin.Context) {
	var body models.FAQ
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

	newFaq := models.FAQ{
		Question:  body.Question,
		Answer:    body.Answer,
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
				jsonField := handlers.GetJSONFieldName(fieldErr, models.FAQ{})
				validationErrors[jsonField] = handlers.CustomErrorMessage(fieldErr, models.FAQ{})
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

func (f *FaqController) GetAll(ctx *gin.Context) {
	// Extract query parameters
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	question := ctx.Query("question")
	answer := ctx.Query("answer")
	keyword := ctx.Query("keyword")
	isActiveStr := ctx.Query("is_active")

	filter := models.FAQ{
		Question: question,
		Answer:   answer,
	}

	if keyword != "" {
		filter.Question = keyword
	}

	if isActiveStr == "" {
		filter.IsActive = true
	} else {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			filter.IsActive = isActive
		}
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

func (f *FaqController) Activate(ctx *gin.Context) {
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

func (f *FaqController) GetById(ctx *gin.Context) {
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

	statusCode, response, errorResponse := f.s.GetById(id)
	if errorResponse.Message != "" {
		ctx.JSON(statusCode, errorResponse)
		return
	}

	ctx.JSON(statusCode, response)
}

func (f *FaqController) Deactivate(ctx *gin.Context) {
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
