package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/auth"
	"github.com/tryffel/market/storage/models"
	"github.com/tryffel/market/storage/repositories"
	"strings"
)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) Create(user *models.User, password string) error {
	hash, err := auth.GetPasswordHash(password)
	if err != nil {
		return Error.Wrap(&err, "failed to hash new password")
	}
	user.Password = hash
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

func (u *UserRepository) UpdatePassword(user *models.User, password string) error {
	hash, err := auth.GetPasswordHash(password)
	if err != nil {
		return Error.Wrap(&err, "failed to hash new password")
	}
	user.Password = hash
	res := u.db.Create(&user).Error
	return getDatabaseError(res)
}

func NewUserRepository(db *gorm.DB) repositories.User {
	u := &UserRepository{
		db: db,
	}
	return u
}
