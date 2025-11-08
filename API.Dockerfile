FROM golang:1.25-alpine AS build

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE ${PORT}

# Set Air for development mode
CMD ["air", "-c", ".air-api.toml"]

