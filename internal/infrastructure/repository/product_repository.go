package repository

import (
	"tapirus_lite/internal/domain/entities"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *entities.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *entities.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(product *entities.Product) error {
	return r.db.Delete(product).Error
}

func (r *ProductRepository) FindAll(products *[]entities.Product) error {
	return r.db.Find(products).Error
}

func (r *ProductRepository) FindByID(id uint) (entities.Product, error) {
	var product entities.Product
	err := r.db.First(&product, id).Error
	return product, err
}
