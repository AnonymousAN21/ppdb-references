package main

import (
	"fmt"
	"spmb/back-end/api/internal/controller"
	"spmb/back-end/api/internal/middleware"
	"spmb/back-end/api/internal/service"
	"spmb/back-end/config"
	"spmb/back-end/config/connection"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// var jwtSecret = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ0NzM4NDksImlhdCI6MTc0NDQzNzg0OSwiaXNfc3R1ZGVudCI6ZmFsc2UsIm5iZiI6MTc0NDQzNzg0OSwib3JnYW5pemF0aW9uX2lkIjoiIiwicm9sZV9pZCI6MCwidXNlcl9pZCI6Ijk4MThmMzMyLTczMDctNDYxNC05NzcxLTNlODQ5NzY5YmYzOCJ9.h1O3rN73fCK5CAk0wAEek6GBFQBoSQ21hqXhNmWgCnk")

func main() {
	jwtgen()

	db := connection.Conn()
	FaqService := service.NewFaqService(db)
	FaqController := controller.NewFaqController(FaqService)
	SettingService := service.NewSettingService(db)
	SettingsController := controller.NewSettingController(SettingService)
	StfService := service.NewStfService(db)
	StfController := controller.NewStfController(StfService)
	AccomplishmentService := service.NewAccomplishmentService(db)
	AccomplishmentController := controller.NewAccomplishmentRefController(AccomplishmentService)
	AccomplishmentOrganizationController := controller.NewAccomplishmentOrgRefController(AccomplishmentService)

	server := gin.Default()

	v1 := server.Group("/v1")

	v1.POST("/faq", middleware.AuthMiddleware(), FaqController.Create)
	v1.GET("/faq", FaqController.GetAll)
	v1.GET("/faq/:id", FaqController.GetById)
	v1.PATCH("/faq/:id/activate", middleware.AuthMiddleware(), FaqController.Activate)
	v1.PATCH("/faq/:id/deactivate", middleware.AuthMiddleware(), FaqController.Deactivate)

	v1.POST("/settings", middleware.AuthMiddleware(), SettingsController.Create)
	v1.GET("/settings", SettingsController.GetAll)
	v1.PATCH("/settings/activate", middleware.AuthMiddleware(), SettingsController.Activate)
	v1.PATCH("/settings/deactivate", middleware.AuthMiddleware(), SettingsController.Deactivate)

	v1.POST("/student_file_type", middleware.AuthMiddleware(), StfController.Create)
	v1.GET("/student_file_type/type", StfController.GetType)
	v1.GET("/student_file_type", StfController.GetAll)
	v1.GET("/student_file_type/:id", StfController.GetById)
	v1.PATCH("/student_file_type/:id/activate", middleware.AuthMiddleware(), StfController.Activate)
	v1.PATCH("/student_file_type/:id/deactivate", middleware.AuthMiddleware(), StfController.Deactivate)

	v1.POST("/accomplishment", middleware.AuthMiddleware(), AccomplishmentController.CreateRef)
	v1.GET("/accomplishment", AccomplishmentController.GetAllRef)
	v1.PATCH("/accomplishment/:id/activate", middleware.AuthMiddleware(), AccomplishmentController.Activate)
	v1.PATCH("/accomplishment/:id/deactivate", middleware.AuthMiddleware(), AccomplishmentController.Deactivate)
	v1.PATCH("/accomplishment/:id", middleware.AuthMiddleware(), AccomplishmentController.UpdateRef)

	v1.POST("/accomplishment-org", middleware.AuthMiddleware(), AccomplishmentOrganizationController.CreateOrganizationRef)
	v1.GET("/accomplishment-org", AccomplishmentOrganizationController.GetAllOrganizationRef)
	v1.PATCH("/accomplishment-org/:id/activate", middleware.AuthMiddleware(), AccomplishmentOrganizationController.Activate)
	v1.PATCH("/accomplishment-org/:id/deactivate", middleware.AuthMiddleware(), AccomplishmentOrganizationController.Deactivate)
	v1.PATCH("/accomplishment-org/:id", middleware.AuthMiddleware(), AccomplishmentOrganizationController.UpdateRefOrg)
	v1.POST("/accomplishment/calculation/:type", controller.Calc)

	server.Run(":8083")
}

func jwtgen() {
	userID := uuid.MustParse("3fd97d6e-430b-4d65-9fc6-6a8d83d0f2cb")

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.LoadConfig().JwtToken))
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT Token:")
	fmt.Println(signed)
}
