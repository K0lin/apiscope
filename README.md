# APIScope

A modern, web-based OpenAPI documentation platform built with Go and Gin. Upload, version, and share your API specifications with beautiful, interactive Swagger UI documentation.

## Features

- **📤 Easy Upload**: Upload OpenAPI/Swagger files or paste YAML/JSON content directly
- **🔗 Shareable Links**: Generate permanent, shareable links for your API documentation
- **📋 Version Control**: Support for multiple versions of your API specifications
- **👁️ Interactive Docs**: Beautiful Swagger UI for testing and exploring your APIs
- **🛠️ SDK Generation**: Generate client SDKs in multiple programming languages using OpenAPI Generator
- **💾 File Storage**: Secure local file storage with automatic organization
- **🗄️ Database Integration**: Redis for fast metadata and version tracking
- **🎨 Modern UI**: Clean, professional interface with responsive design

## Quick Start

### Prerequisites

- Go 1.19 or later
- Redis 6.0 or later
- Git

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/k0lin/apiscope.git
   cd apiscope
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Start Redis:**

   ```bash
   redis-server
   ```

4. **Configure environment (optional):**
   Create a `.env` file in the root directory:

   ```env
   PORT=8080
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   STORAGE_PATH=./storage/documents

   # OpenAPI Generator Configuration
   OPENAPI_GENERATOR_ENABLED=true
   OPENAPI_GENERATOR_SERVER=https://api.openapi-generator.tech
   ```

5. **Run the application:**

   ```bash
   go run cmd/server/main.go
   ```

6. **Open your browser:**
   Navigate to `http://localhost:8080`

## Usage

### Uploading API Specifications

1. **Via File Upload:**
   - Click "Upload File" tab
   - Drag and drop or click to select your OpenAPI YAML/JSON file
   - Optionally set a custom name, description, and version
   - Click "Generate Documentation Link"

2. **Via Content Paste:**
   - Click "Paste Content" tab
   - Paste your OpenAPI specification directly into the text area
   - Optionally set metadata
   - Click "Generate Documentation Link"

### Managing Versions

- View all versions of a document from the document viewer
- Add new versions using the "Add New Version" button
- Each version gets a unique URL for sharing

### SDK Generation

When enabled, APIScope provides built-in SDK generation capabilities:

1. **Enable SDK Generation:**
   - Set `OPENAPI_GENERATOR_ENABLED=true` in your `.env` file
   - Configure `OPENAPI_GENERATOR_SERVER` to point to your OpenAPI Generator instance

2. **Generate SDKs:**
   - Navigate to any document viewer page
   - Select a programming language from the SDK dropdown
   - Enter a package name for your generated SDK
   - Click "Generate & Download" to create and download the SDK

3. **Supported Languages:**
   - Python, Java, JavaScript, TypeScript, Go, PHP, Ruby, C#, and many more
   - Full list depends on your OpenAPI Generator server configuration

**Note:** For the SDK generation to work properly when using localhost, ensure your OpenAPI Generator server can reach your APIScope instance. Consider using network IP addresses instead of localhost when deploying.

### API Endpoints

APIScope provides REST API endpoints for programmatic access:

- `GET /api/document/{id}/content` - Get document content
- `GET /api/document/{id}/content?version={version}` - Get specific version
- `GET /api/document/{id}/versions` - List all versions

## Project Structure

```text
apiscope/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and setup
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
├── storage/documents/   # File storage directory
├── web/
│   ├── static/          # CSS, JS, and assets
│   └── templates/       # HTML templates
├── go.mod
├── go.sum
└── README.md
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `REDIS_ADDR` | `localhost:6379` | Redis server address |
| `REDIS_PASSWORD` | `` | Redis password (if required) |
| `STORAGE_PATH` | `./storage/documents` | File storage directory |
| `OPENAPI_GENERATOR_ENABLED` | `false` | Enable/disable SDK generation feature |
| `OPENAPI_GENERATOR_SERVER` | `` | URL of OpenAPI Generator server |

### File Upload Limits

- **Maximum file size**: 50MB
- **Supported formats**: YAML (.yaml, .yml), JSON (.json)
- **OpenAPI versions**: 2.x (Swagger) and 3.x supported

## Development

### Building

```bash
go build -o apiscope cmd/server/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Style

This project follows standard Go conventions and uses:

- `gofmt` for code formatting
- Standard Go project layout
- Gin web framework
- GORM for database operations

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Technologies Used

- **Backend**: Go, Gin web framework
- **Database**: Redis for high-performance data storage
- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **API Documentation**: Swagger UI
- **SDK Generation**: OpenAPI Generator integration
- **File Processing**: YAML/JSON parsing

## Support

If you find this project helpful, please consider:

- ⭐ Starring the repository
- 🐛 Reporting bugs or issues
- 💡 Suggesting new features
- 📖 Contributing to the documentation

## Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Redis](https://redis.io/)
- [Go YAML library](https://gopkg.in/yaml.v3)
