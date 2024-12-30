# RelayAPI Server

[中文文档](README_CN.md)

The RelayAPI server is the core component responsible for API proxying, token validation, and request forwarding.

## Quick Start

```bash
# Clone the repository
git clone https://github.com/relayapi/RelayAPI.git

# Enter the server directory
cd server

# Run the server with default configuration
go run cmd/server/main.go -rai ./rai -d
```

Command line options:
- `-rai`: Client configuration directory path (default: current directory)
- `-config`: Server configuration file path (default: config.json)
- `-d`: Enable debug mode, logs will be written to debug.log

## Configuration

### Server Configuration (`config.json`)

The server configuration file controls server behavior, including network settings, rate limits, and logging options.

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

### Client Configuration (`.rai` files)

Client configuration files contain encryption settings used by both the server and SDK. The server monitors these files in the `-rai` directory.

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

You can customize the `crypto` section as needed, just ensure the server and SDK use the same configuration.

## Development

### Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go         # Server entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── crypto/            # Encryption implementation
│   ├── handlers/          # Request handlers
│   ├── middleware/        # Middleware components
│   ├── models/            # Data models
│   └── services/          # Business logic
├── rai/                   # Client configuration files
└── config.json            # Server configuration
```

### Core Components

1. **Configuration Management** (`internal/config/`)
   - Load and validate server configuration
   - Monitor client configuration files
   - Handle configuration updates

2. **Encryption** (`internal/crypto/`)
   - Implement AES encryption/decryption
   - Manage encryption keys
   - Handle token generation and validation

3. **Request Handling** (`internal/handlers/`)
   - Handle incoming API requests
   - Validate tokens
   - Forward requests to AI providers

4. **Middleware** (`internal/middleware/`)
   - Authentication
   - Rate limiting
   - Logging
   - Request/Response transformation

### Adding New Features

1. **New AI Provider**
   ```go
   // internal/handlers/provider.go
   func (h *Handler) handleProviderRequest(c *gin.Context) {
       // Implement provider-specific handling logic
   }
   ```

2. **New Middleware**
   ```go
   // internal/middleware/custom.go
   func CustomMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // Implement middleware logic
       }
   }
   ```

3. **New Configuration Option**
   ```go
   // internal/config/config.go
   type Config struct {
       // Add new configuration fields
   }
   ```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/crypto
```

## Deployment

1. **Build**
   ```bash
   go build -o relayapi cmd/server/main.go
   ```

2. **Run**
   ```bash
   ./relayapi -rai /path/to/rai/dir -config /path/to/config.json
   ```

3. **Monitor**
   - View `debug.log` for detailed logs when running with `-d`
   - Monitor server status through `/health` endpoint
   - Check console output for real-time statistics

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.