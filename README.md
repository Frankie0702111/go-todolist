# Project description
這是一個簡單的待辦事項專案
It is a simple todo list project

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
> Receiving HTTP requests calling requests and services

- Request
> 協助 controller 驗證 request parameters
> Assist controller validation request parameters

- Service
> 協助 controller 處理業務邏輯
> Assist controller with business logic

- Entity
> 協助 service 調用 sql query
> Assist service in calling sql query

- Model
> 作為返回的對象
> As a returned object

- Migration
> 建立 datatable 詳細資訊
> Create datatable details

- Router
> API 路由位置
> API route locations

- Utils
> 模組化程式碼置放處，提供專案調用
> Modular code placement for project calls