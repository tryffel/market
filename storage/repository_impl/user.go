package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/market/storage/models"
	"github.com/tryffel/market/storage/repositories"
	"strings"
)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) Search(text string, userId string, paging repositories.Paging) (*[]models.User, error) {
	// TODO: limit to common groups
	lower := strings.ToLower(text)
	query := u.db.Where("id LIKE %?% or lower_name LIKE %?% or email LIKE %?%", lower)
	query = PagedSorted(query, paging)
	users := &[]models.User{}
	res := query.Find(&users)
	return users, getDatabaseError(res.Error)
}

func (u *UserRepository) Create(user *models.User, passwordHash string) error {
	user.Password = passwordHash
	res := u.db.Create(&user).Error
	return getDatabaseError(res)
}

func (u *UserRepository) Update(user *models.User) error {
	err := u.db.Save(&user).Error
	return getDatabaseError(err)
}

func (u *UserRepository) Delete(user *models.User) error {
	return u.db.Delete(&user).Error
}

func (u *UserRepository) FindById(id string) (*models.User, error) {
	user := &models.User{}
	res := u.db.Where("id = ?", id).First(&user)
	return user, getDatabaseError(res.Error)
}

func (u *UserRepository) FindByNid(id int) (*models.User, error) {
	user := &models.User{}
	res := u.db.Where("nid = ?", id).First(&user)
	return user, getDatabaseError(res.Error)
}

func (u *UserRepository) FindByName(name string) (*models.User, error) {
	user := &models.User{}
	res := u.db.Where("lower_name = ?", strings.ToLower(name)).First(&user)
	return user, getDatabaseError(res.Error)
}

func (u *UserRepository) UpdatePassword(user *models.User, passwordHash string) error {
	user.Password = passwordHash
	res := u.db.Create(&user).Error
	return getDatabaseError(res)
}

func NewUserRepository(db *gorm.DB) repositories.User {
	u := &UserRepository{
		db: db,
	}
	return u
}
