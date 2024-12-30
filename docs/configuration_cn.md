# RelayAPI 配置指南

RelayAPI 使用两种配置文件来管理服务器和客户端设置。本指南将详细说明它们的用途、格式和使用方法。

## 服务器配置 (`config.json`)

服务器配置文件是运行 RelayAPI 服务器所**必需**的。它包含了服务器运行所需的所有设置，包括网络设置、速率限制和日志选项。

### 基本用法

```bash
# 使用默认配置文件启动服务器
relayapi

# 使用自定义配置文件启动服务器
relayapi --config=/path/to/config.json
```

### 配置结构

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8840,
    "read_timeout": 30,
    "write_timeout": 30,
    "max_header_bytes": 1048576,
    "debug": false
  },
  "log": {
    "console": true,
    "database": {
      "enabled": true,
      "type": "postgres",
      "connection_string": "user=postgres password=postgres dbname=relayapi host=localhost port=5432 sslmode=disable"
    },
    "web": {
      "enabled": false,
      "callback_url": "http://example.com/log"
    },
    "parquet": {
      "enabled": false,
      "file_path": "/path/to/logs/output.parquet"
    }
  },
  "rate_limit": {
    "requests_per_second": 20,
    "burst": 40,
    "ip_limit": {
      "requests_per_second": 10,
      "burst": 20
    }
  }
}
```

### 配置选项

#### 服务器设置
- `host`: 服务器监听地址
- `port`: 服务器监听端口
- `read_timeout`: 请求读取超时时间（秒）
- `write_timeout`: 响应写入超时时间（秒）
- `max_header_bytes`: 请求头部最大大小
- `debug`: 启用调试模式（默认：false）

#### 日志设置
- `console`: 启用控制台日志
- `database`: 数据库日志配置
  - `enabled`: 启用数据库日志
  - `type`: 数据库类型（postgres、mysql、sqlite）
  - `connection_string`: 数据库连接字符串
- `web`: Web 回调日志配置
  - `enabled`: 启用 Web 回调日志
  - `callback_url`: 日志回调 URL
- `parquet`: Parquet 文件日志配置
  - `enabled`: 启用 Parquet 文件日志
  - `file_path`: Parquet 日志文件保存路径

#### 速率限制
- `requests_per_second`: 全局请求速率限制
- `burst`: 全局突发限制
- `ip_limit`: 每个 IP 的速率限制
  - `requests_per_second`: 每个 IP 的请求速率限制
  - `burst`: 每个 IP 的突发限制

## 客户端配置 (`default.rai`)

客户端配置文件包含 SDK 运行所需的设置，包括加密设置和服务器连接信息。如果不存在，将自动生成默认配置。

### 基本用法

```typescript
// 从文件加载
const client = new RelayAPIClient('default.rai');

// 或直接传入配置对象
const client = new RelayAPIClient({
  version: "1.0.0",
  server: {
    host: "http://localhost",
    port: 8840,
    base_path: "/relayapi/"
  },
  crypto: {
    method: "aes",
    aes_key: "your-aes-key",
    aes_iv_seed: "your-iv-seed"
  }
});
```

### 配置结构

```json
{
  "version": "1.0.0",
  "server": {
    "host": "http://localhost",
    "port": 8840,
    "base_path": "/relayapi/"
  },
  "crypto": {
    "method": "aes",
    "aes_key": "your-aes-key",
    "aes_iv_seed": "your-iv-seed"
  }
}
```

### 配置选项

#### 版本
- `version`: 配置版本（当前：1.0.0）

#### 服务器连接
- `server.host`: RelayAPI 服务器主机地址
- `server.port`: RelayAPI 服务器端口
- `server.base_path`: API 端点的基础路径

#### 加密设置
- `crypto.method`: 加密方法（当前支持：aes）
- `crypto.aes_key`: AES 加密密钥
- `crypto.aes_iv_seed`: AES IV 种子

### 自动生成逻辑

在以下情况下，`default.rai` 文件将被自动生成：
1. 当前目录不存在该文件
2. SDK 初始化时未提供配置对象
3. 服务器启动时未指定客户端配置文件

自动生成的配置将使用安全的随机值作为加密密钥，并使用默认的服务器设置。

## 配置加载优先级

1. 命令行参数（--config、--rai）
2. 环境变量
3. 当前目录中的配置文件
4. 自动生成的默认配置

## 安全注意事项

1. 确保 `config.json` 和 `default.rai` 文件的安全
2. 不要将加密密钥提交到版本控制系统
3. 在生产环境中使用环境变量存储敏感信息
4. 定期轮换加密密钥
5. 在生产环境中使用强随机加密密钥

## 最佳实践

1. 为开发和生产环境使用不同的配置
2. 保持配置文件的备份
3. 记录任何自定义配置的更改
4. 在生产环境中监控速率限制
5. 定期检查和更新配置

更多具体的配置示例和使用场景，请参考：
- [服务器配置指南](../server/README.md)
- [JavaScript SDK 指南](../backend-sdk/JavaScript/README.md)
- [Python SDK 指南](../backend-sdk/python/README.md)
``` 