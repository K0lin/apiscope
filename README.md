# APIScope

A modern, web-based OpenAPI documentation platform built with Go and Gin. Upload, validate, version, and share your API specifications with beautiful, interactive Swagger UI documentation. Includes optional SDK generation, version lifecycle operations, live server editing, and configurable CORS.

## Features

- **📤 Easy Upload**: Upload OpenAPI/Swagger files or paste YAML/JSON content directly.
- **🧪 Validation**: Early validation rejects malformed or structurally empty specs (ensures `openapi`/`swagger`, `info`, and minimal paths/components).
- **🔗 Shareable Links**: Generate permanent, shareable links for your API documentation.
- **📋 Version Control**: Multiple versions per document with automatic latest tracking and chronological ordering.
- **♻️ Version Deletion (Optional)**: Delete individual versions safely with automatic re‑promotion of newest remaining version.
- **⬇️ Version Download (Optional)**: Download the raw stored YAML for the currently selected version.
- **🛠️ SDK Generation**: Generate client SDKs in multiple languages via OpenAPI Generator (toggleable).
- **🧩 Live Servers Editing (Optional)**: Temporarily add/remove `servers` entries client‑side for quick local testing (non‑persistent) and download modified spec.
- **🔁 Auto Server Origin Adjust (Optional)**: When enabled, the first server entry matching the spec's original host:port is auto-rewritten to the current viewer origin (helps when specs hardcode a different localhost port).
- **🚫 Strip Servers (Optional)**: Completely remove all `servers` entries from displayed specs (read-only view, disables Try It Out requests, overrides server editing & auto-adjust).
- **🔐 One-Time Share Slug (Optional)**: Allow choosing a memorable or randomly generated share link `/share/{slug}` per document (immutable once set).
- **🌐 CORS Configuration**: Fine‑grained control over origins, methods, headers, credentials, and max age.
- **🗄️ File Storage**: Local organized storage per document/version ID.
- **⚡ Redis Metadata**: Fast document + version metadata tracking in Redis.
- **🩺 Health Endpoint**: Simple `/health` JSON endpoint for monitoring.
- **🎨 Modern UI**: Clean, responsive, minimal dependencies.

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

   # Feature Toggles
   OPENAPI_GENERATOR_ENABLED=true
   OPENAPI_GENERATOR_SERVER=https://api.openapi-generator.tech
   ALLOW_VERSION_DELETION=false
   ALLOW_VERSION_DOWNLOAD=true
   ALLOW_SERVER_EDITING=false
   AUTO_ADJUST_SERVER_ORIGIN=false
   STRIP_OPENAPI_SERVERS=false
   ALLOW_CUSTOM_SHARE_LINK=false

   # CORS
   ALLOWED_ORIGINS=*
   CORS_ALLOW_CREDENTIALS=false
   CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
   CORS_ALLOWED_HEADERS=Authorization,Content-Type,Accept,Origin
   CORS_EXPOSE_HEADERS=Content-Length
   CORS_MAX_AGE=600
   CORS_DEBUG=false
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

- Versions sorted newest-first.
- Latest is auto-flagged; adding a new version promotes it.
- Selecting an older version updates the view while preserving dropdown selection.
- (Optional) If `ALLOW_VERSION_DELETION=true`, a Delete button appears to remove the selected version (cannot undo). Latest is re‑assigned automatically if removed.
- (Optional) If `ALLOW_VERSION_DOWNLOAD=true`, a Download button provides the raw YAML file of the selected version.

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

### API Endpoints (Core)

APIScope provides REST API endpoints for programmatic access:

- `GET /api/document/{id}/content` – Get latest version content (YAML/JSON as originally stored)
- `GET /api/document/{id}/content?version={version}` – Get a specific version
- `GET /api/document/{id}/versions` – List all versions
- `DELETE /api/document/{id}/version/{version}` – (If enabled) delete specific version
- `GET /api/document/{id}/version/{version}/download` – (If enabled) download stored file
- `GET /health` – Health status JSON
- `POST /api/document/{id}/share` – (If enabled) set a one-time share slug (body: `{ "slug": "optional-custom" }`) returns `{ share_slug, url }`
- `GET /share/{slug}` – Resolve a share slug to the underlying document view (redirects to `/view/{id}`)

### Live Servers Editing (Client‑Side)
If `ALLOW_SERVER_EDITING=true` you can add/remove `servers` entries directly in the viewer for ad‑hoc testing (not persisted). You may then download the modified spec for local reuse.
If `STRIP_OPENAPI_SERVERS=true`, this feature is automatically disabled.

### Read-Only Mode (Strip Servers)
If `STRIP_OPENAPI_SERVERS=true`:

- All `servers` entries are removed from the rendered spec.
- Try It Out / Execute buttons are disabled (no outbound calls).
- `ALLOW_SERVER_EDITING` and `AUTO_ADJUST_SERVER_ORIGIN` are ignored.
- Ideal for public/internal sharing where execution should be blocked.

To combine with version download: You can still download raw stored specs (if `ALLOW_VERSION_DOWNLOAD=true`).

### CORS Configuration
Customize CORS via environment variables. Example tightened production config:
```env
ALLOWED_ORIGINS=https://docs.example.com,https://app.example.com
CORS_ALLOW_CREDENTIALS=true
```
Notes:
- When `CORS_ALLOW_CREDENTIALS=true`, avoid wildcard `*`; the middleware will echo the request origin.
- Preflight cache duration controlled by `CORS_MAX_AGE` (seconds).

### Health Check
`GET /health` returns a simple JSON body: `{ "status": "ok", "time": "<RFC3339>" }`.

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

### Core Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `REDIS_ADDR` | `localhost:6379` | Redis server address |
| `REDIS_PASSWORD` | (empty) | Redis password (if required) |
| `STORAGE_PATH` | `./storage/documents` | File storage root |
| `OPENAPI_GENERATOR_ENABLED` | `false` | Toggle SDK generation feature |
| `OPENAPI_GENERATOR_SERVER` | public generator URL | OpenAPI Generator server URL |
| `ALLOW_VERSION_DELETION` | `false` | Enable Delete Version button/API |
| `ALLOW_VERSION_DOWNLOAD` | `true` | Enable Download Version button/API |
| `ALLOW_SERVER_EDITING` | `false` | Enable client-side servers editor |
| `AUTO_ADJUST_SERVER_ORIGIN` | `false` | Auto-rewrite first server origin to current host/port |
| `STRIP_OPENAPI_SERVERS` | `false` | Strip all servers; disables Try It Out & overrides editing/auto-adjust |
| `ALLOW_CUSTOM_SHARE_LINK` | `false` | Permit one-time assignment of a custom or generated share slug `/share/{slug}` |
| `ALLOWED_ORIGINS` | `*` | Comma-separated allowed CORS origins |
| `CORS_ALLOW_CREDENTIALS` | `false` | Allow credentialed CORS requests |
| `CORS_ALLOWED_METHODS` | defaults list | Allowed CORS methods |
| `CORS_ALLOWED_HEADERS` | defaults list | Allowed CORS request headers |
| `CORS_EXPOSE_HEADERS` | `Content-Length` | Exposed response headers |
| `CORS_MAX_AGE` | `600` | Preflight cache seconds |
| `CORS_DEBUG` | `false` | Verbose CORS logging |

### File Upload & Validation

- **Maximum file size**: 50MB
- **Supported formats**: YAML (.yaml, .yml), JSON (.json)
- **OpenAPI versions**: Swagger 2.x and OpenAPI 3.x supported
- **Validation rules**: Ensures format correctness, required `info` fields, version field, and at least one of `paths` or `components`.

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

- `gofmt` for formatting
- Gin web framework
- Redis for metadata (no ORM layer / GORM removed from docs)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Technologies Used

- **Backend**: Go (Gin)
- **Data Store**: Redis (metadata & version tracking)
- **Storage**: Local filesystem (YAML/JSON specs)
- **Documentation UI**: Swagger UI
- **SDK Generation**: OpenAPI Generator (optional)
- **Validation**: YAML & JSON parsing + structural checks
- **UI**: HTML5, CSS3, Vanilla JS

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
