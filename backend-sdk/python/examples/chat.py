#!/usr/bin/env python3

import os
from relayapi import RelayAPIClient

# 创建客户端实例
client = RelayAPIClient("../../server/default.rai")

# 创建令牌
token = client.create_token(
    api_key=os.getenv("API_KEY", "your-api-key"),
    max_calls=100,
    expire_days=1,
    provider="openai"
)

# 生成 API URL
base_url = client.generate_api_url_with_token(token)
print("Base URL:", base_url)

# 使用 OpenAI 客户端
from openai import OpenAI

openai_client = OpenAI(
    base_url=base_url,
    api_key="not-needed"  # API key 已包含在 token 中
)

# 发送聊天请求
response = openai_client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "What is the capital of France?"}
    ]
)

print("\nChat Response:")
print(response.choices[0].message.content)

# 使用 RelayAPI 客户端直接发送请求
response = client.chat_completions(
    token=token,
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "What is the capital of France?"}
    ],
    model="gpt-3.5-turbo"
)

print("\nDirect API Response:")
print(response['choices'][0]['message']['content'])