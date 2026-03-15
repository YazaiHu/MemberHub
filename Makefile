.PHONY: help run build test clean migrate-up migrate-down docker-up docker-down

# 默认配置
CONFIG_PATH ?= configs/config.dev.yaml

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## 运行API服务
	CONFIG_PATH=$(CONFIG_PATH) go run cmd/api/main.go

worker: ## 运行Worker服务
	CONFIG_PATH=$(CONFIG_PATH) go run cmd/worker/main.go

run-all: ## 同时运行API和Worker服务
	@echo "Starting API and Worker services..."
	@CONFIG_PATH=$(CONFIG_PATH) go run cmd/api/main.go & CONFIG_PATH=$(CONFIG_PATH) go run cmd/worker/main.go

build: ## 编译项目
	@echo "Building API server..."
	@go build -o bin/api cmd/api/main.go
	@echo "Building Worker server..."
	@go build -o bin/worker cmd/worker/main.go
	@echo "Build complete!"

test: ## 运行测试
	go test -v -cover ./...

clean: ## 清理构建文件
	rm -rf bin/
	rm -rf logs/
	go clean

deps: ## 下载依赖
	go mod download
	go mod tidy

migrate-up: ## 执行数据库迁移
	@echo "Running database migrations..."
	@mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS membership_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	@mysql -u root -p membership_dev < migrations/001_create_users_and_admins.up.sql
	@mysql -u root -p membership_dev < migrations/002_create_stores.up.sql
	@mysql -u root -p membership_dev < migrations/003_create_points.up.sql
	@mysql -u root -p membership_dev < migrations/004_create_recharge.up.sql
	@mysql -u root -p membership_dev < migrations/005_create_coupons.up.sql
	@mysql -u root -p membership_dev < migrations/006_create_promotions_and_logs.up.sql
	@echo "Migrations complete!"

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back migrations..."
	@mysql -u root -p -e "DROP DATABASE IF EXISTS membership_dev;"
	@echo "Rollback complete!"

docker-up: ## 启动Docker服务（MySQL, Redis, RabbitMQ）
	docker compose -f deployments/docker/docker-compose.yml up -d

docker-down: ## 停止Docker服务
	docker compose -f deployments/docker/docker-compose.yml down

docker-logs: ## 查看Docker日志
	docker compose -f deployments/docker/docker-compose.yml logs -f

fmt: ## 格式化代码
	go fmt ./...
	goimports -w .

lint: ## 代码检查
	golangci-lint run ./...

.DEFAULT_GOAL := help
