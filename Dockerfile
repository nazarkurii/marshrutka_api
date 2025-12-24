#building the binary
FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GIN_MODE=release go build -o ./cmd/bin/app ./cmd/api/main.go

#copying binary into final image
FROM alpine:3.19

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/cmd/bin/app .
COPY ./static /app

CMD ["./app"]