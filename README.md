# Project description
這是一個簡單的待辦事項專案 <br>
備註: <br>
1. jwt-token TTL 時間預設 900 秒，中介層驗證 token TTL 成功時，計算該 token TTL 是否低於 5 分鐘，低於則重新給予一組新的 token。
2. token 加入白名單機制，經由 redis 管理，使用者登出一併清除 redis token。
3. 日誌輸出，預設採用日期分割檔案，如果需要按容量大小分割檔案請取消註解 `utils/log/logBySize.go`。

It is a simple todo list project <br>
Note: <br>
1. jwt-token TTL time is preset to 900 seconds. When the middleware verifies the success of token TTL, it will calculate whether the token TTL is lower than 5 minutes, and if it is lower, a new set of token will be given again.
2. The token is added to the whitelist mechanism, managed by redis, and the user logout to clear the redis token as well.
3. Log output, default date split file, if you need to split the file by size please uncomment `utils/log/logBySize.go`.

# Project plugins
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [MySQL](https://github.com/go-gorm/mysql)
- [GoDotEnv](https://github.com/joho/godotenv)
- [crypto](https://pkg.go.dev/golang.org/x/crypto)
- [smapping](https://github.com/mashingan/smapping)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [go-redis](https://github.com/go-redis/redis)
- [zap](https://github.com/uber-go/zap)
- [lumberjack](https://github.com/natefinch/lumberjack)
- [file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)
- [oauth2](https://github.com/golang/oauth2)
- [gjson](https://github.com/tidwall/gjson)

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

# Specify batch up or down (If you want to go down to a specific file, it is recommended to open a new folder)
# migrate -database "mysql://{user}:{password}@tcp({host}:{port})/{project name}" -path ./migrations up {number}
migrate -database "mysql://root:@tcp(127.0.0.1:3306)/go-todolist" -path ./migrations up 1
# migrate -database "mysql://{user}:{password}@tcp({host}:{port})/{project name}" -path ./migrations down {number}
migrate -database "mysql://root:@tcp(127.0.0.1:3306)/go-todolist" -path ./migrations down 1
```

## Run go
```bash
go run main.go
```

# Folders structure
```
├── controller
│   ├── categoryController.go
│   ├── googleOauthController.go
│   └── userController.go
├── entity
│   ├── categoryEntity.go
│   ├── redisEntity.go
│   └── userEntity.go
├── log
│   ├── info.log
│   └── error.log
├── middleware
│   ├── cors.go
│   ├── jwt.go
│   └── rateLimiter.go
├── migration
│   ├── 20221129000000_create_users_table.down.sql
│   ├── 20221129000000_create_users_table.up.sql
│   ├── 20221129000001_create_categories_table.down.sql
│   ├── 20221129000001_create_categories_table.up.sql
│   ├── 20221129000002_create_tasks_table.down.sql
│   ├── 20221129000002_create_tasks_table.up.sql
│   ├── 20221129000003_create_category_task_table.down.sql
│   └── 20221129000003_create_category_task_table.up.sql
├── model
│   ├── category.go
│   └── user.go
├── request
│   ├── categoryRequest.go
│   ├── publicRequest.go
│   └── userRequest.go
├── router
│   └── api.go
├── services
│   ├── jwtService.go
│   └── userService.go
├── utils
│   ├── gorm
│   │   └── gorm.go
│   ├── log
│   │   ├── logByDate.go
│   │   └── logBySize.go
│   ├── paginator
│   │   └── paginator.go
│   ├── redis
│   │   └── redis.go
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

- Log
> 日誌檔位置 <br>
> Location of logger files

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
