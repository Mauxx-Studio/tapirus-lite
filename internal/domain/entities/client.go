package entities

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
