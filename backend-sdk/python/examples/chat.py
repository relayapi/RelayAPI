#!/usr/bin/env python3

import os
from relayapi import RelayAPIClient

def main():
    # 创建客户端实例
    client = RelayAPIClient("../../../default.rai")
    
    # 创建访问令牌
    api_key = os.getenv("API_KEY", "sk-573af3eca24f492a83d5e64894ed91f5")
    token = client.create_token(
        api_key=api_key,
        max_calls=100,
        expire_days=1,
        provider="dashscope"  # 或 "openai"
    )
    
    # 检查服务器状态
    print("\n检查服务器状态...")
    try:
        status = client.health_check()
        print(f"服务器状态: {status}")
    except Exception as e:
        print(f"健康检查失败: {e}")
        return
    
    # 发送聊天请求
    print("\n发送聊天请求...")
    messages = [
        {"role": "user", "content": "你好！请介绍一下自己。"}
    ]
    
    try:
        response = client.chat_completions(
            token=token,
            messages=messages,
            model="qwen-vl-max"
        )
        print("\n助手回复:")
        print(response["choices"][0]["message"]["content"])
    except Exception as e:
        print(f"请求失败: {e}")

if __name__ == "__main__":
    main()