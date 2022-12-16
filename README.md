# Project description
簡易的備忘錄專案

# Project plugins
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [MySQL](https://github.com/go-gorm/mysql)
- [GoDotEnv](https://github.com/joho/godotenv)
- [crypto](https://pkg.go.dev/golang.org/x/crypto)
- [smapping](https://github.com/mashingan/smapping)

# How to build
## Install migrate
> https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

## Generate vendor and migrations
```bash
go mod vendor

# Up migration
# migrate -database "mysql://{user}:{password}@tcp({host}:{port})/{project name}" -path ./migrations up
migrate -database "mysql://root:@tcp(127.0.0.1:3306)/go-todolist" -path ./migrations up

# Down migration
# migrate -database "mysql://{user}:{password}@tcp({host}:{port})/{project name}" -path ./migrations down
migrate -database "mysql://root:@tcp(127.0.0.1:3306)/go-todolist" -path ./migrations down
```

## Run go
```bash
go run main.go
```

# Folders structure
```
├── controller
│   └── userController.go
├── entity
│   └── userEntity.go
├── migration
│   ├── 20221129000000_create_users_table.down.sql
│   └── 20221129000000_create_users_table.up.sql
├── model
│   └── user.go
├── request
│   └── userRequest.go
├── router
│   └── api.go
├── services
│   └── userService.go
├── utils
│   ├── gorm
│   │   └── gorm.go
│   └── responses
│       └── response.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
├── main.go
└── README.md
```

# Folder definition
- Controller
> 接收 HTTP request 調用 request & service

- Request
> 輔助 controller 驗證 request parameters

- Service
> 輔助 controller 處理業務邏輯

- Entity
> 輔助 service 調用 sql query

- Model
> 作為返回的對象

- Migration
> 建立 datatable 詳細資訊

- Router
> API 路由位置

- Utils
> 放置模組化位置，提供整個專案調用
