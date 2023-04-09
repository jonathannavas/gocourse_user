package user

import (
	"fmt"
	"log"

	"strings"

	"github.com/jonathannavas/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *domain.User) error {

	// user.ID = uuid.New().String()

	// result := repo.db.Debug().Create(user)
	// if result.Error != nil {
	// 	repo.log.Println(result.Error)
	// 	return result.Error
	// }

	if err := repo.db.Debug().Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("User created with id: ", user.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	var user []domain.User

	tx := repo.db.Model(&user)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *repo) Get(id string) (*domain.User, error) {
	user := domain.User{
		ID: id,
	}
	result := repo.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *repo) Delete(id string) error {
	user := domain.User{
		ID: id,
	}
	result := repo.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = firstName
	}

	if lastName != nil {
		values["last_name"] = lastName
	}

	if email != nil {
		values["email"] = email
	}

	if phone != nil {
		values["phone"] = phone
	}

	result := repo.db.Model(&domain.User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, nil
	}
	return int(count), nil
}
