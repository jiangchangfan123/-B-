package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) Create(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *UserDao) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := d.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) ExistsUsername(username string) (bool, error) {
	var count int64
	err := d.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
