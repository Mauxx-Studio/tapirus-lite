package services

import (
	"errors"
	"strings"
	"tapirus_lite/internal/domain/entities"
	"tapirus_lite/internal/infrastructure/repository"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) CreateClient(name, phone, email, cuit, address string) (*entities.Client, error) {
	//Validacion
	name = strings.TrimSpace(name)
	if len(name) <= 2 {
		return nil, errors.New("el nombre debe tener más de dos letras")
	}

	//Crea Cliente
	client := &entities.Client{
		Name:    name,
		Phone:   phone,
		Email:   email,
		CUIT:    cuit,
		Address: address,
	}
	if err := s.repo.Create(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) UpdateClient(client *entities.Client, name, phone, email, cuit, address string) error {
	//Validacion
	name = strings.TrimSpace(name)
	if len(name) <= 2 {
		return errors.New("el nombre debe tener más de dos letras")
	}
	// Actualiza cliente
	client.Name = name
	client.Phone = phone
	client.Email = email
	client.CUIT = cuit
	client.Address = address

	return s.repo.Update(client)
}

func (s *ClientService) DeleteClient(client *entities.Client) error {
	return s.repo.Delete(client)
}

func (s *ClientService) GetAllClients() ([]entities.Client, error) {
	var clients []entities.Client
	err := s.repo.FindAll(&clients)
	return clients, err
}

func (s *ClientService) GetByID(id uint) (entities.Client, error) {
	return s.repo.FindByID(id)
}
