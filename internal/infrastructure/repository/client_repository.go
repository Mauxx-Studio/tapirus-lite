package repository

import (
	"tapirus_lite/internal/domain/entities"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(client *entities.Client) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) Update(client *entities.Client) error {
	return r.db.Save(client).Error
}

func (r *ClientRepository) Delete(client *entities.Client) error {
	return r.db.Delete(client).Error
}

func (r *ClientRepository) FindAll(clients *[]entities.Client) error {
	return r.db.Find(clients).Error
}

func (r *ClientRepository) FindByID(id uint) (entities.Client, error) {
	var client entities.Client
	err := r.db.First(&client, id).Error
	return client, err
}
