# RelayAPI Python SDK

RelayAPI Python SDK is a client library for interacting with the RelayAPI server. It provides simple interfaces for generating API URLs, creating tokens, and sending various API requests.

## Installation

# Install from source

```bash
cd relayapi/backend-sdk/python
pip install -e .
```

Install using pip (coming soon):

```bash
pip install relayapi-sdk
```

## Configuration

The SDK requires a configuration object for initialization. You can load the configuration from a file (`.rai`) or pass the configuration object directly. Example configuration format:

```python
config = {
    "version": "1.0.0",
    "server": {
        "host": "http://localhost",
        "port": 8080,
        "base_path": "/relayapi/"
    },
    "crypto": {
        "method": "aes",
        "aes_key": "your-aes-key",
        "aes_iv_seed": "your-iv-seed"
    }
}
```

## Usage Examples

### Basic Usage

```python
from relayapi import RelayAPIClient
from openai import OpenAI

# Create client instance (using config object)
client = RelayAPIClient(config)

# Create token
token = client.create_token(
    api_key="your-api-key",
    max_calls=100,
    expire_days=1,
    provider="openai"
)

# Generate API URL
base_url = client.generate_api_url_with_token(token)
print("Base URL:", base_url)
# Output example: http://localhost:8080/relayapi/?token=xxxxx&rai_hash=xxxxx

# Use this URL as the base URL for OpenAI API in frontend code
openai_client = OpenAI(
    base_url=base_url,
    api_key="not-needed"  # Actual API key is included in the token
)
```

### Chat Request

```python
response = client.chat_completions(
    token=token,
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "What is the capital of France?"}
    ],
    model="gpt-3.5-turbo"
)
```

### Image Generation

```python
response = client.images_generations(
    token=token,
    prompt="A beautiful sunset over Paris",
    model="dall-e-3",
    size="1024x1024",
    quality="standard",
    n=1
)
```

### Embedding Generation

```python
response = client.embeddings(
    token=token,
    input="The quick brown fox jumps over the lazy dog",
    model="text-embedding-ada-002"
)
```

### Health Check

```python
status = client.health_check()
```

### URL Generation

The `generate_api_url_with_token` method is used to generate a complete API URL with token and hash parameters:

```python
# Generate base URL (without specifying API type)
base_url = client.generate_api_url_with_token(token)

# Generate URLs for specific APIs
chat_url = client.generate_api_url_with_token(token, 'chat_completions')
image_url = client.generate_api_url_with_token(token, 'images_generations')
embedding_url = client.generate_api_url_with_token(token, 'embeddings')
```

Parameters:
- `token` (str): Encrypted token string
- `api_type` (str, optional): API type, defaults to empty string

The method automatically adds the token and configuration hash as URL parameters.

## API Reference

### RelayAPIClient

#### Constructor

```python
RelayAPIClient(config: Union[str, Dict[str, Any]] = "default.rai")
```

- `config`: Configuration file path (string) or configuration object (dictionary)

#### Methods

##### create_token

Create and encrypt access token.

```python
create_token(
    api_key: str,
    max_calls: int = 100,
    expire_days: int = 1,
    provider: str = "dashscope",
    ext_info: str = ""
) -> str
```

##### generate_api_url_with_token

Generate complete API URL.

```python
generate_api_url_with_token(
    token: str,
    api_type: str = ""
) -> str
```

##### chat_completions

Send chat request.

```python
chat_completions(
    token: str,
    messages: List[Dict[str, str]],
    model: str = "gpt-3.5-turbo",
    **kwargs: Any
) -> Dict[str, Any]
```

##### images_generations

Generate images.

```python
images_generations(
    token: str,
    prompt: str,
    n: int = 1,
    size: str = "1024x1024",
    **kwargs: Any
) -> Dict[str, Any]
```

##### embeddings

Generate embedding vectors.

```python
embeddings(
    token: str,
    input: Union[str, List[str]],
    model: str = "text-embedding-ada-002",
    **kwargs: Any
) -> Dict[str, Any]
```

##### health_check

Check server health status.

```python
health_check() -> Dict[str, Any]
```

## Error Handling

All methods in the SDK will throw exceptions when errors occur. It's recommended to use try-except blocks to handle potential errors:

```python
try:
    response = client.chat_completions(...)
except Exception as e:
    print(f"Error: {e}")
```

## Example Programs

Check the example programs in the `examples` directory for more usage methods:

- `chat.py`: Demonstrates how to use the SDK for chat
- `url_generation.py`: Demonstrates how to generate and use API URLs

## License

MIT 