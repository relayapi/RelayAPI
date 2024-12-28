import json
from typing import Dict, Any, Optional, List, Union

import requests

from .token import TokenGenerator

class RelayAPIClient:
    """RelayAPI 客户端，用于调用 API 服务"""

    def __init__(self, config_path: str = "default.rai"):
        """
        初始化 RelayAPI 客户端
        
        Args:
            config_path: .rai 配置文件路径
        """
        self.token_generator = TokenGenerator(config_path)
        self.headers = {
            "Content-Type": "application/json"
        }

    def create_token(
        self,
        api_key: str,
        max_calls: int = 100,
        expire_days: int = 1,
        provider: str = "dashscope",
        ext_info: str = ""
    ) -> str:
        """
        创建并加密访问令牌
        
        Args:
            api_key: API 密钥
            max_calls: 最大调用次数
            expire_days: 过期天数
            provider: API 提供者 (openai/dashscope)
            ext_info: 扩展信息
            
        Returns:
            str: 加密后的令牌字符串
        """
        token_data = self.token_generator.create_token(
            api_key=api_key,
            max_calls=max_calls,
            expire_days=expire_days,
            provider=provider,
            ext_info=ext_info
        )
        return self.token_generator.encrypt_token(token_data)

    def chat_completions(
        self,
        token: str,
        messages: List[Dict[str, str]],
        model: str = "qwen-vl-max",
        **kwargs: Any
    ) -> Dict[str, Any]:
        """
        调用聊天补全 API
        
        Args:
            token: 加密的访问令牌
            messages: 对话消息列表
            model: 模型名称
            **kwargs: 其他参数
            
        Returns:
            Dict[str, Any]: API 响应
        """
        url = self.token_generator.get_server_url("/chat/completions")
        data = {
            "model": model,
            "messages": messages,
            **kwargs
        }
        
        response = requests.post(
            f"{url}?token={token.strip()}",
            headers=self.headers,
            json=data
        )
        response.raise_for_status()
        return response.json()

    def images_generations(
        self,
        token: str,
        prompt: str,
        n: int = 1,
        size: str = "1024x1024",
        **kwargs: Any
    ) -> Dict[str, Any]:
        """
        调用图像生成 API
        
        Args:
            token: 加密的访问令牌
            prompt: 图像描述
            n: 生成图像数量
            size: 图像尺寸
            **kwargs: 其他参数
            
        Returns:
            Dict[str, Any]: API 响应
        """
        url = self.token_generator.get_server_url("/images/generations")
        data = {
            "prompt": prompt,
            "n": n,
            "size": size,
            **kwargs
        }
        
        response = requests.post(
            f"{url}?token={token.strip()}",
            headers=self.headers,
            json=data
        )
        response.raise_for_status()
        return response.json()

    def embeddings(
        self,
        token: str,
        input: Union[str, List[str]],
        model: str = "text-embedding-ada-002",
        **kwargs: Any
    ) -> Dict[str, Any]:
        """
        调用文本嵌入 API
        
        Args:
            token: 加密的访问令牌
            input: 输入文本或文本列表
            model: 模型名称
            **kwargs: 其他参数
            
        Returns:
            Dict[str, Any]: API 响应
        """
        url = self.token_generator.get_server_url("/embeddings")
        data = {
            "model": model,
            "input": input,
            **kwargs
        }
        
        response = requests.post(
            f"{url}?token={token.strip()}",
            headers=self.headers,
            json=data
        )
        response.raise_for_status()
        return response.json()

    def health_check(self) -> Dict[str, Any]:
        """
        检查服务器健康状态
        
        Returns:
            Dict[str, Any]: 健康状态信息
        """
        url = self.token_generator.get_server_url("")
        url = f"{url}/health"
        print(url)
        response = requests.get(url)
        response.raise_for_status()
        return response.json()