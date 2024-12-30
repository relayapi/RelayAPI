# RelayAPI Python SDK

RelayAPI Python SDK 是一个用于与 RelayAPI 服务器进行交互的客户端库。它提供了简单的接口来生成 API URL、创建令牌，以及发送各种 API 请求。

## 安装

# 从源代码安装

cd relayapi/backend-sdk/python
pip install -e .


使用 pip 安装（即将支持）：

```bash
pip install relayapi-sdk
```

## 配置

SDK 需要一个配置对象来初始化。你可以从配置文件（`.rai`）加载配置，或直接传入配置对象。配置格式示例：

```python
config = {
    "version": "1.0.0",
    "server": {
        "host": "http://localhost",
        "port": 8080,
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

```python
from relayapi import RelayAPIClient
from openai import OpenAI

# 创建客户端实例（使用配置对象）
client = RelayAPIClient(config)

# 创建令牌
token = client.create_token(
    api_key="your-api-key",
    max_calls=100,
    expire_days=1,
    provider="openai"
)

# 生成 API URL
base_url = client.generate_api_url_with_token(token)
print("Base URL:", base_url)
# 输出示例: http://localhost:8080/relayapi/?token=xxxxx&rai_hash=xxxxx

# 在前端代码中将此 URL 用作 OpenAI API 的基础 URL
openai_client = OpenAI(
    base_url=base_url,
    api_key="not-needed"  # 实际的 API key 已包含在 token 中
)
```

### 聊天请求

```python
response = client.chat_completions(
    token=token,
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "What is the capital of France?"}
    ],
    model="gpt-3.5-turbo"
)
```

### 图像生成

```python
response = client.images_generations(
    token=token,
    prompt="A beautiful sunset over Paris",
    model="dall-e-3",
    size="1024x1024",
    quality="standard",
    n=1
)
```

### 嵌入向量生成

```python
response = client.embeddings(
    token=token,
    input="The quick brown fox jumps over the lazy dog",
    model="text-embedding-ada-002"
)
```

### 健康检查

```python
status = client.health_check()
```

### URL 生成

`generate_api_url_with_token` 方法用于生成包含令牌和哈希参数的完整 API URL：

```python
# 生成基础 URL（不指定 API 类型）
base_url = client.generate_api_url_with_token(token)

# 生成特定 API 的 URL
chat_url = client.generate_api_url_with_token(token, 'chat_completions')
image_url = client.generate_api_url_with_token(token, 'images_generations')
embedding_url = client.generate_api_url_with_token(token, 'embeddings')
```

参数：
- `token` (str)：加密的令牌字符串
- `api_type` (str, 可选)：API 类型，默认为空字符串

该方法会自动将令牌和配置哈希作为 URL 参数添加。

## API 参考

### RelayAPIClient

#### 构造函数

```python
RelayAPIClient(config: Union[str, Dict[str, Any]] = "default.rai")
```

- `config`: 配置文件路径（字符串）或配置对象（字典）

#### 方法

##### create_token

创建并加密访问令牌。

```python
create_token(
    api_key: str,
    max_calls: int = 100,
    expire_days: int = 1,
    provider: str = "dashscope",
    ext_info: str = ""
) -> str
```

##### generate_api_url_with_token

生成完整的 API URL。

```python
generate_api_url_with_token(
    token: str,
    api_type: str = ""
) -> str
```

##### chat_completions

发送聊天请求。

```python
chat_completions(
    token: str,
    messages: List[Dict[str, str]],
    model: str = "gpt-3.5-turbo",
    **kwargs: Any
) -> Dict[str, Any]
```

##### images_generations

生成图像。

```python
images_generations(
    token: str,
    prompt: str,
    n: int = 1,
    size: str = "1024x1024",
    **kwargs: Any
) -> Dict[str, Any]
```

##### embeddings

生成嵌入向量。

```python
embeddings(
    token: str,
    input: Union[str, List[str]],
    model: str = "text-embedding-ada-002",
    **kwargs: Any
) -> Dict[str, Any]
```

##### health_check

检查服务器健康状态。

```python
health_check() -> Dict[str, Any]
```

## 错误处理

SDK 中的所有方法都会在发生错误时抛出异常。建议使用 try-except 块来处理可能的错误：

```python
try:
    response = client.chat_completions(...)
except Exception as e:
    print(f"Error: {e}")
```

## 示例程序

查看 `examples` 目录中的示例程序，了解更多使用方法：

- `chat.py`: 展示了如何使用 SDK 进行聊天
- `url_generation.py`: 展示了如何生成和使用 API URL

## 许可证

MIT
```