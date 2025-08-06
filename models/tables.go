package models

import "time"

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"` // Nombre del producto
	Description string  // Descripción del producto
	Price       float64 `gorm:"not null"` // Precio del producto
	Stock       int     `gorm:"not null"` // Cantidad en stock
	Unit        string  // Ejemplo: "kg", "un", "lts"
	TotalSold   float64 // Acumulado de ventas
}

type Client struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"` // Nombre del cliente
	Phone         string // Opcional, útil para WhatsApp
	Email         string // Opcional
	CUIT          string // Opcional
	Address       string // Opcional
	TotalOrders   int    // Total de pedidos realizados
	PendingOrders int    // Pedidos pendientes
}

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
