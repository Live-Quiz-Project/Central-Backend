FROM --platform=linux/amd64 golang:latest

WORKDIR /app

# Set GOPROXY and GOSUMDB environment variables
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

EXPOSE 8080

ENV USE_ENV_FILE=FALSE

CMD ["./app"]