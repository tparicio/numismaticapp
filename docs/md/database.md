# Database Documentation

The application uses **PostgreSQL** as its relational database. The schema is designed to support the core entity `Coin` and its related data like `Images` and `Groups`.

## ER Diagram

```mermaid
erDiagram
    COINS ||--o{ COIN_IMAGES : has
    GROUPS ||--o{ COINS : contains

    COINS {
        UUID id PK
        VARCHAR name
        UUID group_id FK
        VARCHAR mint
        BIGINT mintage
        VARCHAR country
        INTEGER year
        VARCHAR face_value
        VARCHAR currency
        VARCHAR material
        TEXT description
        VARCHAR km_code
        DECIMAL min_value
        DECIMAL max_value
        VARCHAR grade
        JSONB gemini_details
        TEXT personal_notes
        DATE acquired_at
        DATE sold_at
        NUMERIC price_paid
        NUMERIC sold_price
    }

    GROUPS {
        SERIAL id PK
        VARCHAR name
        TEXT description
    }

    COIN_IMAGES {
        UUID id PK
        UUID coin_id FK
        ENUM image_type "original, crop, thumbnail, sample"
        ENUM side "front, back"
        VARCHAR path
        VARCHAR mime_type
    }
```

## Tables

### `coins`
The central table storing all numismatic data.
- **Primary Key**: `id` (UUID v4)
- **JSONB**: `gemini_details` stores the raw analysis result from the AI model, allowing for schema-less flexibility for AI data.
- **Indexes**: `country`, `year` for faster filtering.

### `coin_images`
Stores metadata about the images associated with a coin.
- **Types**:
    - `original`: The raw upload.
    - `crop`: The background-removed, circular crop.
    - `thumbnail`: A smaller version for list views.
    - `sample`: Reference images (unused currently).

### `groups`
Simple categorization for coins (e.g., "My Gold Collection", "Swap List").

## Data Access Strategy

We use **sqlc** to generate type-safe Go code from SQL queries.
- **Queries Location**: `internal/infrastructure/db/queries/`
- **Generated Code**: `internal/infrastructure/db/`
- **Migration**: Schema is defined in `schema.sql`.

## Enums
- `image_type`: Ensures data integrity for image categorization.
- `coin_side`: `front` (obverse) or `back` (reverse).
