FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/api
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY migrations ./migrations

ENV GIN_MODE=release
ENV HTTP_PORT=8080

EXPOSE 8080

CMD ["/app/server"]
