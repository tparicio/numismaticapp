# Infrastructure Documentation

This section details the external services and infrastructure components required to run the Numismatic App.

## External Services

### Google Gemini AI
The core intelligence of the application.
- **Purpose**: Analyzes coin images to extract metadata (Country, Year, Value, etc.).
- **Integration**: `internal/infrastructure/gemini/client.go` using the official Google Generative AI SDK.
- **Configuration**:
    - `GEMINI_API_KEY`: API Key.
    - `GEMINI_MODEL`: Model name (e.g., `gemini-1.5-flash`).

### Rembg
An external service for background removal.
- **Purpose**: Removes the background from coin photos to create clean cutouts.
- **Integration**: `internal/infrastructure/image/rembg.go` sends HTTP requests to a local or dockerized `rembg` server.
- **URL**: Configured via `REMBG_URL`.

## Storage
The application currently supports **Local Filesystem** storage.
- **Path**: Configurable, defaults to `./storage`.
- **Structure**:
    - `/original`: Full resolution uploads.
    - `/crop`: Processed images.
    - `/thumbnails`: Optimization for UI.
- **Future**: Interface design allows easy swapping for S3 or GCS.

## Containerization
The application is designed to be containerized using Docker.

### `docker-compose`
A typical stack includes:
1.  **App**: The Go API server.
2.  **Postgres**: Database.
3.  **Rembg**: Python service for background removal.
4.  **Frontend**: Static file serving (bundled with App or separate Nginx).

## Observability
- **Logging**: Fiber's built-in logger is used for request logging.
- **Health Check**: `GET /api/v1/health` checks database connectivity.
