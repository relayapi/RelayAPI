import json
from typing import Dict, Any, Optional, List, Union

import requests

from .token import TokenGenerator

class RelayAPIClient:
    """RelayAPI 客户端，用于调用 API 服务"""

    def __init__(self, config: Union[str, Dict[str, Any]] = "default.rai"):
        """
        初始化 RelayAPI 客户端
        
        Args:
            config: .rai 配置文件路径或配置字典对象
        """
        self.token_generator = TokenGenerator(config)
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

    def generate_api_url_with_token(self, token: str, api_type: str) -> str:
        """
        根据令牌和 API 类型生成完整的 API URL
        
        Args:
            token: 加密的访问令牌
            api_type: API 类型 (chat_completions/images_generations/embeddings)
            
        Returns:
            str: 完整的 API URL，包含令牌参数
        """
        api_paths = {
            'chat_completions': '/chat/completions',
            'images_generations': '/images/generations',
            'embeddings': '/embeddings'
        }
        
        if api_type not in api_paths:
            raise ValueError(f"不支持的 API 类型: {api_type}")
            
        base_url = self.token_generator.get_server_url(api_paths[api_type])
        return f"{base_url}?token={token.strip()}"

    def generate_api_url(
        self,
        api_key: str,
        api_type: str,
        max_calls: int = 100,
        expire_days: int = 1,
        provider: str = "dashscope",
        ext_info: str = ""
    ) -> str:
        """
        一步生成带有新令牌的 API URL
        
        Args:
            api_key: API 密钥
            api_type: API 类��� (chat_completions/images_generations/embeddings)
            max_calls: 最大调用次数
            expire_days: 过期天数
            provider: API 提供者 (openai/dashscope)
            ext_info: 扩展信息
            
        Returns:
            str: 完整的 API URL，包含新生成的令牌参数
        """
        # 创建新令牌
        token = self.create_token(
            api_key=api_key,
            max_calls=max_calls,
            expire_days=expire_days,
            provider=provider,
            ext_info=ext_info
        )
        
        # 生成完整 URL
        return self.generate_api_url_with_token(token, api_type)