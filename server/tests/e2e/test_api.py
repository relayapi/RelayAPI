#!/usr/bin/env python3

import json
import base64
import time
from datetime import datetime, timedelta,timezone
from typing import Dict, Any

import requests
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from Crypto.Util.Padding import pad, unpad

class TokenGenerator:
    def __init__(self, config_path: str = "../../config.json"):
        # 读取配置文件
        with open(config_path, 'r') as f:
            self.config = json.load(f)
        
        # 获取加密配置
        self.crypto_config = self.config['crypto']
        if self.crypto_config['method'] != 'aes':
            raise ValueError("Only AES encryption is supported")
        
        # 解码 AES 密钥和 IV 种子
        self.key = bytes.fromhex(self.crypto_config['aes_key'])
        self.iv_seed = self.crypto_config['aes_iv_seed'].encode()

    def generate_iv(self) -> bytes:
        """生成 IV"""
        iv = get_random_bytes(AES.block_size)
        # 使用 IV 种子进行混合
        return bytes(a ^ b for a, b in zip(iv, self.iv_seed))

    def encrypt_token(self, token_data: Dict[str, Any]) -> str:
        """加密令牌数据"""
        # 序列化令牌数据
        json_data = json.dumps(token_data, separators=(',', ':')).encode()
        
        # 生成 IV
        iv = self.generate_iv()
        
        # 创建 AES 加密器
        cipher = AES.new(self.key, AES.MODE_CBC, iv)
        
        # 加密数据
        padded_data = pad(json_data, AES.block_size)
        encrypted_data = cipher.encrypt(padded_data)
        
        # 组合 IV 和加密数据
        result = iv + encrypted_data
        
        # Base64 URL 安全编码
        return base64.urlsafe_b64encode(result).decode()

def create_test_token() -> Dict[str, Any]:
    """创建测试令牌数据"""
    now = datetime.now(timezone.utc)
    return {
        "id": "test-token-1",
        "api_key": "sk-test123456789",
        "max_calls": 100,
        "used_calls": 0,
        "expire_time": (now + timedelta(days=1)).isoformat(),
        "created_at": now.isoformat(),
        "updated_at": now.isoformat(),
        "ext_info": "test token"
    }

def test_openai_api():
    """测试 OpenAI API 代理"""
    # 创建令牌生成器
    generator = TokenGenerator()
    
    # 创建测试令牌
    token_data = create_test_token()
    encrypted_token = generator.encrypt_token(token_data)
    
    # API 配置
    base_url = f"http://localhost:{generator.config['server']['port']}"
    headers = {
        "Content-Type": "application/json"
    }
    
    # 测试健康检查
    print("\n测试健康检查...")
    resp = requests.get(f"{base_url}/health")
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.json()}")
    assert resp.status_code == 200
    
    # 测试 OpenAI API
    print("\n测试 OpenAI API...")
    data = {
        "model": "gpt-3.5-turbo",
        "messages": [{"role": "user", "content": "Hello!"}]
    }
    
    # 使用加密令牌调用 API
    encrypted_token = encrypted_token.strip()
    url = f"{base_url}/api/openai/v1/chat/completions?token={encrypted_token}"
    print(f"请求 URL: {url}")
    print(f"加密令牌: {encrypted_token}")
    print(f"令牌长度: {len(encrypted_token)}")
    print(f"令牌字节: {[ord(c) for c in encrypted_token[:10]]}...")
    
    try:
        resp = requests.post(url, headers=headers, json=data)
        print(f"状态码: {resp.status_code}")
        print(f"响应: {resp.text}")
    except Exception as e:
        print(f"请求失败: {e}")
    
    # 测试无效令牌
    print("\n测试无效令牌...")
    invalid_url = f"{base_url}/api/openai/v1/chat/completions?token=invalid_token"
    try:
        resp = requests.post(invalid_url, headers=headers, json=data)
        print(f"状态码: {resp.status_code}")
        print(f"响应: {resp.json()}")
    except Exception as e:
        print(f"请求失败: {e}")

if __name__ == "__main__":
    test_openai_api() 