FROM golang:1.24.1 AS builder

WORKDIR /app
COPY . .
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -tags=release -mod=readonly -v -o server

# 実行ステージ
FROM alpine:3.21.3
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tokyo

COPY --from=builder /app/server /server
EXPOSE 8080

ENTRYPOINT ["/server"]
