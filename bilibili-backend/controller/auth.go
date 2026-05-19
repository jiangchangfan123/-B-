package controller

import (
	"regexp"

	"bilibili-backend/service"
	"bilibili-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

// 初始化得到authService实例，其实也是userDao中对数据库操作的对象
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// 注册时需要用到的字段为一个类型
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// 登录时用到的字段为一个类型
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 用户注册的函数
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterReq
	//检查参数是否符合要求，不符合的话返回参数错误的类型
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	//用户名字不合法就返回
	if !validateUsername(req.Username) {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	//检查密码的格式
	if len(req.Password) < 6 || len(req.Password) > 30 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	//检查邮箱的格式
	if !validateEmail(req.Email) {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	//一切正常后进行用户注册的服务
	user, err := ctrl.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		if err.Error() == "用户名已存在" {
			utils.Error(c, utils.CodeUsernameExists)
			return
		}
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	//成功后返回用户的id和用户名
	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// 用户登录的函数
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	//用户名或者密码不能为空
	if req.Username == "" || req.Password == "" {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	user, err := ctrl.authService.Login(req.Username, req.Password)
	if err != nil {
		if err.Error() == "账号被禁用" {
			utils.Error(c, utils.CodeUserForbidden)
			return
		}
		utils.Error(c, utils.CodeUserNotFound)
		return
	}

	//给用户创建Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	//返回成功的状态码
	utils.OK(c, gin.H{
		"accessToken": token,
		"userInfo": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
			"email":    user.Email,
			"sign":     user.Sign,
			"coins":    user.Coins,
		},
	})
}

func (ctrl *AuthController) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.CodeUnauthorized)
		return
	}

	user, err := ctrl.authService.GetUserByID(userID.(uint64))
	if err != nil {
		utils.Error(c, utils.CodeTokenInvalid)
		return
	}

	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"role":     user.Role,
		"email":    user.Email,
		"sign":     user.Sign,
		"coins":    user.Coins,
	})
}

// 检查用户名是否合法
func validateUsername(u string) bool {
	if len(u) < 3 || len(u) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, u)
	return matched
}

// 检查邮箱格式是否合法
func validateEmail(e string) bool {
	matched, _ := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, e)
	return matched
}
