package main

import (
	"fmt"
	"log"

	"bilibili-backend/config"
	"bilibili-backend/controller"
	"bilibili-backend/dao"
	"bilibili-backend/middleware"
	"bilibili-backend/model"
	"bilibili-backend/router"
	"bilibili-backend/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.C.DB.Username, config.C.DB.Password, config.C.DB.Host,
		config.C.DB.Port, config.C.DB.Database, config.C.DB.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	// DAO / Service / Controller
	userDao := dao.NewUserDao(db)
	authService := service.NewAuthService(userDao)
	authCtrl := controller.NewAuthController(authService)

	// Gin
	gin.SetMode(config.C.Server.Mode)
	r := gin.Default()
	r.Use(middleware.CORS())

	router.Setup(r, authCtrl)

	addr := ":" + config.C.Server.Port
	fmt.Println("🚀 Server running on http://localhost" + addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
