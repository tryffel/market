package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/market/storage/models"
	"github.com/tryffel/market/storage/repositories"
)

type GroupRepository struct {
	db *gorm.DB
}

func (g *GroupRepository) Create(group *models.Group) error {
	res := g.db.Create(&group)
	return getDatabaseError(res.Error)
}

func (g *GroupRepository) Remove(group *models.Group) error {
	panic("implement me")
}

func (g *GroupRepository) Update(group *models.Group) error {
	err := g.db.Save(&group).Error
	return getDatabaseError(err)
}

func (g *GroupRepository) FindById(id string) (*models.Group, error) {
	group := &models.Group{}
	err := g.db.Where("id = ?", id).First(&group).Error
	return group, getDatabaseError(err)
}

func (g *GroupRepository) UserBelongsToGroup(user int, group int) (bool, error) {
	panic("implement me")
}

func (g *GroupRepository) AddUserToGroup(user int, group int) error {
	ug := &models.UserGroup{
		UserNid:  user,
		GroupNid: group,
	}

	res := g.db.Create(&ug)
	return getDatabaseError(res.Error)

}

func (g *GroupRepository) DelUserFromGroup(user int, group int) error {
	panic("implement me")
}

func (g *GroupRepository) ListUsers(group int, paging repositories.Paging) (*[]models.User, error) {
	panic("implement me")
}

func NewGroupRepository(db *gorm.DB) repositories.Group {
	return &GroupRepository{db: db}
}
