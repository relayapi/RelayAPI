# RelayAPI JavaScript SDK

RelayAPI JavaScript SDK is a client library for interacting with RelayAPI servers. It provides simple interfaces for generating API URLs, creating tokens, and sending various API requests.

## Installation

Install using npm:

```bash
npm install relayapi-sdk
```

## Configuration

The SDK requires a configuration object for initialization. You can load the configuration from a file (`.rai`) or pass it directly as an object. Example configuration format:

```json
{
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

```javascript
import { RelayAPIClient } from 'relayapi-sdk';
import fs from 'fs/promises';
import { OpenAI } from 'openai';

// Load configuration from file
const configContent = await fs.readFile('config.rai', 'utf-8');
const config = JSON.parse(configContent);

// Create client instance
const client = new RelayAPIClient(config);

// Create token
const token = client.createToken({
    apiKey: 'your-api-key',
    maxCalls: 100,
    expireSeconds: 3600,
    provider: 'openai'
});

// Generate API URL
const baseUrl = client.generateUrl(token);
console.log('Base URL:', baseUrl);
// Output example: http://localhost:8080/relayapi/?token=xxxxx&rai_hash=xxxxx

// Use this URL as OpenAI API base URL in your frontend code
const openai = new OpenAI({
    baseURL: baseUrl,
    apiKey: 'not-needed' // The actual API key is already in the token
});
```

### Chat Request

```javascript
const response = await client.chat({
    messages: [
        { role: 'system', content: 'You are a helpful assistant.' },
        { role: 'user', content: 'What is the capital of France?' }
    ],
    model: 'gpt-3.5-turbo',
    temperature: 0.7,
    maxTokens: 1000,
    token: token
});
```

### Image Generation

```javascript
const response = await client.generateImage({
    prompt: 'A beautiful sunset over Paris',
    model: 'dall-e-3',
    size: '1024x1024',
    quality: 'standard',
    n: 1,
    token: token
});
```

### Embedding Generation

```javascript
const response = await client.createEmbedding({
    input: 'The quick brown fox jumps over the lazy dog',
    model: 'text-embedding-ada-002',
    token: token
});
```

### Health Check

```javascript
const status = await client.healthCheck();
```

## API Reference

### RelayAPIClient

#### Constructor

```javascript
new RelayAPIClient(config)
```

- `config`: String (config file path) or Object (config object)

#### Methods

##### createToken(options)

Creates a new token.

- `options.apiKey`: API key
- `options.maxCalls`: Maximum number of calls (default: 100)
- `options.expireSeconds`: Seconds until expiration (default: 86400, 24 hours)
- `options.provider`: Provider name or URL. When a URL is provided, it will be used directly as the provider endpoint. Supported provider names: 'dashscope', 'openai', etc.
- `options.extInfo`: Extended information (optional)

##### generateUrl(endpoint, token)

Generates an API URL.

- `endpoint`: API endpoint path
- `token`: Token string

##### chat(options)

Sends a chat request.

- `options.messages`: Array of messages
- `options.model`: Model name (default: 'gpt-3.5-turbo')
- `options.temperature`: Temperature value (default: 0.7)
- `options.maxTokens`: Maximum tokens (default: 1000)
- `options.token`: Token string

##### generateImage(options)

Generates images.

- `options.prompt`: Image description
- `options.model`: Model name (default: 'dall-e-3')
- `options.size`: Image size (default: '1024x1024')
- `options.quality`: Image quality (default: 'standard')
- `options.n`: Number of images to generate (default: 1)
- `options.token`: Token string

##### createEmbedding(options)

Generates embeddings.

- `options.input`: Input text
- `options.model`: Model name (default: 'text-embedding-ada-002')
- `options.token`: Token string

##### healthCheck()

Checks server health status.

## Error Handling

All methods in the SDK will throw exceptions when errors occur. It's recommended to use try-catch blocks to handle potential errors:

```javascript
try {
    const response = await client.chat({...});
} catch (error) {
    console.error('Error:', error.message);
}
```

## Example Programs

Check the example programs in the `examples` directory for more usage examples:

- `chat.js`: Demonstrates how to use the SDK for chat, image generation, and embeddings

## License

MIT

### Generate URL

The `generateUrl` method is used to generate a complete API URL with token and hash parameters:

```javascript
const url = client.generateUrl(token);  // Use default empty endpoint
const url = client.generateUrl(token, 'chat/completions');  // Specify endpoint
```

Parameters:
- `token` (string): The encrypted token string
- `endpoint` (string, optional): The API endpoint path, defaults to empty string

The method will automatically add the token and configuration hash as URL parameters.
