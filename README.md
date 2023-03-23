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

# Contents
 - [Software requirements](#software-requirements)
 - [Project plugins](#project-plugins)
 - [How to build project](#how-to-build-project)
 - [Folder structure](#folder-structure)
 - [Folder definition](#folder-definition)
 - [How to get telegram notifications](#how-to-get-telegram-notifications)

# Software requirement
 - **Database**
    - [MySQL](https://dev.mysql.com/downloads/mysql/): v8.0
    - [Redis](https://redis.io/download/): v6.2
 - **Programming language**
    - [Go](https://go.dev/dl/): v1.19
 - **Deveops**
    - [Docker GUI](https://www.docker.com/products/docker-desktop/)
 - **Other**
    - [Postman](https://www.postman.com/downloads/)
 
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
- [aws sdk for go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [uuid](https://github.com/gofrs/uuid)

# How to build project
## 1.Clone GitHub project to local
```bash
git clone https://github.com/Frankie0702111/go-todolist.git
```

## 2.Generate .env and set up environment
```bash
cd go-todolist
cp .env.example .env

# Set up basic information, such as database, Google Oauth, JWT
vim .env
```

## 3.Build docker image and start
```bash
cd go-todolist

# Create docker image
docker compose build --no-cache

# Run docker
docker compose up -d

# Stop docker
docker compose stop
```

## 4.Generate db migrations
```bash
cd go-todolist

# Up migration
make migrate-up

# Down migration
make migrate-down

# Specify batch up or down (If you want to go down to a specific file, it is recommended to open a new folder)
make migrate-up number=1
make migrate-down number=1
```

# Folder structure
```
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── controller
│   ├── categoryController.go
│   ├── googleOauthController.go
│   ├── taskController.go
│   └── userController.go
├── docker-compose.yaml
├── entity
│   ├── categoryEntity.go
│   ├── redisEntity.go
│   ├── s3Entity.go
│   ├── taskEntity.go
│   └── userEntity.go
├── go.mod
├── go.sum
├── log
├── main.go
├── middleware
│   ├── cors.go
│   ├── jwt.go
│   └── rateLimiter.go
├── migrations
│   ├── 20221129000000_create_users_table.down.sql
│   ├── 20221129000000_create_users_table.up.sql
│   ├── 20221129000001_create_categories_table.down.sql
│   ├── 20221129000001_create_categories_table.up.sql
│   ├── 20221129000002_create_tasks_table.down.sql
│   └── 20221129000002_create_tasks_table.up.sql
├── model
│   ├── category.go
│   ├── task.go
│   └── user.go
├── request
│   ├── categoryRequest.go
│   ├── publicRequest.go
│   ├── taskRequest.go
│   └── userRequest.go
├── router
│   └── api.go
├── services
│   ├── categoryService.go
│   ├── jwtService.go
│   ├── taskService.go
│   └── userService.go
└── utils
    ├── aws
    │   └── s3.go
    ├── civilDatetime
    │   └── civilDatetime.go
    ├── gorm
    │   └── gorm.go
    ├── log
    │   ├── logByDate.go
    │   └── logBySize.go
    ├── paginator
    │   └── paginator.go
    ├── redis
    │   └── redis.go
    └── responses
        └── response.go
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

# How to get telegram notifications
1. Sign in to your Telegram account.
2. Search for this account "@GolangToDoListBot" and press "Start".
3. Input this command "/setaccess {Email}" to set your configuration.
4. Input this command "/tasklist" to check for incomplete tasks.
5. To be continued ...