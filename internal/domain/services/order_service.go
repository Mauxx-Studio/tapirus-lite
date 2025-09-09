package services

import (
	"errors"
	"strings"
	"tapirus_lite/internal/domain/entities"
	"tapirus_lite/internal/infrastructure/repository"
	"time"
)

type OrderService struct {
	orderRepo      *repository.OrderRpository
	productService *ProductService
	clientService  *ClientService
}

func NewOrderService(orderRepo *repository.OrderRpository, productService *ProductService, clientService *ClientService) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		clientService:  clientService,
	}
}

func (s *OrderService) GetAllProducts() ([]entities.Product, error) {
	return s.productService.GetAllProducts()
}

func (s *OrderService) GetAllClients() ([]entities.Client, error) {
	return s.clientService.GetAllClients()
}

func (s *OrderService) GetAllOrders() ([]entities.Order, error) {
	var orders []entities.Order
	err := s.orderRepo.FindAll(&orders)
	return orders, err
}

func (s *OrderService) CreateOrder(clientID uint, items []entities.OrderItem, deliveryDate time.Time, note string) (*entities.Order, error) {
	// Validaciones
	if len(items) == 0 {
		return nil, errors.New("debe ingresar al menos un ítem")
	}
	if clientID == 0 {
		return nil, errors.New("debe indicar un cliente")
	}
	if deliveryDate.IsZero() {
		return nil, errors.New("fecha invalida")
	}

	// Validar Cliente
	client, err := s.clientService.GetByID(clientID)
	if err != nil {
		return nil, errors.New("cliente no encontrado")
	}

	// Validar ítems y calcular monto total
	var amount float64
	for i, item := range items {
		product, err := s.productService.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("producto no encontrado, " + err.Error())
		}
		if product.Stock < int(item.Quantity) {
			return nil, errors.New("stock insuficiente para: " + product.Name)
		}
		amount += product.Price * item.Quantity
		items[i].Product = product
	}

	// Crear orden
	order := &entities.Order{
		ClientID:     clientID,
		ClientName:   client.Name,
		DeliveryDate: deliveryDate,
		Note:         strings.TrimSpace(note),
		Completed:    false,
		Amount:       amount,
		Items:        items,
	}

	// Guardar orden en la base de datos
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Actualizar stock y contadores
	for _, item := range items {
		product, _ := s.productService.GetByID(item.ProductID)
		product.Stock -= int(item.Quantity)
		product.TotalSold += item.Quantity
		if err := s.productService.UpdateProduct(&product, product.Name, product.Description, product.Unit, product.Price, product.Stock); err != nil {
			return nil, err
		}
	}
	client.PendingOrders++
	client.TotalOrders++
	if err := s.clientService.UpdateClient(&client, client.Name, client.Phone, client.Email, client.CUIT, client.Address); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(order *entities.Order, clientID uint, items []entities.OrderItem, deliveryDate time.Time, note string) error {
	// Vaslidaciones
	if len(items) == 0 {
		return errors.New("debe ingresar al menos un prosucto")
	}
	if clientID == 0 {
		return errors.New("debe ingresar un cliente")
	}
	if deliveryDate.IsZero() {
		return errors.New("fecha inválida")
	}

	// Validar cliente
	client, err := s.clientService.GetByID(clientID)
	if err != nil {
		return errors.New("cliente no encontrado")
	}

	// Validar ítems y calcular moto total
	var amount float64
	for i, item := range items {
		product, err := s.productService.GetByID(item.ProductID)
		if err != nil {
			return errors.New("producto no encontrado, " + err.Error())
		}
		// Restaurar stock de ítems originales
		for _, oldItem := range order.Items {
			oldProduct, _ := s.productService.GetByID(oldItem.ProductID)
			oldProduct.Stock += int(oldItem.Quantity)
			s.productService.UpdateProduct(&oldProduct, oldProduct.Name, oldProduct.Description, oldProduct.Unit, oldProduct.Price, oldProduct.Stock)
		}
		// Verificar stock para nuevos ítems
		if product.Stock < int(item.Quantity) {
			return errors.New("stock insuficiente para " + product.Name)
		}
		amount += product.Price * item.Quantity
		items[i].Product = product // Asignar producto para relación
	}

	// Actualizar orden
	order.ClientID = clientID
	order.ClientName = client.Name
	order.DeliveryDate = deliveryDate
	order.Note = strings.TrimSpace(note)
	order.Amount = amount
	order.Items = items

	// Guardar cambios
	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	// Actualizar stock y contadores
	for _, item := range items {
		product, _ := s.productService.GetByID(item.ProductID)
		product.Stock -= int(item.Quantity)
		product.TotalSold += item.Quantity
		if err := s.productService.UpdateProduct(&product, product.Name, product.Description, product.Unit, product.Price, product.Stock); err != nil {
			return err
		}
	}

	// Actualizar contadores del cliente si cambió, creo que esto no va
	if order.ClientID != clientID {
		oldClient, _ := s.clientService.GetByID(order.ClientID)
		if oldClient.PendingOrders > 0 {
			oldClient.PendingOrders--
		}
		s.clientService.UpdateClient(&oldClient, oldClient.Name, oldClient.Phone, oldClient.Email, oldClient.CUIT, oldClient.Address)
		client.PendingOrders++
		client.TotalOrders++
		s.clientService.UpdateClient(&client, client.Name, client.Phone, client.Email, client.CUIT, client.Address)
	}
	return nil
}

func (s *OrderService) DeleteOrder(order *entities.Order) error {
	// Restaurar stock de los items
	for _, item := range order.Items {
		product, err := s.productService.GetByID(item.ProductID)
		if err != nil {
			return err
		}
		product.Stock += int(item.Quantity)
		product.TotalSold -= item.Quantity
		if err := s.productService.UpdateProduct(&product, product.Name, product.Description, product.Unit, product.Price, product.Stock); err != nil {
			return err
		}
	}

	// Actualizar contadores de clientes
	client, err := s.clientService.GetByID(order.ClientID)
	if err != nil {
		return err
	}
	if client.PendingOrders > 0 {
		client.PendingOrders--
	}
	if client.TotalOrders > 0 {
		client.TotalOrders++
	}
	if err := s.clientService.UpdateClient(&client, client.Name, client.Phone, client.Email, client.CUIT, client.Address); err != nil {
		return err
	}

	// Eliminar orden
	if err := s.orderRepo.Delete(order); err != nil {
		return err
	}

	return nil
}
