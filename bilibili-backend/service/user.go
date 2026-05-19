package service

import (
	"errors"

	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (s *UserService) GetUserByID(id uint64) (*model.User, error) {
	return s.userDao.GetByID(id)
}

func (s *UserService) UpdateProfile(userID uint64, sign, nickname string) error {
	user, err := s.userDao.GetByID(userID)
	if err != nil {
		return err
	}
	if sign != "" {
		user.Sign = sign
	}
	if nickname != "" {
		user.Nickname = nickname
	}
	return s.userDao.Update(user)
}

func (s *UserService) UpdatePassword(userID uint64, oldPwd, newPwd string) error {
	user, err := s.userDao.GetByID(userID)
	if err != nil {
		return err
	}
	if !utils.CheckPassword(oldPwd, user.Password) {
		return errors.New("旧密码错误")
	}
	hash, err := utils.HashPassword(newPwd)
	if err != nil {
		return err
	}
	user.Password = hash
	return s.userDao.Update(user)
}

func (s *UserService) UpdateAvatar(userID uint64, avatarURL string) error {
	user, err := s.userDao.GetByID(userID)
	if err != nil {
		return err
	}
	user.Avatar = avatarURL
	return s.userDao.Update(user)
}
