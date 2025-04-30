package main

import (
	"spmb/back-end/api/internal/controller"
	"spmb/back-end/api/internal/middleware"
	"spmb/back-end/api/internal/service"
	"spmb/back-end/config/connection"

	"github.com/gin-gonic/gin"
)

func main() {

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
