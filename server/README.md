# RelayAPI Server

安全的 API 代理服务器，用于保护敏感的 API 密钥。

## 编译

```bash
# 设置 Go 环境
export GOROOT=/usr/local/go
go mod download && go mod tidy

# 编译
cd server
go build -o bin/relayapi-server cmd/server/main.go
```

## 配置

创建 `config.json` 文件：

```json
{
  "server": {
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "your_password",
    "dbname": "relayapi"
  },
  "crypto": {
    "private_key_path": "keys/private.pem",
    "public_key_path": "keys/public.pem"
  }
}
```

## 运行

```bash
# 直接运行
./bin/relayapi-server

# 或使用 go run
go run cmd/server/main.go
```

## API 使用

1. 健康检查：
```bash
curl http://localhost:8080/health
```

2. OpenAI API 代理：
```bash
curl -X POST http://localhost:8080/api/openai/v1/chat/completions \
  -H "X-API-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello"}]}'
``` 