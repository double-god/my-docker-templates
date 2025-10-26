package alioss // 包名必须和文件夹名一致

import (
	"fmt"
	"net/http"

	// 需要 os 来读取环境变量
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
)

// Config 结构体用于持有从 main 包传递过来的配置
// 或者，你也可以让这个包直接读取环境变量，但从 main 传递更清晰
type Config struct {
	StsEndpoint     string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	OssRegion       string
	RoleArn         string // 可选
}

// OssHandler 结构体，持有配置信息
type OssHandler struct {
	Config Config
}

// NewOssHandler 创建一个新的 OssHandler 实例
func NewOssHandler(cfg Config) *OssHandler {
	return &OssHandler{Config: cfg}
}

// GetOssStsCredentials 是处理获取 STS 凭证的 Gin Handler
// 注意：它现在是 OssHandler 的一个方法，可以访问 h.Config
func (h *OssHandler) GetOssStsCredentials(c *gin.Context) {
	// 1. 创建 STS 客户端 (使用传入的配置)
	client, err := sts.NewClientWithAccessKey(h.Config.OssRegion, h.Config.AccessKeyID, h.Config.AccessKeySecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建 STS Client 失败: %v", err)})
		return
	}

	// 2. 创建 AssumeRole 请求
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.DurationSeconds = requests.NewInteger(3600) // 1 小时有效期

	// --- 【简化版：使用 Policy 限制权限】 ---
	policy := fmt.Sprintf(`{
        "Statement": [
            {
                "Action": ["oss:PutObject"],
                "Effect": "Allow",
                "Resource": ["acs:oss:*:*:%s/uploads/*"]
            }
        ],
        "Version": "1"
    }`, h.Config.BucketName)
	request.Policy = policy
	// 如果使用 RAM Role，设置 request.RoleArn = h.Config.RoleArn 和 session name

	// 3. 发起请求
	response, err := client.AssumeRole(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AssumeRole 失败: %v", err)})
		return
	}

	// 4. 返回给前端
	c.JSON(http.StatusOK, gin.H{
		"accessKeyId":     response.Credentials.AccessKeyId,
		"accessKeySecret": response.Credentials.AccessKeySecret,
		"stsToken":        response.Credentials.SecurityToken,
		"expiration":      response.Credentials.Expiration,
		"bucket":          h.Config.BucketName,
		"region":          h.Config.OssRegion,
	})
}

// RegisterOssRoutes 函数用于在 Gin 引擎上注册 OSS 相关的路由
// 它接收 Gin 引擎和配置作为参数
func RegisterOssRoutes(router *gin.Engine, cfg Config) {
	handler := NewOssHandler(cfg) // 创建 Handler 实例
	// 注册路由，将请求交给 handler 的方法处理
	router.GET("/oss/sts-credentials", handler.GetOssStsCredentials)
}
