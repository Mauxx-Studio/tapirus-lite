package db

import (
	"fmt"
	"tapirus_lite/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBSetup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("tapirus_lite.db"), &gorm.Config{})
	if err != nil {
		panic("Error conectando a SQLite: " + err.Error())
	}
	db.AutoMigrate(&models.Product{}, &models.Client{}, &models.Order{}, &models.OrderItem{})

	var count int64
	db.Model(&models.Client{}).Count(&count)
	if count == 0 {
		defaultClient := models.Client{
			Name: "Consumidor final",
		}
		if err := db.Create(&defaultClient).Error; err != nil {
			panic("failed to create default client")
		}
		fmt.Println("Cliente 'Consumidor final' creado con ID:", defaultClient.ID)
	}

	return db
}
