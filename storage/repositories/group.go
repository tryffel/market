package repositories

import "github.com/tryffel/market/storage/models"

type Group interface {
	Create(group *models.Group) error
	Remove(group *models.Group) error
	Update(group *models.Group) error
	FindById(id string) (*models.Group, error)
	UserBelongsToGroup(user int, group int) (bool, error)
	AddUserToGroup(user int, group int) error
	DelUserFromGroup(user int, group int) error
	ListUsers(group int, paging Paging) (*[]models.User, error)
}
