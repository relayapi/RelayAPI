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
    "port": 8840
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "your_password",
    "dbname": "relayapi"
  },
  "crypto": {
    "method": "aes",
    "key_size": 256,
    "private_key_path": "keys/private.pem",
    "public_key_path": "keys/public.pem",
    "aes_key": "",
    "aes_iv_seed": ""
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
curl http://localhost:8840/health
```

2. OpenAI API 代理：
```bash
# 使用 URL 参数传递令牌
curl -X POST "http://localhost:8840/api/openai/v1/chat/completions?token=your_token" \
  -H "Content-Type: application/json" \
  -d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello"}]}'

# 或者在路径中包含令牌
curl -X POST "http://localhost:8840/api/openai/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello"}]}'
```

## 客户端使用

对于使用 OpenAI 官方客户端的应用，只需要：

1. 将 base URL 修改为 RelayAPI 服务器地址
2. 将 API Key 作为 URL 参数传递

例如，使用 OpenAI Node.js 客户端：

```javascript
import OpenAI from 'openai';

const openai = new OpenAI({
  baseURL: 'http://localhost:8840/api/openai/v1',
  apiKey: 'your_token', // 这里的令牌会自动被添加到 URL 参数中
});

const response = await openai.chat.completions.create({
  model: 'gpt-3.5-turbo',
  messages: [{ role: 'user', content: 'Hello!' }],
});
```

或者使用 Python 客户端：

```python
from openai import OpenAI

client = OpenAI(
    base_url='http://localhost:8840/api/openai/v1',
    api_key='your_token',  # 这里的令牌会自动被添加到 URL 参数中
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello!"}]
) 