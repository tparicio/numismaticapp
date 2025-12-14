# 🪙 NumismaticApp

<div align="center">

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/tparicio/numismaticapp?sort=semver&logo=docker&label=Docker%20Hub)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/antonioparicio/numismaticapp)
![Vue Version](https://img.shields.io/badge/vue-3.x-42b883.svg?logo=vue.js)

**Gestiona tu colección de monedas con el poder de la Inteligencia Artificial.**

[Ver en DockerHub](https://hub.docker.com/r/tparicio/numismaticapp) • [Reportar Bug](https://github.com/antonioparicio/numismaticapp/issues)

</div>

---

## 📋 Descripción

**NumismaticApp** es una aplicación web moderna diseñada para coleccionistas de monedas. Utiliza la IA de **Google Gemini** para analizar fotografías de monedas, extraer automáticamente metadatos (país, año, valor, ceca) y evaluar su estado de conservación.

Olvídate de introducir datos manualmente. Simplemente sube una foto de tu moneda y deja que la IA haga el trabajo pesado, organizando tu colección en una base de datos segura y presentándola en un galería visualmente atractiva.

## ✨ Características Principales

*   **🤖 Análisis con IA:** Identificación automática de monedas y evaluación de grado (estado de conservación) mediante Google Gemini Vision.
*   **🖼️ Procesamiento de Imagen:** Recorte automático a círculo, eliminación de fondo y rotación inteligente con `libvips`.
*   **📁 Gestión de Colección:** Crea, edita y organiza tus monedas en grupos personalizados.
*   **📊 Dashboard Interactivo:** Visualiza estadísticas de tu colección, distribución por países, materiales y valor total.
*   **🔍 Integración con Numista:** Enlaza tus monedas con la base de datos de Numista para obtener información detallada y referencias.
*   **📱 Diseño Responsivo:** Interfaz moderna y adaptable construida con Vue 3 y DaisyUI.
*   **🐳 Docker Ready:** Despliegue sencillo y consistente mediante contenedores Docker.

## 🛠️ Tecnologías

| Componente | Tecnología | Descripción |
| :--- | :--- | :--- |
| **Backend** | ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) | API RESTful rápida y eficiente con Fiber. |
| **Frontend** | ![Vue.js](https://img.shields.io/badge/Vue.js-35495E?style=flat-square&logo=vue.js&logoColor=4FC08D) | SPA reactiva y ligera. |
| **Base de Datos** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat-square&logo=postgresql&logoColor=white) | Persistencia robusta y relacional. |
| **AI** | ![Gemini](https://img.shields.io/badge/Google%20Gemini-8E75B2?style=flat-square&logo=google&logoColor=white) | Motor de análisis visual. |
| **Imágenes** | **libvips** | Procesamiento de imágenes de alto rendimiento. |

## 🚀 Guía de Instalación

### Prerrequisitos

*   Docker & Docker Compose
*   Una [API Key de Google Gemini](https://aistudio.google.com/app/apikey)
*   Una [API Key de Numista](https://en.numista.com/api/doc/) (Opcional, para enriquecer datos)

### Opción 1: Docker Compose (Recomendado)

La forma más rápida de empezar es utilizando la imagen pre-construida desde DockerHub.

1.  **Crea un archivo `docker-compose.yml`:**

    ```yaml
    services:
      app:
        image: tparicio/numismaticapp:latest
        ports:
          - "8080:8080"
        environment:
          - GEMINI_API_KEY=tu_api_key_aqui
          - NUMISTA_API_KEY=tu_api_key_numista_opcional
          - REMBG_URL=http://rembg:5000/api/remove
          - POSTGRES_HOST=db
          - POSTGRES_USER=postgres
          - POSTGRES_PASSWORD=secret
          - POSTGRES_DB=numismatic
        depends_on:
          db:
            condition: service_healthy
          rembg:
            condition: service_healthy
        volumes:
          - ./storage:/app/storage

      db:
        image: postgres:15-alpine
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
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:5000"]
          interval: 10s
          timeout: 5s
          retries: 5
          start_period: 10s
        ports:
          - "5000:5000"

    volumes:
      postgres_data:
    ```

2.  **Inicia la aplicación:**

    ```bash
    docker compose up -d
    ```

3.  **Accede al navegador:**
    *   Abre `http://localhost:8080` para ver tu colección.

### Opción 2: Compilación Local

Si prefieres compilar desde el código fuente:

1.  **Clona el repositorio:**
    ```bash
    git clone https://github.com/antonioparicio/numismaticapp.git
    cd numismaticapp
    ```

2.  **Configura el entorno:**
    Crea un archivo `.env` en la raíz:
    ```bash
    GEMINI_API_KEY=tu_api_key_aqui
    NUMISTA_API_KEY=tu_api_key_numista_opcional
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=numismatic
    ```

3.  **Ejecuta con Make:**
    ```bash
    make run
    ```
    Esto levantará los servicios usando el `docker-compose.yml` de desarrollo incluido en el proyecto.

## 📖 Uso

### Añadir una Moneda

1.  Ve a la sección **"Añadir Moneda"**.
2.  Sube una foto del **Anverso** y otra del **Reverso**.
3.  Selecciona el grupo (opcional) o crea uno nuevo.
4.  Haz clic en **"Analizar y Guardar"**.
5.  La IA procesará las imágenes y rellenará los datos automáticamente.

### Gestionar Grupos

1.  Ve a la sección **"Grupos"**.
2.  Crea colecciones temáticas (ej: "Pesetas de Juan Carlos I", "Dólares de Plata").
3.  Asigna tus monedas a estos grupos para mantener tu colección organizada.

## ❓ Solución de Problemas

### Problemas de Permisos en Linux / NAS
Si experimentas errores como `permission denied` al intentar guardar imágenes en `storage/`, es probable que el usuario dentro del contenedor (`appuser`, UID normalment 1000) no tenga permisos de escritura en la carpeta montada desde el host.

**Solución Recomendada:**
Asegúrate de que el contenedor se ejecute con el mismo UID/GID que tu usuario actual en el host. Modifica tu `docker-compose.yml` añadiendo la directiva `user`:

```yaml
services:
  app:
    image: tparicio/numismaticapp:latest
    user: "${UID}:${GID}" # Usa el UID y GID de tu usuario actual
    # ... resto de la configuración
```

Luego, crea un archivo `.env` o exporta las variables antes de levantar el contenedor:
```bash
export UID=$(id -u)
export GID=$(id -g)
docker compose up -d
```

**Solución Alternativa:**
Cambia el propietario de la carpeta `storage` en tu host para que coincida con el usuario del contenedor (o dale permisos de escritura a todos `chmod 777 storage` - no recomendado para producción).

```bash
chown -R 1000:1000 ./storage
```

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Si tienes ideas para mejorar la aplicación:

1.  Haz un Fork del proyecto.
2.  Crea una rama con tu nueva funcionalidad (`git checkout -b feature/AmazingFeature`).
3.  Haz Commit de tus cambios (`git commit -m 'Add some AmazingFeature'`).
4.  Haz Push a la rama (`git push origin feature/AmazingFeature`).
5.  Abre un Pull Request.

## 📄 Licencia

Distribuido bajo la licencia MIT. Ver el archivo `LICENSE` para más información.

---

<div align="center">
  Hecho con ❤️ por <a href="https://github.com/antonioparicio">Antonio Aparicio</a>
</div>
