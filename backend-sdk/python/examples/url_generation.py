#!/usr/bin/env python3

import os
from relayapi import RelayAPIClient

def main():
    # 创建客户端实例
    client = RelayAPIClient("../../../server/default.rai")
    
    # 方法 1: 使用已有令牌生成 URL
    print("\n方法 1: 使用已有令牌生成 URL")
    api_key = os.getenv("API_KEY", "sk-573af3eca24f492a83d5e64894ed91f5")
    token = client.create_token(
        api_key=api_key,
        max_calls=100,
        expire_days=1,
        provider="dashscope"
    )
    
    # 为不同 API 生成 URL
    for api_type in ['chat_completions', 'images_generations', 'embeddings']:
        url = client.generate_api_url_with_token(token, api_type)
        print(f"{api_type} URL: {url}")
    
    # 方法 2: 一步生成带有新令牌的 URL
    print("\n方法 2: 一步生成带有新令牌的 URL")
    for api_type in ['chat_completions', 'images_generations', 'embeddings']:
        url = client.generate_api_url(
            api_key=api_key,
            api_type=api_type,
            max_calls=100,
            provider="dashscope"
        )
        print(f"{api_type} URL: {url}")
    
    # 测试错误处理
    print("\n测试错误处理")
    try:
        url = client.generate_api_url_with_token(token, 'invalid_api')
    except ValueError as e:
        print(f"预期的错误: {e}")

if __name__ == "__main__":
    main() 