import json
import base64
import hashlib
from datetime import datetime, timedelta, timezone
from typing import Dict, Any, Optional, Union

from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from Crypto.Util.Padding import pad

class TokenGenerator:
    """令牌生成器，用于创建和加密访问令牌"""

    def __init__(self, config: Union[str, Dict[str, Any]] = "default.rai"):
        """
        初始化令牌生成器
        
        Args:
            config: .rai 配置文件路径或配置字典对象
        """
        # 读取配置
        if isinstance(config, str):
            with open(config, 'r') as f:
                self.config = json.load(f)
        else:
            self.config = config
        
        # 获取加密配置
        self.crypto_config = self.config['crypto']
        if self.crypto_config['method'] != 'aes':
            raise ValueError("Only AES encryption is supported")
        
        # 解码 AES 密钥和 IV 种子
        self.key = bytes.fromhex(self.crypto_config['aes_key'])
        self.iv_seed = self.crypto_config['aes_iv_seed'].encode()

        # 生成配置哈希
        self.hash = self._generate_config_hash()

    def _generate_config_hash(self) -> str:
        """
        生成配置哈希
        
        Returns:
            str: 配置哈希值
        """
        data = (
            self.crypto_config['method'] +
            self.crypto_config['aes_key'] +
            self.crypto_config['aes_iv_seed']
        ).encode()
        return hashlib.sha256(data).hexdigest()

    def generate_iv(self) -> bytes:
        """
        生成初始化向量 (IV)
        
        Returns:
            bytes: 生成的 IV
        """
        iv = get_random_bytes(AES.block_size)
        # 使用 IV 种子进行混合
        return bytes(a ^ b for a, b in zip(iv, self.iv_seed))

    def create_token(
        self,
        api_key: str,
        max_calls: int = 100,
        expire_days: int = 1,
        provider: str = "dashscope",
        ext_info: str = ""
    ) -> Dict[str, Any]:
        """
        创建令牌数据
        
        Args:
            api_key: API 密钥
            max_calls: 最大调用次数
            expire_days: 过期天数
            provider: API 提供者 (openai/dashscope)
            ext_info: 扩展信息
            
        Returns:
            Dict[str, Any]: 令牌数据
        """
        now = datetime.now(timezone.utc)
        return {
            "id": f"token-{now.timestamp()}",
            "api_key": api_key,
            "max_calls": max_calls,
            "expire_time": (now + timedelta(days=expire_days)).isoformat(),
            "created_at": now.isoformat(),
            "provider": provider,
            "ext_info": ext_info
        }

    def encrypt_token(self, token_data: Dict[str, Any]) -> str:
        """
        加密令牌数据
        
        Args:
            token_data: 令牌数据
            
        Returns:
            str: 加密后的令牌字符串
        """
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
        return base64.urlsafe_b64encode(result).decode().rstrip('=')

    def get_server_url(self, path: Optional[str] = None) -> str:
        """
        获取服务器 URL
        
        Args:
            path: API 路径（可选）
            
        Returns:
            str: 完整的服务器 URL
        """
        base_url = f"{self.config['server']['host']}:{self.config['server']['port']}"
        if not path:
            return base_url
        
        # 确保路径以斜杠开头
        if not path.startswith('/'):
            path = '/' + path
            
        return f"{base_url}{self.config['server']['base_path'].rstrip('/')}{path}"