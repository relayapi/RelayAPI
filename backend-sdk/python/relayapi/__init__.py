"""
RelayAPI Python SDK
~~~~~~~~~~~~~~~~~~

RelayAPI Python SDK 是一个用于访问 RelayAPI Server 的客户端库。

基本用法:

    >>> from relayapi import RelayAPIClient
    >>> client = RelayAPIClient("default.rai")
    >>> token = client.create_token("your-api-key")
    >>> response = client.chat_completions(token, [{"role": "user", "content": "Hello!"}])
    >>> print(response["choices"][0]["message"]["content"])

:copyright: (c) 2024 RelayAPI Team
:license: MIT
"""

from .client import RelayAPIClient
from .token import TokenGenerator

__version__ = "1.0.0"
__all__ = ["RelayAPIClient", "TokenGenerator"]