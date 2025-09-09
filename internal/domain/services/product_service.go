package services

import (
	"errors"
	"strings"
	"tapirus_lite/internal/domain/entities"
	"tapirus_lite/internal/infrastructure/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(name, description, unit string, price float64, stock int) (*entities.Product, error) {
	//Validaciones
	name = strings.TrimSpace(name) //elimina los espacios al principio y al final del string
	unit = strings.TrimSpace(unit)
	if len(name) <= 2 {
		return nil, errors.New("el nombre debe tener más de dos letras")
	}
	if unit == "" {
		return nil, errors.New("la unidad no puede estar vacía")
	}

	//Crear producto
	product := &entities.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Unit:        unit,
		TotalSold:   0,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(product *entities.Product, name, description, unit string, price float64, stock int) error {
	name = strings.TrimSpace(name)
	unit = strings.TrimSpace(unit)
	if len(name) <= 2 {
		return errors.New("el nombre debe tener más de dos letras")
	}
	if unit == "" {
		return errors.New("la unidad no debe estar vacía")
	}

	//Actualizar producto
	product.Name = name
	product.Description = strings.TrimSpace(description)
	product.Unit = unit
	product.Price = price
	product.Stock = stock

	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(product *entities.Product) error {
	return s.repo.Delete(product)
}

func (s *ProductService) GetAllProducts() ([]entities.Product, error) {
	var products []entities.Product
	err := s.repo.FindAll(&products)
	return products, err
}

func (s *ProductService) GetByID(id uint) (entities.Product, error) {
	return s.repo.FindByID(id)
}
