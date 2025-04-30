package controller

import (
	"net/http"
	"spmb/back-end/api/handlers"
	"spmb/back-end/models"
	"spmb/back-end/models/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Calc(ctx *gin.Context) {
	var body models.Point
	SchoolType := ctx.Param("type")

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

	a := body.Accomplishment
	b := body.AccomplishmentOrganization
	c := body.Interview
	d := body.Report_Point

	var bobotA float32
	var bobotB float32
	var bobotC float32
	var bobotD float32

	if SchoolType == "SMK" {
		bobotA = 15.0 / 100
		bobotB = 20.0 / 100
		bobotC = 35.0 / 100
		bobotD = 30.0 / 100
	} else if SchoolType == "SMA" {
		bobotA = 45.0 / 100
		bobotB = 25.0 / 100
		bobotC = 0.0 / 100
		bobotD = 30.0 / 100
	} else {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Gagal Menghitung Nilai",
			Error: map[string]string{
				"Error": "Type harus SMA/SMK",
			},
		})
		return
	}

	point := float32(a)*bobotA + float32(b)*bobotB + float32(c)*bobotC + float32(d)*bobotD

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Berhasil Menghitung Nilai",
		Data: map[string]float32{
			"point": point,
		},
	})
}
