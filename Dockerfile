FROM golang:1.19-alpine

WORKDIR /var/www/app

RUN go install github.com/cosmtrek/air@latest & \
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

COPY go.mod go.sum ./

RUN go mod download

CMD ["air", "-c", ".air.toml"]
