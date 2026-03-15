# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**MemberHub (会员通)** is a unified membership management system for chain stores, built entirely using Claude Code (claude-sonnet-4.5) as a Vibe Coding showcase project.

- **Module Name**: `github.com/YazaiHu/MemberHub`
- **Backend**: Go 1.25+ with Gin framework
- **Frontend**: Vue 3 + Element Plus + Pinia
- **Databases**: MySQL 8.0+, Redis 7.0+
- **Message Queue**: RabbitMQ 3.12+
- **Architecture**: DDD (Domain-Driven Design) with Clean Architecture

## Development Commands

### Backend

```bash
# Run API server (default: configs/config.dev.yaml)
make run

# Run Worker service (for async tasks)
make worker

# Run both API and Worker
make run-all

# Build binaries
make build

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Clean build artifacts
make clean

# Download dependencies
make deps
```

### Frontend (web/admin)

```bash
cd web/admin

# Install dependencies
npm install

# Run dev server (http://localhost:3000)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Database

```bash
# Run all migrations (creates membership_dev database)
make migrate-up

# Drop database
make migrate-down
```

### Docker Services

```bash
# Start MySQL, Redis, RabbitMQ
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs
```

### Testing

```bash
# Complete system API test
./test_system.sh

# Frontend-backend integration test
./test_frontend_api.sh
```

## Architecture

### DDD Layered Structure

```
internal/
├── domain/              # Domain Layer (Business Logic)
│   ├── member/         # Each domain has:
│   │   ├── model.go       # - Domain models (entities)
│   │   ├── repository.go  # - Repository interface
│   │   └── service.go     # - Domain services
│   ├── points/
│   ├── recharge/
│   ├── coupon/
│   ├── store/
│   └── promotion/
│
├── application/        # Application Layer (Use Cases)
│   └── */service.go   # Application services orchestrating domain logic
│
├── infrastructure/     # Infrastructure Layer
│   ├── persistence/   # Data access implementations
│   │   ├── mysql/    # MySQL repositories
│   │   └── redis/    # Redis cache & locks
│   ├── messaging/    # RabbitMQ message queue
│   └── wechat/       # WeChat SDK integration
│
├── interfaces/        # Interface Layer (Presentation)
│   ├── http/
│   │   ├── handler/
│   │   │   ├── admin/    # Admin API handlers
│   │   │   └── miniapp/  # Mini-program API handlers
│   │   └── middleware/   # Auth, RBAC, CORS, etc.
│   └── worker/       # Background job consumers
│
└── pkg/              # Shared packages
    ├── config/       # Configuration management (Viper)
    ├── logger/       # Structured logging (Zap)
    ├── errors/       # Custom error types
    └── utils/        # Utilities
```

### Request Flow

1. **HTTP Request** → Middleware (Logger, CORS, Auth, RBAC) → Handler
2. **Handler** → Application Service → Domain Service → Repository
3. **Repository** → Infrastructure (MySQL/Redis) → Return Data
4. **Response** → JSON serialization → Client

### Key Patterns

#### 1. Repository Pattern
Each domain has a repository interface defined in `domain/*/repository.go` and implemented in `infrastructure/persistence/mysql/*_repository.go`.

#### 2. Service Pattern
- **Domain Services** (`domain/*/service.go`): Core business logic, no dependencies on infrastructure
- **Application Services** (`application/*/service.go`): Orchestrate multiple domain services, handle transactions

#### 3. Dependency Injection
Handlers instantiate services with dependencies:
```go
// In cmd/api/main.go or handler constructors
pointsRepo := mysql.NewPointsRepository(db)
pointsService := domain.NewPointsService(pointsRepo, redisClient)
handler := admin.NewPointsHandler(pointsService)
```

## Critical Business Logic

### High-Concurrency Points Exchange
**Location**: `internal/domain/points/service.go`

Uses triple-layer protection against overselling:
1. **Redis atomic operations** for stock pre-deduction
2. **Distributed locks** to prevent duplicate submissions
3. **Optimistic locking** (version field) in database

When modifying points exchange logic, ensure all three mechanisms remain in place.

### Recharge Transaction Consistency
**Location**: `internal/domain/recharge/transaction.go`

Payment callback processing must maintain atomicity:
- Order status update
- Balance increase (with optimistic lock)
- Bonus points grant
- Transaction log creation

All operations wrapped in database transaction with idempotency check (transaction_no).

### Async Coupon Push
**Location**: `internal/interfaces/worker/coupon_push_worker.go`

- Uses RabbitMQ for async processing
- Batch processing (100 users per batch)
- Task status tracking (pending → running → completed)
- Failure retry mechanism

## Configuration

### Environment-Specific Config
- **Development**: `configs/config.dev.yaml`
- **Production**: `configs/config.prod.yaml` (not in repo)

Override config path:
```bash
CONFIG_PATH=configs/config.prod.yaml make run
```

### Required Services

Before running, ensure these services are available:
- **MySQL** (127.0.0.1:3306) - Database: `membership_dev`
- **Redis** (127.0.0.1:6379) - For cache and distributed locks
- **RabbitMQ** (127.0.0.1:5672) - For async tasks

Quick start with Docker:
```bash
make docker-up
make migrate-up
make run
```

### Default Credentials

**Admin Login**:
- Username: `admin`
- Password: `admin123`
- API: `POST /api/admin/auth/login`

If login fails, `test_system.sh` auto-resets the password hash.

## API Structure

### Admin API (`/api/admin`)
Requires JWT authentication via `Authorization: Bearer <token>` header.

**Endpoints**:
- `/auth/login` - Admin login
- `/members/*` - Member management, manual points/balance adjustment
- `/stores/*` - Store CRUD
- `/points/*` - Points exchange rules, redemption records
- `/recharge/*` - Recharge promotions, order management
- `/coupons/*` - Coupon templates, push tasks
- `/promotions/*` - Special products management

### Mini-Program API (`/api/miniapp`)
Requires WeChat user JWT authentication.

**Endpoints**:
- `/auth/login` - WeChat code login
- `/member/*` - User profile, assets overview
- `/points/*` - Points balance, history, exchange
- `/recharge/*` - Recharge orders, balance query
- `/coupons/*` - Available coupons, receive
- `/stores/*` - Store list, nearby stores
- `/promotions/*` - Current special products

## Database Migrations

Migration files in `migrations/` directory:
- `001_create_users_and_admins.up.sql` - User/admin tables
- `002_create_stores.up.sql` - Store tables
- `003_create_points.up.sql` - Points system (4 tables)
- `004_create_recharge.up.sql` - Recharge system (4 tables)
- `005_create_coupons.up.sql` - Coupon system (3 tables)
- `006_create_promotions_and_logs.up.sql` - Promotions and logs

**Important**: Always apply migrations in order. Each migration builds on previous ones.

## Testing Strategy

### Unit Tests
```bash
# Run all tests
make test

# Run tests for specific package
go test -v ./internal/domain/points/...

# Run with coverage
go test -cover ./...
```

### Integration Tests
Use provided shell scripts:
- `test_system.sh` - Tests all 35 API endpoints
- `test_frontend_api.sh` - Tests frontend-backend integration

### Manual Testing
1. Start backend: `make run`
2. Start frontend: `cd web/admin && npm run dev`
3. Access: http://localhost:3000
4. Login with admin/admin123

## Important Notes

### Module Import Paths
All internal imports use: `github.com/YazaiHu/MemberHub/internal/...`

When adding new packages, follow this pattern:
```go
import (
    "github.com/YazaiHu/MemberHub/internal/domain/points"
    "github.com/YazaiHu/MemberHub/internal/pkg/logger"
)
```

### Error Handling
Use custom error types from `internal/pkg/errors/`:
- `ErrNotFound` - Resource not found (404)
- `ErrUnauthorized` - Auth failed (401)
- `ErrForbidden` - Permission denied (403)
- `ErrBadRequest` - Invalid input (400)
- `ErrConflict` - Duplicate/version conflict (409)

### Logging
Use structured logging from `internal/pkg/logger/`:
```go
logger.Info("Processing request", zap.String("user_id", userID))
logger.Error("Failed to save", zap.Error(err))
```

### Transaction Management
For operations requiring multiple database writes:
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// ... perform operations with tx ...

if err := tx.Commit().Error; err != nil {
    return err
}
```

## AI-Generated Codebase

This project was built 100% with Claude Code as a Vibe Coding demonstration:
- **Development Time**: ~8 hours (vs 3-4 weeks traditional)
- **Code Lines**: 29,300+ lines
- **Files**: 140+ files
- **Documentation**: 15+ technical documents

When making changes:
- Follow existing patterns and conventions
- Maintain DDD architecture boundaries
- Keep business logic in domain layer
- Preserve the three-layer concurrency protection in points system
- Ensure transaction atomicity in financial operations

Detailed development notes: `docs/AI_DEVELOPMENT_NOTES.md`
