FROM --platform=linux/amd64 golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

EXPOSE 8080

ENV USE_ENV_FILE=FALSE

CMD ["./app"]