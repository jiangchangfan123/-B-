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

// 初始化AuthService，其实也相当于dao/users.go中可以操控数据库的实例
func NewAuthService(userDao *dao.UserDao) *AuthService {
	return &AuthService{userDao: userDao}
}

// 注册用户的函数
func (s *AuthService) Register(username, password, email string) (*model.User, error) {
	//先查看是否已经有重复的用户名，如果有的话就返回具体的Error
	exists, err := s.userDao.ExistsUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	//对用户输入的密码进行加密
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	//将注册的用户进行实例化
	user := &model.User{
		Username: username,
		Password: hash,
		Email:    email,
		Role:     1,
		Status:   1,
		Coins:    0,
	}
	//将实例化后的用户存到数据库中
	if err := s.userDao.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// 用户登录的函数
func (s *AuthService) Login(username, password string) (*model.User, error) {
	//根据用户输入的用户名去数据库中查找
	user, err := s.userDao.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	//查看此用户是不是被封了
	if user.Status != 1 {
		return nil, errors.New("账号被禁用")
	}
	//检查用户的密码是否输入正确
	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// 根据用户的Id得到该用户的函数
func (s *AuthService) GetUserByID(id uint64) (*model.User, error) {
	return s.userDao.GetByID(id)
}
