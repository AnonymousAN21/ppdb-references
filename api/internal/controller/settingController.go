package controller

import (
	"net/http"
	"spmb/back-end/api/handlers"
	"spmb/back-end/api/internal/service"
	"spmb/back-end/models"
	"spmb/back-end/models/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SettingController struct {
	s service.SettingService
}

func NewSettingController(s service.SettingService) *SettingController {
	return &SettingController{s: s}
}

func (sc *SettingController) GetAll(ctx *gin.Context) {
	code, res, errRes := sc.s.GetAll()
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}
	ctx.JSON(code, res)
}

func (sc *SettingController) Create(ctx *gin.Context) {
	var setting models.Settings
	if err := ctx.ShouldBindJSON(&setting); err != nil {
		validationErrors := make(map[string]string)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range errs {
				jsonField := handlers.GetJSONFieldName(fieldErr, models.Settings{})
				validationErrors[jsonField] = handlers.CustomErrorMessage(fieldErr, models.Settings{})
			}
		}
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Tidak Valid",
			Error:   validationErrors,
		})
		return
	}
	code, res, errRes := sc.s.Create(setting)
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}
	ctx.JSON(code, res)
}

func (sc *SettingController) Activate(ctx *gin.Context) {
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
	code, successResponse, errorResponse := sc.s.Activate(userUUID)

	if errorResponse.Message != "" {
		ctx.JSON(code, errorResponse)
		return
	}

	ctx.JSON(code, successResponse)
}

func (sc *SettingController) Deactivate(ctx *gin.Context) {
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

	code, successResponse, errorResponse := sc.s.Deactivate(userUUID)

	if errorResponse.Message != "" {
		ctx.JSON(code, errorResponse)
		return
	}

	ctx.JSON(code, successResponse)
}
