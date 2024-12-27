
```markdown:README.md
# RelayAPI

RelayAPI 是一个安全的 API 代理层，专门用于解决前端调用需要 API Key 的服务（如 OpenAI）时的安全问题。通过非对称加密和访问控制，实现安全可控的 API 调用。

## 特性

- 🔒 安全性：使用 ECC 非对称加密保护敏感信息
- 🚀 高性能：基于 Go 实现的高并发服务器
- 🎯 精确控制：支持调用次数和时间限制
- 🔌 多语言支持：提供多种语言的加密 SDK
- 🛡️ API Key 保护：敏感信息不在前端暴露
- 📊 使用统计：支持调用量统计和监控

## 架构

```ascii
步骤1: Backend 生成加密令牌
┌─────────┐    ┌──────────────┐
│ Backend ├────┤ RelayAPI SDK │
│         │    │  加密模块     │
└─────────┘    └──────────────┘

步骤2: Frontend 使用令牌调用服务
┌─────────┐    ┌─────────────┐    ┌──────────┐
│ Frontend├────┤ RelayAPI    ├────┤ OpenAI   │
│         │    │ Server      │    │ API      │
└─────────┘    └─────────────┘    └──────────┘
```

## 工作流程

1. Backend 使用 RelayAPI SDK 生成加密令牌
   - 设置调用限制（次数/时间）
   - 使用公钥加密 API Key 和参数
   - 生成加密令牌返回给前端

2. Frontend 调用 API
   - 携带加密令牌请求 RelayAPI Server
   - RelayAPI Server 验证并解密令牌
   - 使用解密后的 API Key 调用目标服务
   - 返回响应结果给前端

## 实现步骤

### 1. RelayAPI Server 开发 (Go)

1. 基础框架搭建
   - 使用 Gin 框架搭建 HTTP 服务
   - 配置数据库（PostgresSQL）用于存储调用记录
   - 实现健康检查接口

2. 加密系统实现
   - 实现 ECC 加密/解密功能
   - 生成并管理公私钥对
   - 实现令牌验证和解密

3. 访问控制
   - 实现调用次数限制
   - 实现时间有效期检查
   - 实现并发控制

4. API 代理功能
   - 实现 OpenAI API 代理
   - 错误处理和重试机制
   - 响应数据转发

### 2. RelayAPI 加密 SDK 开发

1. 核心功能
   - 参数加密
   - 令牌生成
   - 调用限制设置

2. 多语言支持
   - Node.js SDK
   - Python SDK
   - Java SDK
   - Go SDK

### 3. 部署和运维

1. 服务器部署
   - Docker 容器化
   - 负载均衡配置
   - 监控系统搭建

2. 文档编写
   - API 文档
   - SDK 使用文档
   - 部署文档

## 快速开始

### 安装

```bash
# 后端安装加密 SDK
pnpm install @relayapi/sdk

# RelayAPI Server 独立部署
docker pull relayapi/server
```

### 后端使用示例

```typescript
import { RelayAPISDK } from '@relayapi/sdk'

const sdk = new RelayAPISDK({
  publicKey: 'YOUR_PUBLIC_KEY'
})

// 生成加密令牌
const token = await sdk.generateToken({
  apiKey: 'sk-...',
  maxCalls: 100,
  expireTime: '2024-12-31',
  allowedModels: ['gpt-4']
})

// 返回令牌给前端
return { token }
```

### 前端使用示例

```typescript
// 使用后端生成的令牌调用 API
const response = await fetch('https://relay.api.example.com/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    model: 'gpt-4',
    messages: [...]
  })
})
```

## 安全说明

1. 私钥仅在 RelayAPI Server 端保存
2. 后端服务器通过公钥生成加密令牌
3. 令牌包含使用限制和过期时间
4. 所有敏感信息经过加密，前端仅能访问授权范围内的 API

## 贡献指南

欢迎提交 Pull Request 和 Issue。

## 许可证

MIT
```

这个版本更准确地描述了系统的工作流程：
1. 后端使用 SDK 和公钥生成加密令牌
2. 前端使用令牌直接与 RelayAPI Server 通信
3. RelayAPI Server 使用私钥解密并验证令牌
4. 最后代理调用目标 API

这种方式确保了：
- API Key 等敏感信息只在后端和 RelayAPI Server 之间传递
- 前端只能在授权的范围内使用 API
- 通过令牌实现了精确的访问控制

### 具体实现步骤和技术选型

1. **后端技术栈**
   - 主框架：Go + Gin
   - 数据库：Redis（调用计数）+ MongoDB（日志记录）
   - 加密：golang.org/x/crypto（ECC 实现）
   - API 文档：Swagger

2. **开发流程**

第一阶段：核心功能实现
- 搭建基础项目结构
- 实现 ECC 加密解密功能
- 实现基础的 API 代理功能
- 编写单元测试

第二阶段：SDK 开发
- 实现 Node.js SDK
- 实现参数加密功能
- 实现请求构建功能
- 编写 SDK 文档

第三阶段：完善功能
- 实现调用限制
- 实现监控统计
- 实现多语言 SDK
- 完善错误处理

第四阶段：部署和优化
- 容器化部署
- 性能优化
- 监控系统
- 文档完善

这个项目的关键在于安全性和性能的平衡，建议先实现核心功能的原型，然后再逐步完善其他功能。
