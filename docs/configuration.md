# RelayAPI Configuration Guide

RelayAPI uses two types of configuration files to manage server and client settings. This guide explains their purposes, formats, and usage.

## Server Configuration (`config.json`)

The server configuration file is **required** for running the RelayAPI server. It contains all necessary settings for server operation, including network settings, rate limits, and logging options.

### Basic Usage

```bash
# Start server with default config file
relayapi

# Start server with custom config file
relayapi --config=/path/to/config.json
```

### Configuration Structure

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8840,
    "read_timeout": 30,
    "write_timeout": 30,
    "max_header_bytes": 1048576,
    "debug": false
  },
  "log": {
    "console": true,
    "database": {
      "enabled": true,
      "type": "postgres",
      "connection_string": "user=postgres password=postgres dbname=relayapi host=localhost port=5432 sslmode=disable"
    },
    "web": {
      "enabled": false,
      "callback_url": "http://example.com/log"
    },
    "parquet": {
      "enabled": false,
      "file_path": "/path/to/logs/output.parquet"
    }
  },
  "rate_limit": {
    "requests_per_second": 20,
    "burst": 40,
    "ip_limit": {
      "requests_per_second": 10,
      "burst": 20
    }
  }
}
```

### Configuration Options

#### Server Settings
- `host`: Server listening address
- `port`: Server listening port
- `read_timeout`: Request read timeout in seconds
- `write_timeout`: Response write timeout in seconds
- `max_header_bytes`: Maximum size of request headers
- `debug`: Enable debug mode (default: false)

#### Logging Settings
- `console`: Enable console logging
- `database`: Database logging configuration
  - `enabled`: Enable database logging
  - `type`: Database type (postgres, mysql, sqlite)
  - `connection_string`: Database connection string
- `web`: Web callback logging configuration
  - `enabled`: Enable web callback logging
  - `callback_url`: URL for log callbacks
- `parquet`: Parquet file logging configuration
  - `enabled`: Enable Parquet file logging
  - `file_path`: Path to save Parquet log files

#### Rate Limiting
- `requests_per_second`: Global request rate limit
- `burst`: Global burst limit
- `ip_limit`: Per-IP rate limiting
  - `requests_per_second`: Per-IP request rate limit
  - `burst`: Per-IP burst limit

## Client Configuration (`default.rai`)

The client configuration file contains settings for SDK operation, including encryption settings and server connection information. If not present, a default configuration will be auto-generated.

### Basic Usage

```typescript
// Load from file
const client = new RelayAPIClient('default.rai');

// Or pass configuration object directly
const client = new RelayAPIClient({
  version: "1.0.0",
  server: {
    host: "http://localhost",
    port: 8840,
    base_path: "/relayapi/"
  },
  crypto: {
    method: "aes",
    aes_key: "your-aes-key",
    aes_iv_seed: "your-iv-seed"
  }
});
```

### Configuration Structure

```json
{
  "version": "1.0.0",
  "server": {
    "host": "http://localhost",
    "port": 8840,
    "base_path": "/relayapi/"
  },
  "crypto": {
    "method": "aes",
    "aes_key": "your-aes-key",
    "aes_iv_seed": "your-iv-seed"
  }
}
```

### Configuration Options

#### Version
- `version`: Configuration version (current: 1.0.0)

#### Server Connection
- `server.host`: RelayAPI server host
- `server.port`: RelayAPI server port
- `server.base_path`: Base path for API endpoints

#### Encryption Settings
- `crypto.method`: Encryption method (currently supports: aes)
- `crypto.aes_key`: AES encryption key
- `crypto.aes_iv_seed`: AES IV seed for encryption

### Auto-generation Logic

The `default.rai` file will be automatically generated if:
1. The file doesn't exist in the current directory
2. The SDK is initialized without a configuration object
3. The server is started without specifying a client configuration file

The auto-generated configuration will use secure random values for encryption keys and default server settings.

## Configuration Loading Priority

1. Command-line arguments (--config, --rai)
2. Environment variables
3. Configuration files in the current directory
4. Auto-generated default configuration

## Security Considerations

1. Keep your `config.json` and `default.rai` files secure
2. Never commit encryption keys to version control
3. Use environment variables for sensitive information in production
4. Regularly rotate encryption keys
5. Use strong, random encryption keys in production

## Best Practices

1. Use different configurations for development and production
2. Keep backups of your configuration files
3. Document any custom configuration changes
4. Monitor rate limits in production
5. Regularly review and update configurations

For more specific configuration examples and use cases, please refer to:
- [Server Configuration Guide](../server/README.md)
- [JavaScript SDK Guide](../backend-sdk/JavaScript/README.md)
- [Python SDK Guide](../backend-sdk/python/README.md) 