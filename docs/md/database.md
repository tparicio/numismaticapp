# Database Schema Documentation

## Tables

### `coins`

Stores the main information about each coin in the collection.

| Column | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `country` | VARCHAR(255) | Country of origin |
| `year` | INTEGER | Year of minting |
| `face_value` | VARCHAR(100) | Face value (e.g., "1 Dollar") |
| `currency` | VARCHAR(100) | Currency unit |
| `material` | VARCHAR(100) | Material composition |
| `description` | TEXT | Detailed description |
| `km_code` | VARCHAR(50) | Krause-Mishler catalog code |
| `min_value` | DECIMAL(10, 2) | Estimated minimum value |
| `max_value` | DECIMAL(10, 2) | Estimated maximum value |
| `grade` | VARCHAR(50) | Condition grade |
| `sample_image_url_front` | TEXT | URL/Path to the front sample image |
| `sample_image_url_back` | TEXT | URL/Path to the back sample image |
| `notes` | TEXT | User notes |
| `gemini_details` | JSONB | Raw JSON data from Gemini analysis |
| `created_at` | TIMESTAMP | Creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

### `coin_images`

Stores metadata for all images associated with a coin (originals and processed).

| Column | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `coin_id` | UUID | Foreign Key to `coins.id` |
| `image_type` | ENUM | Type of image: `original_front`, `original_back`, `processed_front`, `processed_back` |
| `path` | VARCHAR(255) | File system path to the image |
| `extension` | VARCHAR(10) | File extension (e.g., "jpg") |
| `size` | BIGINT | File size in bytes |
| `width` | INTEGER | Image width in pixels |
| `height` | INTEGER | Image height in pixels |
| `mime_type` | VARCHAR(50) | MIME type (e.g., "image/jpeg") |
| `created_at` | TIMESTAMP | Creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

## Enums

### `image_type`
- `original_front`
- `original_back`
- `processed_front`
- `processed_back`
