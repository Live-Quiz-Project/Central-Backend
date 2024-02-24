FROM golang:alpine3.18

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -o app ./cmd/main.go

EXPOSE 8080

ENV USE_ENV_FILE=FALSE

CMD ["./app"]