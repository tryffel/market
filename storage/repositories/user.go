package repositories

import "github.com/tryffel/market/storage/models"

// User repository
type User interface {
	// CRUD
	Create(user *models.User, passwordHash string) error
	Update(user *models.User) error
	Delete(user *models.User) error
	// Find operations
	FindById(id string) (*models.User, error)
	FindByNid(id int) (*models.User, error)
	FindByName(name string) (*models.User, error)
	// Search user by any field. Limit search to groups that user is part of, if configured so
	Search(text string, userId string, paging Paging) (*[]models.User, error)

	UpdatePassword(user *models.User, passwordHash string) error
}
