package entities

import "time"

type Order struct {
	ID           uint        `gorm:"primaryKey"`
	ClientID     uint        // Relación con Client
	Client       Client      `gorm:"foreignKey:ClientID"`
	ClientName   string      // Nombre del cliente para consulta rápida
	DeliveryDate time.Time   `gorm:"not null"` // Fecha de entrega
	Note         string      // Campo de texto libre
	Completed    bool        `gorm:"default:false"` // Estado de cumplimiento
	Amount       float64     // Monto total del pedido
	Items        []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    // Relación con Order
	ProductID uint    // Relación con Product
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  float64 // Cambiado de int a float64
}
