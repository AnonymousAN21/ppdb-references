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

type AccomplishmentRefController struct {
	s service.AccomplishmentService
}

// NewPointSimulationController creates a new instance of PointSimulationController.
func NewAccomplishmentRefController(s service.AccomplishmentService) *AccomplishmentRefController {
	return &AccomplishmentRefController{s: s}
}

// Calculate handles the calculation of points based on accomplishments and report.
// func (f *PointSimulationController) Calculate(ctx *gin.Context) {
// 	var body models.Accomplishment

// 	if err := ctx.ShouldBindJSON(&body); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
// 			Success: false,
// 			Message: "Tidak Valid",
// 			Error:   map[string]string{"error": "kesalahan Server"},
// 		})
// 	}
// 	points, err := f.s.Calculate(body)

// 	if err.Message != "" {
// 		ctx.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"points": points})
// }

func (f *AccomplishmentRefController) GetAllRef(ctx *gin.Context) {
	var filter models.AccomplishmentRef
	level := ctx.Query("level")
	rank := ctx.Query("ranking")
	isActive := ctx.Query("is_active")

	if isActive == "" {
		filter.IsActive = true
	} else {
		isActive, err := strconv.ParseBool(isActive)
		if err == nil {
			filter.IsActive = isActive
		}
	}

	if level != "" {
		filter.Level = level
	}

	if rank != "" {
		ranking, err := strconv.Atoi(rank)
		if err == nil {
			filter.Ranking = ranking
		}
	}

	statusCode, res, err := f.s.GetRef(filter)
	if err.Message != "" {
		ctx.JSON(statusCode, err)
		return
	}

	ctx.JSON(statusCode, res)
}

func (f *AccomplishmentRefController) CreateRef(ctx *gin.Context) {
	var body models.AccomplishmentRef

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

	newRef := models.AccomplishmentRef{
		Level:           body.Level,
		Ranking:         body.Ranking,
		SoloPoint:       body.SoloPoint,
		SmallGroupPoint: body.SmallGroupPoint,
		BigGroupPoint:   body.BigGroupPoint,
		IsActive:        true,
		CreatedBy:       userUUID,
		UpdatedBy:       userUUID,
	}

	if err := ctx.ShouldBindJSON(&newRef); err != nil {
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
	code, res, err := f.s.CreateRef(newRef)

	// Check if there is an error message
	if err.Message != "" {
		ctx.JSON(code, err)
		return
	}

	ctx.JSON(code, res)
}

func (f *AccomplishmentRefController) Activate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Error:   map[string]string{"id": err.Error()},
		})
		return
	}

	code, res, errRes := f.s.ActivateAccRef(id)
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
}

func (f *AccomplishmentRefController) Deactivate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Error:   map[string]string{"id": err.Error()},
		})
		return
	}

	code, res, errRes := f.s.DeactivateAccRef(id)
	if errRes.Message != "" {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
}

func (f *AccomplishmentRefController) UpdateRef(ctx *gin.Context) {
	var body models.UpdateAccomplishmentRef
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

	code, res, errResp := f.s.UpdateRef(body, id)

	if errResp.Message != "" {
		ctx.JSON(code, errResp)
		return
	}

	ctx.JSON(code, res)

}
