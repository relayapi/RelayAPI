#!/usr/bin/env python3

import os
from relayapi import RelayAPIClient
from openai import OpenAI

# 使用配置字典直接创建客户端
config = {
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

# 创建客户端实例
client = RelayAPIClient(config)

# 创建令牌
token = client.create_token(
    api_key="your-api-key",
    max_calls=100,
    expire_days=1,
    provider="openai"
)

# 生成基础 URL（不指定 API 类型）
base_url = client.generate_api_url_with_token(token)
print("Base URL:", base_url)
# 输出示例: http://localhost:8840/relayapi/?token=xxxxx&rai_hash=xxxxx

# 生成聊天 API URL
chat_url = client.generate_api_url_with_token(token, 'chat_completions')
print("Chat API URL:", chat_url)

# 生成图像 API URL
image_url = client.generate_api_url_with_token(token, 'images_generations')
print("Image API URL:", image_url)

# 生成嵌入 API URL
embedding_url = client.generate_api_url_with_token(token, 'embeddings')
print("Embedding API URL:", embedding_url)

# 使用 OpenAI 客户端示例
openai_client = OpenAI(
    base_url=base_url,
    api_key="not-needed"  # API key 已包含在 token 中
)

if __name__ == "__main__":
    main() 