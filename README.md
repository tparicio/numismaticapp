# 🪙 NumismaticApp

<div align="center">

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/tparicio/numismaticapp?sort=semver&logo=docker&label=Docker%20Hub)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/tparicio/numismaticapp/main)
![Vue Version](https://img.shields.io/badge/vue-3.x-42b883.svg?logo=vue.js)

**Manage your coin collection with the power of Artificial Intelligence.**

[View on DockerHub](https://hub.docker.com/r/tparicio/numismaticapp) • [Report Bug](https://github.com/antonioparicio/numismaticapp/issues)

</div>

---

## 📋 Description

**NumismaticApp** is a modern web application designed for coin collectors.
- **Numista API Integration:** Get detailed data and cross-references for your coins directly from the largest numismatic database.
- **AI Identification:** Upload photos of your coins and let Google Gemini AI identify and extract key details (country, year, value, mint, rule, etc.).
- **Automatic Grading:** Uses AI to estimate the conservation state (grade) of your coins.

Forget about manual data entry. Simply upload a photo of your coin and let AI do the heavy lifting, organizing your collection into a secure database and presenting it in a visually appealing gallery.

## ✨ Key Features

*   **🤖 AI Analysis:** Automatic coin identification and grading using Google Gemini Vision.
*   **🖼️ Image Processing:** Auto-crop to circle, background removal, and smart rotation using `vips` (Alpine optimized).
*   **📁 Collection Management:** Create, edit, and organize your coins into custom groups.
*   **📊 Interactive Dashboard:** Visualize collection statistics, distribution by country, material, and total value.
*   **🔍 Numista Integration:** Link your coins with the Numista database for detailed information and references.
*   **📱 Responsive Design:** Modern and adaptable interface built with Vue 3 and DaisyUI.
*   **🐳 Docker Ready:** Simple and consistent deployment using Docker containers (Secure Alpine-based image).

## 🛠️ Technology Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Backend** | ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) | Fast and efficient RESTful API with Fiber (Go 1.25). |
| **Frontend** | ![Vue.js](https://img.shields.io/badge/Vue.js-35495E?style=flat-square&logo=vue.js&logoColor=4FC08D) | Reactive and lightweight SPA. |
| **Database** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat-square&logo=postgresql&logoColor=white) | Robust relational persistence. |
| **AI** | ![Gemini](https://img.shields.io/badge/Google%20Gemini-8E75B2?style=flat-square&logo=google&logoColor=white) | Visual analysis engine. |
| **Images** | **libvips** | High-performance image processing (Alpine optimized). |
| **Background Removal** | **rembg** | Smart background removal. |
| **External Data** | **Numista API** | Coin information and catalogs. |

## 🚀 Installation Guide

### Prerequisites

*   Docker & Docker Compose
*   A [Google Gemini API Key](https://aistudio.google.com/app/apikey)
*   A [Numista API Key](https://en.numista.com/api/doc/) (Optional, for data enrichment)

### Option 1: Docker Compose (Recommended)

The fastest way to start is using the pre-built image from DockerHub.

1.  **Create a storage directory (important for persistence):**

    ```bash
    mkdir -p ./storage
    ```

    **For NAS deployments (Synology, QNAP, etc.):** The container runs as a non-root user (`appuser`, UID 1001), so ensure the directory has correct permissions:

    ```bash
    # Set owner to UID 1001 (appuser)
    mkdir -p ./storage
    chown 1001:1001 ./storage
    chmod 755 ./storage
    ```

2.  **Create a `docker-compose.yml` file:**

    ```yaml
    services:
      app:
        image: tparicio/numismaticapp:latest
        restart: unless-stopped
        user: "1001:1001" # Default appuser, but good to be explicit
        ports:
          - "8080:8080"
        environment:
          - GEMINI_API_KEY=your_api_key_here
          - NUMISTA_API_KEY=your_optional_numista_key
          - REMBG_URL=http://rembg:5000/api/remove
          - POSTGRES_HOST=db
          - POSTGRES_USER=postgres
          - POSTGRES_PASSWORD=secret
          - POSTGRES_DB=numismatic
        depends_on:
          db:
            condition: service_healthy
        volumes:
          - ./storage:/app/storage

      db:
        image: postgres:15-alpine
        restart: unless-stopped
        environment:
          - POSTGRES_USER=postgres
          - POSTGRES_PASSWORD=secret
          - POSTGRES_DB=numismatic
        volumes:
          - postgres_data:/var/lib/postgresql/data
        healthcheck:
          test: ["CMD-SHELL", "pg_isready -U postgres"]
          interval: 10s
          timeout: 5s
          retries: 5
          start_period: 30s

      rembg:
        image: danielgatis/rembg:latest
        command: s --host 0.0.0.0 --port 5000
        ports:
          - "5000:5000"

    volumes:
      postgres_data:
    ```

3.  **Start the application:**

    ```bash
    docker compose up -d
    ```

4.  **Access in browser:**
    *   Open `http://localhost:8080` to view your collection.

### Option 2: Local Build

If you prefer building from source:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/tparicio/numismaticapp.git
    cd numismaticapp
    ```

2.  **Configure environment:**
    Create a `.env` file in the root directory:
    ```bash
    GEMINI_API_KEY=your_api_key
    NUMISTA_API_KEY=your_optional_key
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=numismatic
    ```

3.  **Run with Make:**
    ```bash
    make run
    ```
    This will start services using the development `docker-compose.yml`.

## 📖 Usage

### Adding a Coin

1.  Go to **"Add Coin"** section.
2.  Upload a photo of the **Obverse** and **Reverse**.
3.  Select a group (optional) or create a new one.
4.  Click **"Analyze and Save"**.
5.  AI will process images and fill in data automatically.

### Managing Groups

1.  Go to **"Groups"** section.
2.  Create thematic collections (e.g., "Silver Dollars", "Ancient Rome").
3.  Assign your coins to these groups to keep your collection organized.

## ❓ Troubleshooting

### Persistence & Permissions on NAS (Synology, QNAP, etc.)

#### Issue: Data lost on re-deploy
**Cause:** `storage` directory is not correctly mounted as a persistent volume.

**Solution:**
1. Ensure the `storage` directory exists on the host **before** first deployment.
2. Verify `docker-compose.yml` includes the volume mapping:
   ```yaml
   volumes:
     - ./storage:/app/storage
   ```

#### Issue: Permission Denied errors
**Cause:** The container runs as non-root user `appuser` (UID 1001), but the host directory belongs to root or another user.

**Solution:**
1. Change ownership of the folder on the host machine:
   ```bash
   chown -R 1001:1001 ./storage
   ```
2. Or explicitly set the user in `docker-compose.yml` if you need a different UID:
   ```yaml
   user: "1026:100" # Example user:group in Synology
   ```

## 🤝 Contribution

Contributions are welcome! If you have ideas to improve the application:

1.  Fork the project.
2.  Create a feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Abre un Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` file for more information.

---

<div align="center">
  Made with ❤️ by <a href="https://github.com/tparicio">Toni Paricio</a> with help from 🚀 Antigravity & ✨ Gemini</div>
