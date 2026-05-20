package main

import (
	"fmt"
	"log"
	"os/exec"

	"bilibili-backend/config"
	"bilibili-backend/consumer"
	"bilibili-backend/controller"
	"bilibili-backend/dao"
	"bilibili-backend/middleware"
	"bilibili-backend/model"
	"bilibili-backend/router"
	"bilibili-backend/service"
	"bilibili-backend/utils"

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
	if err := db.AutoMigrate(&model.User{}, &model.Video{}, &model.VideoHistory{}); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	// MinIO 初始化
	if err := utils.InitMinIO(); err != nil {
		log.Printf("[WARN] MinIO 初始化失败: %v", err)
		log.Println("[WARN] 文件上传功能将不可用")
	}

	// RabbitMQ 初始化
	if err := utils.InitRabbitMQ(); err != nil {
		log.Printf("[WARN] RabbitMQ 初始化失败: %v", err)
		log.Println("[WARN] 异步转码功能将不可用")
	} else {
		// 启动转码消费者
		go consumer.StartTranscodeConsumer(db)
	}

	// FFmpeg 环境检查
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Println("[WARN] FFmpeg 未安装，视频转码和封面截取功能将不可用")
	} else {
		log.Println("[FFmpeg] 检测到 ffmpeg")
	}

	// DAO
	userDao := dao.NewUserDao(db)
	videoDao := dao.NewVideoDao(db)
	_ = dao.NewHistoryDao(db)

	// Service
	authService := service.NewAuthService(userDao)
	userService := service.NewUserService(userDao)
	videoService := service.NewVideoService(videoDao, userDao)

	// Controller
	authCtrl := controller.NewAuthController(authService)
	userCtrl := controller.NewUserController(userService)
	uploadCtrl := controller.NewUploadController(userService)
	videoCtrl := controller.NewVideoController(videoService, userService)

	// Gin
	gin.SetMode(config.C.Server.Mode)
	r := gin.Default()
	r.Use(middleware.CORS())

	router.Setup(r, authCtrl, userCtrl, uploadCtrl, videoCtrl)

	addr := ":" + config.C.Server.Port
	fmt.Println("🚀 Server running on http://localhost" + addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
