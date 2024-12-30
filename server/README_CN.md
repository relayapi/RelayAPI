# RelayAPI 服务器

[English Documentation](README.md)

RelayAPI 服务器是处理 API 代理、令牌验证和请求转发的核心组件。

## 快速开始

```bash
# 克隆仓库
git clone https://github.com/relayapi/RelayAPI.git

# 进入服务器目录
cd server

# 使用默认配置运行服务器
go run cmd/server/main.go -rai ./rai -d
```

命令行选项：
- `-rai`：客户端配置文件目录路径（默认：当前目录）
- `-config`：服务器配置文件路径（默认：config.json）
- `-d`：启用调试模式，日志将写入 debug.log

## 配置

### 服务器配置（`config.json`）

服务器配置文件控制服务器行为，包括网络设置、速率限制和日志选项。

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

### 客户端配置（`.rai` 文件）

客户端配置文件包含加密设置，由服务器和 SDK 共同使用。服务器会监控 `-rai` 目录中的这些文件。

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

您可以根据需要自定义 `crypto` 部分，只需确保服务器和 SDK 使用相同的配置即可。

## 开发

### 项目结构

```
server/
├── cmd/
│   └── server/
│       └── main.go         # 服务器入口点
├── internal/
│   ├── config/            # 配置管理
│   ├── crypto/            # 加密实现
│   ├── handlers/          # 请求处理器
│   ├── middleware/        # 中间件组件
│   ├── models/            # 数据模型
│   └── services/          # 业务逻辑
├── rai/                   # 客户端配置文件
└── config.json            # 服务器配置
```

### 核心组件

1. **配置管理**（`internal/config/`）
   - 加载和验证服务器配置
   - 监控客户端配置文件
   - 处理配置更新

2. **加密**（`internal/crypto/`）
   - 实现 AES 加密/解密
   - 管理加密密钥
   - 处理令牌生成和验证

3. **请求处理**（`internal/handlers/`）
   - 处理传入的 API 请求
   - 验证令牌
   - 转发请求到 AI 提供商

4. **中间件**（`internal/middleware/`）
   - 认证
   - 速率限制
   - 日志记录
   - 请求/响应转换

### 添加新功能

1. **新的 AI 提供商**
   ```go
   // internal/handlers/provider.go
   func (h *Handler) handleProviderRequest(c *gin.Context) {
       // 实现提供商特定的处理逻辑
   }
   ```

2. **新的中间件**
   ```go
   // internal/middleware/custom.go
   func CustomMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // 实现中间件逻辑
       }
   }
   ```

3. **新的配置选项**
   ```go
   // internal/config/config.go
   type Config struct {
       // 添加新的配置字段
   }
   ```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 运行特定包的测试
go test ./internal/crypto
```

## 部署

1. **构建**
   ```bash
   go build -o relayapi cmd/server/main.go
   ```

2. **运行**
   ```bash
   ./relayapi -rai /path/to/rai/dir -config /path/to/config.json
   ```

3. **监控**
   - 使用 `-d` 运行时查看 `debug.log` 获取详细日志
   - 通过 `/health` 端点监控服务器状态
   - 查看控制台输出获取实时统计信息

## 贡献

1. Fork 本仓库
2. 创建您的特性分支
3. 提交您的更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](../LICENSE) 文件。 