# AI Coding Agent Instructions for exam-project1

## 项目概述

这是一个基于Go后端和React前端的文件上传系统，使用阿里云OSS进行文件存储。系统采用前后端分离架构，通过STS临时凭证实现安全的文件上传。

## 项目结构

```
.
├── backgend/              # Go后端服务
│   ├── internal/         # 内部包
│   │   └── alioss/      # 阿里云OSS相关功能
│   └── main.go          # 主程序入口
└── frontend/             # React前端应用
    ├── src/             # 源代码
    │   ├── ossUploader.js  # OSS上传组件
    │   └── App.jsx      # 主应用组件
    └── vite.config.js   # Vite配置
```

## 关键概念

### 后端服务 (Go)
- 使用Gin框架构建RESTful API
- 通过环境变量进行配置（查看`.env`文件模板）
- 主要API端点：`/oss/sts-credentials`用于获取临时上传凭证

### 前端应用 (React + Vite)
- 基于Vite构建系统
- 使用ali-oss SDK处理文件上传
- 环境变量通过`import.meta.env.VITE_*`访问

## 开发工作流

### 后端开发
1. 配置`.env`文件（参考`backgend/main.go`中的必要环境变量）
2. 新增API时，遵循`internal/`包组织结构
3. 确保实现健康检查端点

### 前端开发
1. 使用`VITE_API_BASE_URL`环境变量配置API地址
2. 组件开发遵循`src/`目录结构
3. 文件上传相关逻辑集中在`ossUploader.js`

## 关键依赖

### 后端
- github.com/gin-gonic/gin
- github.com/joho/godotenv
- 阿里云SDK

### 前端
- ali-oss
- axios
- React + Vite

## 集成要点

1. OSS上传流程：
   - 前端请求STS凭证
   - 使用凭证初始化OSS客户端
   - 使用multipartUpload处理大文件上传

2. 环境变量：
   - 后端：使用`.env`文件
   - 前端：使用`VITE_`前缀的环境变量

3. 错误处理：
   - 后端统一使用Gin的错误处理中间件
   - 前端在`ossUploader.js`中集中处理上传错误

## 注意事项

1. 所有OSS相关配置必须通过环境变量提供
2. 前端API调用必须使用`API_BASE_URL`环境变量
3. 遵循目录结构约定，保持代码组织清晰