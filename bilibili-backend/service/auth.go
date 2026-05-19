package service

import (
	"errors"

	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"
)

type AuthService struct {
	userDao *dao.UserDao
}

func NewAuthService(userDao *dao.UserDao) *AuthService {
	return &AuthService{userDao: userDao}
}

func (s *AuthService) Register(username, password, email string) (*model.User, error) {
	exists, err := s.userDao.ExistsUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: hash,
		Email:    email,
		Role:     1,
		Status:   1,
		Coins:    0,
	}
	if err := s.userDao.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	user, err := s.userDao.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在或密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账号被禁用")
	}
	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("用户不存在或密码错误")
	}
	return user, nil
}

func (s *AuthService) GetUserByID(id uint64) (*model.User, error) {
	return s.userDao.GetByID(id)
}
