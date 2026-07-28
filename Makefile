.PHONY: dev infra start stop clean db-init db-migrate help

# ============================================
# AI-Growth-Engine Makefile
# ============================================

# 启动基础设施 (PostgreSQL + Redis)
infra:
	docker compose up -d postgres redis

# 停止基础设施
stop:
	docker compose down

# 清理基础设施 (含数据)
clean:
	docker compose down -v

# 数据库初始化
db-init:
	docker compose exec postgres psql -U postgres -d ai_growth_engine -f /docker-entrypoint-initdb.d/init.sql

# 启动全部开发服务
dev: infra

# 默认帮助
help:
	@echo "AI-Growth-Engine"
	@echo ""
	@echo "  make infra      启动基础设施 (PostgreSQL + Redis)"
	@echo "  make stop       停止基础设施"
	@echo "  make clean      清理基础设施 (含数据卷)"
	@echo "  make dev        启动全部开发服务"
