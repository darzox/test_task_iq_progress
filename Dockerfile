# Build stage
FROM golang:1.22.5 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server -v ./cmd/server

# Deployment stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /server ./server
COPY cmd/server/migrations ./migrations

COPY docs ./docs
EXPOSE 8080

ENTRYPOINT ["/app/server"]
