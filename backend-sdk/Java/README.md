# RelayAPI Java SDK

RelayAPI Java SDK 是一个用于与 RelayAPI 服务器交互的客户端库。它提供了简单的接口用于生成 API URL、创建令牌和发送各种 API 请求。

## 编译和打包

### 环境要求
- JDK 11 或更高版本
- Maven 3.6 或更高版本

### 编译步骤

1. 克隆代码库：
```bash
git clone https://github.com/relayapi/RelayAPI.git
cd RelayAPI/backend-sdk/Java
```

2. 编译和打包：
```bash
mvn clean package
```
成功后会在 `target` 目录下生成 `relayapi-sdk-1.0.0.jar` 文件。

3. 安装到本地 Maven 仓库（可选）：
```bash
mvn install
```

### 运行测试程序

1. 准备配置文件
在项目根目录创建 `config.rai` 文件，内容如下：
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

2. 配置环境变量
复制 `.env.example` 文件为 `.env` 并设置您的 API key：
```bash
cp .env.example .env
# 编辑 .env 文件，设置您的 API key
```

3. 运行示例程序：
```bash
# 编译并运行
mvn compile exec:java -Dexec.mainClass="com.github.relayapi.sdk.examples.ChatExample"

# 或者使用 Java 命令运行
java -cp target/relayapi-sdk-1.0.0.jar com.github.relayapi.sdk.examples.ChatExample
```

注意：运行示例程序前，请确保：
- 已经启动了 RelayAPI 服务器
- 配置文件 `config.rai` 中的服务器地址和端口正确
- 已经在 `.env` 文件中设置了正确的 API key

## 安装

将以下依赖添加到您的 Maven 项目的 `pom.xml` 文件中：

```xml
<dependency>
    <groupId>com.github.relayapi</groupId>
    <artifactId>relayapi-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

## 配置

SDK 需要一个配置对象进行初始化。您可以从文件（`.rai`）加载配置，或直接传递配置对象。配置格式示例：

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

## 使用示例

### 基本用法

```java
import com.relayapi.sdk.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.nio.file.Files;
import java.nio.file.Paths;

// 从文件加载配置
String configContent = Files.readString(Paths.get("config.rai"));
ObjectMapper mapper = new ObjectMapper();
Config config = mapper.readValue(configContent, Config.class);

// 创建客户端实例
RelayAPIClient client = new RelayAPIClient(config);

// 创建令牌
TokenOptions options = new TokenOptions("your-api-key")
    .setMaxCalls(100)
    .setExpireSeconds(3600)
    .setProvider("openai");
String token = client.createToken(options);

// 生成 API URL
String baseUrl = client.generateUrl(token, "");
System.out.println("Base URL: " + baseUrl);
// 输出示例：http://localhost:8840/relayapi/?token=xxxxx&rai_hash=xxxxx
```

### 聊天请求

```java
import java.util.List;
import java.util.Map;

RelayAPIClient.ChatRequest chatRequest = new RelayAPIClient.ChatRequest(token)
    .setMessages(List.of(
        Map.of("role", "system", "content", "You are a helpful assistant."),
        Map.of("role", "user", "content", "What is the capital of France?")
    ))
    .setModel("gpt-3.5-turbo")
    .setTemperature(0.7)
    .setMaxTokens(1000);

Map<String, Object> response = client.chat(chatRequest);
```

### 图像生成

```java
RelayAPIClient.ImageRequest imageRequest = new RelayAPIClient.ImageRequest(token)
    .setPrompt("A beautiful sunset over Paris")
    .setModel("dall-e-3")
    .setSize("1024x1024")
    .setQuality("standard")
    .setN(1);

Map<String, Object> response = client.generateImage(imageRequest);
```

### 生成嵌入向量

```java
RelayAPIClient.EmbeddingRequest embeddingRequest = new RelayAPIClient.EmbeddingRequest(token)
    .setInput("The quick brown fox jumps over the lazy dog")
    .setModel("text-embedding-ada-002");

Map<String, Object> response = client.createEmbedding(embeddingRequest);
```

### 健康检查

```java
Map<String, Object> status = client.healthCheck();
```

## API 参考

### RelayAPIClient

#### 构造函数

```java
new RelayAPIClient(Config config)
```

- `config`: 配置对象

#### 方法

##### createToken(TokenOptions options)

创建新令牌。

- `options.apiKey`: API 密钥
- `options.maxCalls`: 最大调用次数（默认：100）
- `options.expireSeconds`: 过期秒数（默认：86400，24小时）
- `options.provider`: 提供商名称或 URL。当提供 URL 时，它将直接用作提供商端点。支持的提供商名称：'dashscope'、'openai' 等
- `options.extInfo`: 扩展信息（可选）

##### generateUrl(String token, String endpoint)

生成 API URL。

- `token`: 令牌字符串
- `endpoint`: API 端点路径

##### chat(ChatRequest request)

发送聊天请求。

- `request.messages`: 消息数组
- `request.model`: 模型名称（默认：'gpt-3.5-turbo'）
- `request.temperature`: 温度值（默认：0.7）
- `request.maxTokens`: 最大令牌数（默认：1000）
- `request.token`: 令牌字符串

##### generateImage(ImageRequest request)

生成图像。

- `request.prompt`: 图像描述
- `request.model`: 模型名称（默认：'dall-e-3'）
- `request.size`: 图像大小（默认：'1024x1024'）
- `request.quality`: 图像质量（默认：'standard'）
- `request.n`: 生成图像数量（默认：1）
- `request.token`: 令牌字符串

##### createEmbedding(EmbeddingRequest request)

生成嵌入向量。

- `request.input`: 输入文本
- `request.model`: 模型名称（默认：'text-embedding-ada-002'）
- `request.token`: 令牌字符串

##### healthCheck()

检查服务器健康状态。

## 错误处理

SDK 中的所有方法在发生错误时都会抛出异常。建议使用 try-catch 块来处理潜在的错误：

```java
try {
    Map<String, Object> response = client.chat(chatRequest);
} catch (IOException e) {
    System.err.println("Error: " + e.getMessage());
}
```

## 示例程序

查看 `examples` 目录中的示例程序以获取更多使用示例：

- `ChatExample.java`: 演示如何使用 SDK 进行聊天、图像生成和嵌入向量生成

## 许可证

MIT 