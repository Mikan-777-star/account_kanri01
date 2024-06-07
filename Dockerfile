# ベースイメージとしてGoを指定
FROM golang:1.20-alpine AS builder

# アプリケーションディレクトリを作成
WORKDIR /app

# 必要なファイルをコピー

COPY . ./
# Check if go.mod and go.sum exist, then remove them
RUN [ -f go.mod ] && rm go.mod || true
RUN [ -f go.sum ] && rm go.sum || true

RUN go mod init example.com/mymodule && \
    go mod tidy && \
    go mod vendor

RUN go mod download


# アプリケーションをビルド
RUN go build -o account_management




FROM alpine:latest

WORKDIR /app

RUN apk add --update-cache --no-cache bash curl && \
    curl -o wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh 

COPY --from=builder  /app/account_management .

RUN chmod 111 ./wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "mysql:3306", "--","./account_management"]