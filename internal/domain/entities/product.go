package entities

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"` // Nombre del producto
	Description string  // Descripción del producto
	Price       float64 `gorm:"not null"` // Precio del producto
	Stock       int     `gorm:"not null"` // Cantidad en stock
	Unit        string  // Ejemplo: "kg", "un", "lts"
	TotalSold   float64 // Acumulado de ventas
}
