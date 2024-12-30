# 支持的服务商列表

RelayAPI 目前支持 90+ 个 AI 服务商的 API 代理。以下是完整的支持列表：

## 🤖 通用 AI 模型服务

| 服务商 | API 终端 | 说明 |
|-------|---------|------|
| OpenAI | api.openai.com/v1 | GPT-4, GPT-3.5, DALL-E 3 等 |
| Anthropic | api.anthropic.com/v1 | Claude 系列模型 |
| Google AI | generativelanguage.googleapis.com/v1beta | PaLM, Gemini 等 |
| Mistral AI | api.mistral.ai/v1 | Mistral 系列模型 |
| Cohere | api.cohere.ai/v1 | Command 系列模型 |
| AI21 Labs | api.ai21.com/v1 | Jurassic 系列模型 |
| Hugging Face | api-inference.huggingface.co/models | 开源模型推理 |

## ☁️ 云服务商 AI

| 服务商 | API 终端 | 说明 |
|-------|---------|------|
| Azure OpenAI | api.cognitive.microsoft.com/v1 | 微软 Azure 平台 |
| AWS | comprehend.us-east-1.amazonaws.com | AWS AI 服务 |
| Google Cloud | aiplatform.googleapis.com | Google Cloud AI |
| 阿里云 | ai.aliyun.com/api/v1 | 通义千问等 |
| 百度智能云 | aip.baidubce.com | 文心一言等 |
| 腾讯云 | api.ai.qq.com | 混元大模型等 |
| 华为云 | api.hicloud.com/ai/v1 | 盘古大模型等 |

## 🎯 专业领域 AI

### 图像处理
- Stability AI (api.stability.ai/v1)
- Replicate (api.replicate.com/v1)
- RunwayML (api.runwayml.com/v1)
- DeepAI (api.deepai.org/api)
- Clarifai (api.clarifai.com/v2)

### 语音识别/处理
- AssemblyAI (api.assemblyai.com/v2)
- Speechmatics (asr.api.speechmatics.com/v2)
- Rev.ai (api.rev.ai/speechtotext/v1)
- Otter.ai (api.otter.ai/v1)

### 自然语言处理
- DeepL (api.deepl.com/v2)
- Wolfram Alpha (api.wolframalpha.com/v1)
- Wit.ai (api.wit.ai/v1)
- DialogFlow (dialogflow.googleapis.com/v2)

### 机器学习平台
- H2O.ai (api.h2o.ai/v1)
- DataRobot (api.datarobot.com/v2)
- BigML (bigml.io/andromeda)
- Algorithmia (api.algorithmia.com/v1)

## 🚀 企业 AI 服务

### 咨询公司 AI
- Accenture AI (api.accenture.com/v1)
- Deloitte AI (api.deloitte.com/v1)
- EY AI (api.ey.com/v1)
- KPMG AI (api.kpmg.com/v1)
- PwC AI (api.pwc.com/v1)
- BCG AI (api.bcg.com/v1)
- McKinsey AI (api.mckinsey.com/v1)

### 科技公司 AI
- IBM Watson (api.ai.ibm.com/v1)
- Salesforce Einstein (api.einstein.ai/v2)
- Oracle AI (api.oracle.com/v1)
- SAP Leonardo (api.sap.com/v1)

### 航空航天 AI
- NASA JPL (api.jpl.nasa.gov/v1)
- ESA (api.esa.int/v1)
- SpaceX (api.spacex.com/v1)
- Boeing (api.boeing.com/v1)
- Airbus (api.airbus.com/v1)

## 🔧 使用说明

1. 在生成令牌时，通过 `provider` 参数指定服务商：
```typescript
const token = await relay.generateToken({
  provider: 'openai',  // 这里指定服务商
  apiKey: 'your-api-key',
  // ... 其他参数
})
```

2. 前端使用时只需修改 baseURL：
```typescript
const client = new OpenAI({
  baseURL: 'https://your-relay-server.com',
  apiKey: token
})
```

## 📝 注意事项

1. 不同服务商的 API 格式可能不同，请参考各自的官方文档
2. 部分服务商可能需要额外的认证参数，请在生成令牌时包含
3. 建议在测试环境中先验证 API 调用是否正常
4. 留意各服务商的 API 限制和计费规则

## 🆕 添加新的服务商

如果你需要支持新的 AI 服务商，可以：

1. 在 GitHub 上提交 Issue 或 PR
2. 提供服务商的 API 文档
3. 说明使用场景和需求
4. 我们会评估并尽快支持

## 🔄 更新记录

- 2024-01: 新增支持 Google Gemini API
- 2023-12: 新增支持 Mistral AI API
- 2023-11: 新增支持 Anthropic Claude 2.1
- 2023-10: 新增支持 OpenAI GPT-4 Turbo 