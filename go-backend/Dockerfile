# Build stage
FROM golang:1.19-alpine3.15 as build-env
RUN apk add --no-cache git gcc
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openapi ./cmd/main.go

# Final stage
FROM alpine:3.15
COPY --from=build-env /app/openapi .
EXPOSE 8080/tcp
USER 1001
ENTRYPOINT ["./openapi"]