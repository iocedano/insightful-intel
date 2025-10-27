FROM golang:1.25-alpine AS build

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download


COPY . .

CMD ["air", "-c", ".air-cli.toml"]

