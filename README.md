# Project description
這是一個簡單的待辦事項專案 <br>
備註: <br>
1. jwt-token TTL 時間預設 900 秒，中介層驗證 token TTL 成功時，計算該 token TTL 是否低於 5 分鐘，低於則重新給予一組新的 token。
2. logout 僅回傳成功的資訊，須由前端開發者主動清除使用者 token 資訊 (未開發白名單或黑名單機制)。

It is a simple todo list project <br>
Note: <br>
1. jwt-token TTL time is preset to 900 seconds. When the middleware verifies the success of token TTL, it will calculate whether the token TTL is lower than 5 minutes, and if it is lower, a new set of token will be given again.
2. The logout only returns successful information, and the front-end developer must take the initiative to clear the user token information (no whitelist or blacklist mechanism has been developed).

# Project plugins
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [MySQL](https://github.com/go-gorm/mysql)
- [GoDotEnv](https://github.com/joho/godotenv)
- [crypto](https://pkg.go.dev/golang.org/x/crypto)
- [smapping](https://github.com/mashingan/smapping)
- [golang-jwt](https://github.com/golang-jwt/jwt)

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
├── middleware
│   └── jwt.go
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
│   ├── jwtService.go
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
> 接收 HTTP request 調用 request & service <br>
> Receiving HTTP requests calling requests and services

- Entity
> 協助 service 調用 sql query <br>
> Assist service in calling sql query

- Middleware
> 中介層，負責過濾進入的資料 <br>
> Intermediary layer, responsible for filtering incoming data

- Migration
> 建立 datatable 詳細資訊 <br>
> Create datatable details

- Model
> 作為返回的對象 <br>
> As a returned object

- Request
> 協助 controller 驗證 request parameters <br>
> Assist controller validation request parameters

- Router
> API 路由位置 <br>
> API route locations

- Service
> 協助 controller 處理業務邏輯 <br>
> Assist controller with business logic

- Utils
> 模組化程式碼置放處，提供專案調用 <br>
> Modular code placement for project calls
