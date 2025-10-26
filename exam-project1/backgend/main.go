package main

import (
	"log"
	"net/http" // 需要导入 net/http (虽然 Gin 封装了，但好习惯)
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// --- 导入你新创建的包 ---
	// 使用我们约定的通用模块名 "my-backend-app"
	"my-backend-app/internal/alioss"
)

// --- 全局配置变量 ---
var ossConfig alioss.Config // 使用 alioss 包定义的 Config 结构体

// --- 配置加载函数 (保持不变) ---
func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("警告：无法加载 .env 文件:", err)
	}

	// 加载配置到全局变量 ossConfig
	ossConfig = alioss.Config{
		StsEndpoint:     os.Getenv("STS_ENDPOINT"),
		AccessKeyID:     os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"),
		BucketName:      os.Getenv("OSS_BUCKET_NAME"),
		OssRegion:       os.Getenv("OSS_REGION"),
		RoleArn:         os.Getenv("ALIBABA_CLOUD_RAM_ROLE_ARN"), // 可选
	}

	if ossConfig.AccessKeyID == "" || ossConfig.AccessKeySecret == "" || ossConfig.BucketName == "" || ossConfig.OssRegion == "" || ossConfig.StsEndpoint == "" {
		log.Fatal("错误：必要的 OSS/STS 环境变量 (ACCESS_KEY, SECRET_KEY, BUCKET_NAME, OSS_REGION, STS_ENDPOINT) 未设置！请检查 .env 文件或系统环境变量。")
	}
	log.Println("配置加载成功.")
}

func main() {
	loadConfig() // 首先加载配置

	router := gin.Default()

	// --- 注册健康检查路由 (可选，但推荐) ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --- 在这里添加你【考核项目】的路由注册 ---
	// 例如:
	// userHandler := handlers.NewUserHandler(db) // 假设你有 handlers 包和数据库连接 db
	// router.GET("/users", userHandler.GetUsers)
	// router.POST("/users", userHandler.CreateUser)
	// -----------------------------------------

	// --- 调用 alioss 包来注册 OSS 相关的路由 ---
	// 把加载好的配置传递给它
	alioss.RegisterOssRoutes(router, ossConfig)

	// --- 获取监听端口 (可选，从环境变量获取更灵活) ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默认端口
	}
    listenAddr := ":" + port

	log.Printf("启动 API 服务器，监听在 %s", listenAddr)
	// router.Run() 默认监听 8080，改为 Run(listenAddr) 更灵活
	if err := router.Run(listenAddr); err != nil {
        log.Fatalf("启动 Gin 服务器失败: %v", err)
    }
}

// --- 在这里添加你【考核项目】的业务处理函数 (Handlers) ---
// 例如:
// func (h *UserHandler) GetUsers(c *gin.Context) { ... }
// func (h *UserHandler) CreateUser(c *gin.Context) { ... }
// -------------------------------------------------------