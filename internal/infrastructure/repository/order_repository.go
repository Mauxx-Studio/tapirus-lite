package repository

import (
	"tapirus_lite/internal/domain/entities"

	"gorm.io/gorm"
)

type OrderRpository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRpository {
	return &OrderRpository{db: db}
}

func (r *OrderRpository) FindAll(orders *[]entities.Order) error {
	return r.db.Preload("Items").Find(orders).Error
}

func (r *OrderRpository) FindByID(id uint) (entities.Order, error) {
	var order entities.Order
	err := r.db.Preload("Items").First(&order, id).Error
	return order, err
}

func (r *OrderRpository) Create(order *entities.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRpository) Update(order *entities.Order) error {
	//Eliminar Items existentes
	if err := r.db.Where("order_id = ?", order.ID).Delete(&entities.OrderItem{}).Error; err != nil {
		return err
	}
	//Guardar con nuevos items
	return r.db.Save(order).Error
}

func (r *OrderRpository) Delete(order *entities.Order) error {
	return r.db.Delete(order).Error
}
