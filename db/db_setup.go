package db

import (
	"fmt"
	"tapirus_lite/internal/domain/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBSetup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("tapirus_lite.db"), &gorm.Config{})
	if err != nil {
		panic("Error conectando a SQLite: " + err.Error())
	}
	db.AutoMigrate(&entities.Product{}, &entities.Client{}, &entities.Order{}, &entities.OrderItem{})

	var count int64
	db.Model(&entities.Client{}).Count(&count)
	if count == 0 {
		defaultClient := entities.Client{
			Name: "Consumidor final",
		}
		if err := db.Create(&defaultClient).Error; err != nil {
			panic("failed to create default client")
		}
		fmt.Println("Cliente 'Consumidor final' creado con ID:", defaultClient.ID)
	}

	return db
}
