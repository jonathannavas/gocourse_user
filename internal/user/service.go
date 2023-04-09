package user

import (
	"log"

	"github.com/jonathannavas/gocourse_domain/domain"
)

type (
	Service interface {
		Create(firstName string, lastName string, email string, phone string) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Get(id string) (*domain.User, error)
		Delete(id string) error
		Update(id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		FirstName string
		LastName  string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName string, lastName string, email string, phone string) (*domain.User, error) {
	log.Println("Create user service")

	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	// si recibo un puntero debo enviar la dirección de memoria a donde esta con el simbolo de &

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	log.Println("Service Get All")
	users, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) Get(id string) (*domain.User, error) {
	log.Println("Get user service by id:", id)
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) Delete(id string) error {
	log.Println("Delete user service by id:", id)
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	return s.repo.Update(id, firstName, lastName, email, phone)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
