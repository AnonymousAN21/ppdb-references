package connection

import (
	"log"
	"spmb/back-end/config"
	"spmb/back-end/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Conn() *gorm.DB {
	dsn := config.LoadConfig().Database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = db.AutoMigrate(
		&models.FAQ{},
		&models.Settings{},
		&models.StudentFileType{},
		&models.AccomplishmentRef{},
		&models.AccomplishmentOrganizationRef{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return db
}
