# RelayAPI Python SDK

RelayAPI Python SDK 是一个用于访问 RelayAPI Server 的客户端库。它提供了简单的接口来创建访问令牌和调用各种 API 服务。

## 安装

```bash
# 从源代码安装
git clone https://github.com/yourusername/relayapi.git
cd relayapi/backend-sdk/python
pip install -e .

# 或者直接使用 pip 安装（即将支持）
pip install relayapi
```

## 配置

SDK 需要一个 `.rai` 配置文件来初始化。默认会在当前目录查找 `default.rai`，也可以指定配置文件路径。

配置文件示例：

```json
{
    "version": "1.0.0",
    "server": {
        "host": "http://localhost",
        "port": 8080,
        "base_path": "/relayapi/"
    },
    "crypto": {
        "method": "aes",
        "aes_key": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
        "aes_iv_seed": "fedcba9876543210"
    }
}
```

## 快速开始

```python
from relayapi import RelayAPIClient

# 方式1：使用配置文件初始化
client = RelayAPIClient("default.rai")

# 方式2：使用配置字典初始化
config = {
    "version": "1.0.0",
    "server": {
        "host": "http://localhost",
        "port": 8080,
        "base_path": "/relayapi/"
    },
    "crypto": {
        "method": "aes",
        "aes_key": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
        "aes_iv_seed": "fedcba9876543210"
    }
}
client = RelayAPIClient(config)

# 生成 API URL
url = client.generate_api_url(
    api_key="your-api-key",
    api_type="chat_completions", 
    max_calls=100,
    expire_days=1,
    provider="dashscope"  # 或 "openai")

# 创建访问令牌
token = client.create_token(
    api_key="your-api-key",
    max_calls=100,
    expire_days=1,
    provider="dashscope"  # 或 "openai"
)

# 发送聊天请求
messages = [
    {"role": "user", "content": "你好！"}
]

response = client.chat_completions(
    token=token,
    messages=messages,
    model="qwen-vl-max"
)

print(response["choices"][0]["message"]["content"])
```

## API 文档

### RelayAPIClient

主要的客户端类，用于创建令牌和调用 API。

#### 初始化

客户端支持两种初始化方式：

1. 使用配置文件初始化
```python
client = RelayAPIClient(config_path: str = "default.rai")
```

2. 使用配置字典初始化
```python
config = {
    "version": "1.0.0",
    "server": {
        "host": "http://localhost",  # 服务器地址
        "port": 8080,                # 服务器端口
        "base_path": "/relayapi/"    # API 基础路径
    },
    "crypto": {
        "method": "aes",             # 加密方式（目前仅支持 aes）
        "aes_key": "hex_encoded_key",# 32字节的16进制编码密钥
        "aes_iv_seed": "iv_seed"     # IV种子（用于生成初始化向量）
    }
}
client = RelayAPIClient(config)
```

#### 方法


0. 一步生成 API URL（包含令牌创建）
```python
url = client.generate_api_url(
    api_key: str,                # API 密钥
    api_type: str,               # API 类型：chat_completions/images_generations/embeddings
    max_calls: int = 100,        # 最大调用次数
    expire_days: int = 1,        # 过期天数
    provider: str = "dashscope", # API 提供者
    ext_info: str = ""          # 扩展信息
) -> str
```

1. 创建令牌
```python
token = client.create_token(
    api_key: str,                # API 密钥
    max_calls: int = 100,        # 最大调用次数
    expire_days: int = 1,        # 过期天数
    provider: str = "dashscope", # API 提供者
    ext_info: str = ""          # 扩展信息
) -> str
```

2. 聊天补全
```python
response = client.chat_completions(
    token: str,                 # 访问令牌
    messages: List[Dict],       # 对话消息列表
    model: str = "qwen-vl-max", # 模型名称
    **kwargs                    # 其他参数
) -> Dict
```

3. 图像生成
```python
response = client.images_generations(
    token: str,                # 访问令牌
    prompt: str,               # 图像描述
    n: int = 1,               # 生成数量
    size: str = "1024x1024",  # 图��尺寸
    **kwargs                  # 其他参数
) -> Dict
```

4. 文本嵌入
```python
response = client.embeddings(
    token: str,                          # 访问令牌
    input: Union[str, List[str]],        # 输入文本
    model: str = "text-embedding-ada-002", # 模型名称
    **kwargs                            # 其他参数
) -> Dict
```

5. 健康检查
```python
status = client.health_check() -> Dict
```

6. 生成 API URL（带令牌）
```python
url = client.generate_api_url_with_token(
    token: str,                # 访问令牌
    api_type: str             # API 类型：chat_completions/images_generations/embeddings
) -> str
```


### TokenGenerator

令牌生成器类，用于创建和加密访问令牌。通常不需要直接使用此类，应该使用 `RelayAPIClient` 的方法。

## 示例

查看 `examples` 目录获取更多示例：

- `chat.py`: 聊天对话示例
- 更多示例正在添加中...

## 错误处理

SDK 使用 Python 的异常机制处理错误：

- `requests.exceptions.HTTPError`: API 调用失败
- `ValueError`: 参数验证失败
- `FileNotFoundError`: 配置文件不存在
- `json.JSONDecodeError`: 配置文件格式错误

## 依赖

- Python >= 3.7
- requests >= 2.25.0
- pycryptodome >= 3.10.0


## 许可证

MIT License
```