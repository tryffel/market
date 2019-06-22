package repositories

import "github.com/tryffel/market/storage/models"

// User repository
type User interface {
	Create(user *models.User, password string) error
	Update(user *models.User) error
	Delete(user *models.User) error
	FindById(id string) (*models.User, error)
	FindByNid(id int) (*models.User, error)
	FindByName(name string) (*models.User, error)
	UpdatePassword(user *models.User, password string) error
}
