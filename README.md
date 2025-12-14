# 🪙 NumismaticApp

<div align="center">

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/tparicio/numismaticapp?sort=semver&logo=docker&label=Docker%20Hub)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/tparicio/numismaticapp/main)
![Vue Version](https://img.shields.io/badge/vue-3.x-42b883.svg?logo=vue.js)

**Gestiona tu colección de monedas con el poder de la Inteligencia Artificial.**

[Ver en DockerHub](https://hub.docker.com/r/tparicio/numismaticapp) • [Reportar Bug](https://github.com/antonioparicio/numismaticapp/issues)

</div>

---

## 📋 Descripción

**NumismaticApp** es una aplicación web moderna diseñada para coleccionistas de monedas.
- **Integración con Numista API:** Obtén datos detallados y referencias cruzadas de tus monedas directamente desde la mayor base de datos numismática.
- **Identificación con IA:** Sube fotos de tus monedas y deja que Google Gemini AI identifique y extraiga los detalles clave (país, año, valor, etc.).
Utiliza la IA de **Google Gemini** para analizar fotografías de monedas, extraer automáticamente metadatos (país, año, valor, ceca) y evaluar su estado de conservación.

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
| **Background Removal** | **rembg** | Eliminación de fondo inteligente. |
| **Datos Externos** | **Numista API** | Información y catálogos de monedas. |

## 🚀 Guía de Instalación

### Prerrequisitos

*   Docker & Docker Compose
*   Una [API Key de Google Gemini](https://aistudio.google.com/app/apikey)
*   Una [API Key de Numista](https://en.numista.com/api/doc/) (Opcional, para enriquecer datos)

### Opción 1: Docker Compose (Recomendado)

La forma más rápida de empezar es utilizando la imagen pre-construida desde DockerHub.

1.  **Crea el directorio de storage (importante para persistencia):**

    ```bash
    mkdir -p ./storage
    ```

    **Para despliegues en NAS (Synology, QNAP, etc.):** Si necesitas usar un UID/GID específico, asegúrate de que el directorio tenga los permisos correctos:

    ```bash
    # Ejemplo para UID 1000 y GID 100
    mkdir -p ./storage
    chown 1000:100 ./storage
    chmod 755 ./storage
    ```

2.  **Crea un archivo `docker-compose.yml`:**

    ```yaml
    services:
      app:
        image: tparicio/numismaticapp:latest
        # Para NAS: especifica tu UID:GID
        user: "1000:100"
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
          # IMPORTANTE: Monta el directorio storage para persistencia
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

3.  **Inicia la aplicación:**

    ```bash
    docker compose up -d
    ```

4.  **Accede al navegador:**
    *   Abre `http://localhost:8080` para ver tu colección.

### Opción 2: Compilación Local

Si prefieres compilar desde el código fuente:

1.  **Clona el repositorio:**
    ```bash
    git clone https://github.com/tparicio/numismaticapp.git
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

### Problemas de Persistencia y Permisos en NAS (Synology, QNAP, etc.)

#### Problema: Se borran los datos al re-desplegar
**Causa:** El directorio `storage` no está correctamente montado como volumen persistente.

**Solución:**
1. Asegúrate de crear el directorio `storage` en el host **antes** del primer despliegue:
   ```bash
   mkdir -p ./storage
   ```

2. Verifica que el `docker-compose.yml` incluya el volumen:
   ```yaml
   volumes:
     - ./storage:/app/storage
   ```

3. **IMPORTANTE:** No borres el directorio `./storage` del host al re-desplegar. Solo ejecuta `docker compose down` y `docker compose up -d`.

#### Problema: Permisos incorrectos después de re-desplegar
**Causa:** El contenedor se ejecuta con un UID/GID específico (ej: `1000:100`) pero el directorio `storage` tiene permisos diferentes.

**Solución para NAS:**
1. Identifica tu UID y GID en el NAS:
   ```bash
   id
   # Ejemplo de salida: uid=1000(usuario) gid=100(users)
   ```

2. Configura el directorio `storage` con los permisos correctos:
   ```bash
   chown 1000:100 ./storage
   chmod 755 ./storage
   ```

3. Añade la directiva `user` en tu `docker-compose.yml`:
   ```yaml
   services:
     app:
       image: tparicio/numismaticapp:latest
       user: "1000:100"  # Usa tu UID:GID específico
       volumes:
         - ./storage:/app/storage
   ```

4. Re-despliega:
   ```bash
   docker compose down
   docker compose up -d
   ```

#### Verificación de Permisos
Después del despliegue, verifica que el contenedor puede escribir en storage:
```bash
# Verifica permisos del directorio
ls -la ./storage

# Prueba de escritura desde el contenedor
docker compose exec app touch /app/storage/test.txt
ls -la ./storage/test.txt
```

### Problemas de Permisos en Linux (Desarrollo Local)
Si experimentas errores como `permission denied` en desarrollo local:

**Solución Recomendada:**
Usa las variables de entorno para que el contenedor use tu UID/GID actual:

```yaml
services:
  app:
    image: tparicio/numismaticapp:latest
    user: "${UID}:${GID}"
    # ... resto de la configuración
```

Luego exporta las variables antes de levantar el contenedor:
```bash
export UID=$(id -u)
export GID=$(id -g)
docker compose up -d
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
  Hecho con ❤️ por <a href="https://github.com/tparicio">Toni Paricio</a> con ayuda de 🚀 Antigravity y ✨ Gemini</div>
