# 阶段 1: 编译
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
RUN npm install
COPY . .
# 运行打包
RUN npm run build

# 阶段 2: 运行 (用 Nginx 托管)
FROM nginx:alpine
# 复制你打包好的 dist 文件夹到 Nginx 的网页根目录
COPY --from=builder /app/dist /usr/share/nginx/html
# (可选) 如果你有自己的 Nginx 配置，可以在这里 COPY

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]