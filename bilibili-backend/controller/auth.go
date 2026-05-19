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

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if !validateUsername(req.Username) {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if len(req.Password) < 6 || len(req.Password) > 30 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if !validateEmail(req.Email) {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	user, err := ctrl.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		if err.Error() == "用户名已存在" {
			utils.Error(c, utils.CodeUsernameExists)
			return
		}
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
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

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"accessToken": token,
		"userInfo": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"avatar":   user.Avatar,
			"role":     user.Role,
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
		"avatar":   user.Avatar,
		"role":     user.Role,
		"email":    user.Email,
		"sign":     user.Sign,
		"coins":    user.Coins,
	})
}

func validateUsername(u string) bool {
	if len(u) < 3 || len(u) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, u)
	return matched
}

func validateEmail(e string) bool {
	matched, _ := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, e)
	return matched
}
