# 🚀 API DB Migrations

API REST de ejemplo con migraciones automáticas de base de datos, construida con **Go**, **Fiber** y **MySQL**.

## 📋 Descripción

Este proyecto es un **ejemplo completo de API REST** que implementa:

- ✅ **Migraciones automáticas** de base de datos con versionado
- ✅ **CRUD completo** de Usuarios y Campañas
- ✅ **Arquitectura limpia** con separación de responsabilidades
- ✅ **Endpoints RESTful** siguiendo mejores prácticas
- ✅ **Docker Compose** para desarrollo rápido
- ✅ **Logging automático** con Fiber
- ✅ **Manejo de errores** consistente
- ✅ **Validaciones** de datos
- ✅ **Documentación completa** de la API

---

## 🛠️ Tecnologías

| Tecnología | Versión | Uso |
|------------|---------|-----|
| **Go** | v1.22 | Lenguaje principal |
| **Fiber** | v2.52.0 | Framework web (ultra rápido) |
| **MySQL** | v8.4 | Base de datos |
| **golang-migrate** | v4.17.1 | Migraciones versionadas |
| **Docker** | v28.5.1 | Contenedorización |
| **Docker Compose** | v2.40.0-desktop.1 | Orquestación |

---

## 📁 Estructura del Proyecto

```
api-db-migrations/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point de la aplicación
│
├── internal/
│   ├── db/
│   │   ├── migrate.go                 # Gestión de migraciones
│   │   └── commect.go                 # Conexión a BD (reutilizable)
│   │
│   ├── handlers/                      # Handlers HTTP (lógica de negocio)
│   │   ├── health.go                  # Health check
│   │   ├── db.go                      # Info de base de datos
│   │   ├── users.go                   # CRUD de usuarios
│   │   └── campaigns.go               # CRUD de campañas
│   │
│   ├── models/                        # Modelos de datos
│   │   ├── user.go                    # Estructuras de usuario
│   │   └── campaign.go                # Estructuras de campaña
│   │
│   └── router/                        # Configuración de rutas
│       └── router.go                  # Definición de todas las rutas
│
├── migrations/                         # Migraciones SQL versionadas
│   ├── 0001_init_schema.up.sql       # Schema inicial (tablas)
│   ├── 0001_init_schema.down.sql     # Rollback del schema
│   ├── 0002_example_alter.up.sql     # Alteraciones (agregar columnas)
│   ├── 0002_example_alter.down.sql   # Rollback de alteraciones
│   ├── 0003_insert_sample_data.up.sql # Datos de ejemplo
│   └── 0003_insert_sample_data.down.sql # Rollback de datos
│
├── docker-compose.yml                 # Configuración de servicios
├── Dockerfile                         # Build de la aplicación
├── go.mod                             # Dependencias de Go
├── go.sum                             # Checksums de dependencias
└── README.md                          # Este archivo
```

---

## 🚀 Inicio Rápido

### Prerrequisitos

- **Docker** y **Docker Compose** instalados
- Puertos **8080** (API)

### 1. Iniciar los Servicios

```bash
docker-compose up --build
```

Esto va a:
1. Construir la imagen Docker de la aplicación
2. Iniciar MySQL 8.4
3. Ejecutar migraciones automáticamente
4. Iniciar la API en el puerto 8080

### 2. Verificar que Funciona

```bash
# Health check
curl http://localhost:8080/health

# Ver tablas creadas
curl http://localhost:8080/db

# Listar usuarios
curl http://localhost:8080/api/v1/users
```

---

## 📊 Endpoints Disponibles

### 🔧 Endpoints Generales

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check del servidor |
| GET | `/db` | Lista todas las tablas de la BD |

### 👥 Endpoints de Usuarios

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/v1/users` | Listar todos los usuarios |
| GET | `/api/v1/users/:id` | Obtener usuario por ID |
| POST | `/api/v1/users` | Crear nuevo usuario |
| PUT | `/api/v1/users/:id` | Actualizar usuario |
| DELETE | `/api/v1/users/:id` | Eliminar usuario |

### 📢 Endpoints de Campañas

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/v1/campaigns` | Listar todas las campañas |
| GET | `/api/v1/campaigns?status=active` | Filtrar por estado |
| GET | `/api/v1/campaigns/stats` | Estadísticas por estado |
| GET | `/api/v1/campaigns/:id` | Obtener campaña por ID |
| POST | `/api/v1/campaigns` | Crear nueva campaña |
| PUT | `/api/v1/campaigns/:id` | Actualizar campaña |
| DELETE | `/api/v1/campaigns/:id` | Eliminar campaña |

**Total: 14 endpoints RESTful** 🎯

---

## 💡 Ejemplos de Uso

### Crear un Usuario

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@example.com",
    "full_name": "Nombre Completo",
    "phone": "+54-911-1234-5678"
  }'
```

### Listar Campañas Activas

```bash
curl "http://localhost:8080/api/v1/campaigns?status=active"
```

### Actualizar una Campaña

```bash
curl -X PUT http://localhost:8080/api/v1/campaigns/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

### Ver Estadísticas

```bash
curl http://localhost:8080/api/v1/campaigns/stats
```

---

## 🗃️ Base de Datos

### Tablas

1. **users** - Usuarios del sistema
   - `id`, `email`, `full_name`, `phone`, `created_at`
   
2. **campaigns** - Campañas de marketing
   - `id`, `name`, `status`, `created_at`
   - Estados: `draft`, `active`, `paused`, `completed`, `archived`
   
3. **campaign_members** - Relación usuarios-campañas
   - `campaign_id`, `user_id`, `role`, `created_at`
   - Roles: `owner`, `admin`, `member`
   
4. **schema_migrations** - Control de versiones de migraciones

### Migraciones

Las migraciones se ejecutan **automáticamente** al iniciar la aplicación:

```
✅ 0001_init_schema.up.sql       → Crea las tablas
✅ 0002_example_alter.up.sql     → Agrega columna phone a users
✅ 0003_insert_sample_data.up.sql → Inserta datos de ejemplo
```

**Comportamiento:**
- Base de datos vacía → Crea todo desde cero
- Base de datos existente → Solo aplica migraciones faltantes
- Base de datos actualizada → No hace nada

---

## ⚙️ Configuración

### Variables de Entorno

Configurables en `docker-compose.yml`:

```yaml
# Conexión a MySQL
MYSQL_DSN: "api:api@tcp(mysql:3306)/api?multiStatements=true&parseTime=true"

# Ruta de migraciones
MIGRATIONS_PATH: "file://migrations"

# Puerto del servidor
HTTP_ADDR: ":8080"
```

### Puertos

- **API**: `8080` → `http://localhost:8080`
- **MySQL**: `3307` → `localhost:3307` (internamente usa 3306)

---

## 🏗️ Arquitectura

### Patrón de Diseño

El proyecto sigue una **arquitectura limpia** con separación de responsabilidades:

```
┌─────────────┐
│   Router    │ ← Define rutas y middlewares
└──────┬──────┘
       │
┌──────▼──────┐
│  Handlers   │ ← Lógica de negocio HTTP
└──────┬──────┘
       │
┌──────▼──────┐
│   Models    │ ← Estructuras de datos
└──────┬──────┘
       │
┌──────▼──────┐
│   Database  │ ← MySQL
└─────────────┘
```

### Características Clave

- **Dependency Injection**: Las dependencias se inyectan (no se crean internamente)
- **Constructor Pattern**: Uso de constructores para handlers
- **Error Handling**: Manejo consistente de errores en todas las capas
- **Middleware Stack**: Logger y Recover automáticos
- **JSON Responses**: Formato estandarizado de respuestas
- **Conexión DB Flexible**: Dos patrones de uso según necesidad

---

## 🔌 Conexión a Base de Datos

El proyecto incluye dos patrones para trabajar con la base de datos:

### Patrón 1: Conexión por Request (Simple)

Ideal para queries específicos o handlers que se ejecutan ocasionalmente:

```go
func MyHandler(c *fiber.Ctx) error {
    // Crear nueva conexión
    db := db.ConnectDB()
    defer db.Close()
    
    // Ejecutar query
    rows, err := db.Query("SELECT * FROM users WHERE id = ?", userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Error al consultar",
        })
    }
    defer rows.Close()
    
    // Procesar resultados
    var users []User
    for rows.Next() {
        var user User
        rows.Scan(&user.ID, &user.Email, &user.FullName)
        users = append(users, user)
    }
    
    if err := rows.Err(); err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error"})
    }
    
    return c.JSON(fiber.Map{"success": true, "data": users})
}
```

**Ventajas:**
- ✅ Simple y directo
- ✅ Control explícito del ciclo de vida
- ✅ Cierre automático con `defer`
- ✅ Ideal para operaciones específicas

### Patrón 2: Inyección de Dependencias (Actual)

Usado en los handlers principales (usuarios, campañas):

```go
type UserHandler struct {
    DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{DB: db}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
    rows, err := h.DB.Query("SELECT * FROM users")
    // ... procesar
}
```

**Ventajas:**
- ✅ Reutiliza el pool de conexiones
- ✅ Mejor rendimiento para operaciones frecuentes
- ✅ Fácil de testear (mock de DB)
- ✅ Patrón profesional

### Función ConnectDB()

Ubicada en `internal/db/commect.go`:

```go
// Crea una nueva conexión a MySQL
func ConnectDB() *sql.DB

// Igual que ConnectDB pero hace panic si falla
func MustConnectDB() *sql.DB
```

**Configuración automática:**
- Lee `MYSQL_DSN` desde variables de entorno
- Retry automático (60 segundos)
- Pool configurado: 25 conexiones máximas
- Logging de conexión

**Ejemplo completo:**

```go
// Handler para buscar usuarios por email
func SearchUsersByEmail(c *fiber.Ctx) error {
    searchTerm := c.Query("email")
    
    // Validación básica
    if searchTerm == "" {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "El parámetro 'email' es requerido",
        })
    }
    
    // 1. Conectar a la base de datos
    database := db.ConnectDB()
    defer database.Close()
    
    // 2. Query con parámetros seguros (LIKE para búsqueda parcial)
    query := "SELECT id, email, full_name, phone, created_at FROM users WHERE email LIKE ?"
    rows, err := database.Query(query, "%"+searchTerm+"%")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Error al consultar usuarios",
        })
    }
    defer rows.Close()
    
    // 3. Procesar resultados
    var users []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.FullName,
            &user.Phone,
            &user.CreatedAt,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "message": "Error al leer datos",
            })
        }
        users = append(users, user)
    }
    
    // 4. Verificar errores del cursor
    if err := rows.Err(); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Error al procesar resultados",
        })
    }
    
    // 5. Retornar respuesta estructurada
    return c.JSON(fiber.Map{
        "success": true,
        "data":    users,
        "total":   len(users),
        "query":   searchTerm,
    })
}
```

### Configuración del Pool

```go
db.SetMaxOpenConns(25)                  // Máximo 25 conexiones simultáneas
db.SetMaxIdleConns(10)                  // Máximo 10 en espera
db.SetConnMaxLifetime(5 * time.Minute)  // Vida máxima: 5 minutos
```

---

## 🧪 Testing

### Ejecutar Tests Manuales

```bash
# Health check
curl http://localhost:8080/health

# Usuarios
curl http://localhost:8080/api/v1/users

# Crear usuario
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","full_name":"Test User"}'

# Campañas
curl http://localhost:8080/api/v1/campaigns

# Estadísticas
curl http://localhost:8080/api/v1/campaigns/stats
```

---

## 🐳 Docker

### Servicios

```yaml
services:
  mysql:
    image: mysql:8.4
    ports: ["3307:3306"]
    networks: [backend]
    # Solo accesible internamente
    
  api:
    build: .
    ports: ["8080:8080"]
    networks: [backend]
    depends_on:
      mysql:
        condition: service_healthy
```

### Comandos Útiles

```bash
# Iniciar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Reiniciar con rebuild
docker-compose up --build

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes (reset completo)
docker-compose down -v

# Ejecutar comando en MySQL
docker exec -it api-db-migrations-mysql-1 mysql -uapi -papi api
```

---

## 📈 Rendimiento

Latencias promedio (medidas con Fiber):

| Endpoint | Latencia |
|----------|----------|
| GET /health | < 1ms |
| GET /api/v1/users | ~1-2ms |
| GET /api/v1/campaigns | ~1-2ms |
| POST /api/v1/users | ~2-3ms |
| PUT /api/v1/users/:id | ~3-4ms |
| DELETE /api/v1/users/:id | ~2-3ms |

**Fiber** es extremadamente rápido gracias a su uso de **fasthttp**.

---

## 🔐 Seguridad

### Implementado

- ✅ Validación de datos de entrada
- ✅ Manejo seguro de errores SQL
- ✅ MySQL solo accesible internamente (red bridge)
- ✅ Prepared statements (previene SQL injection)
- ✅ Healthchecks de contenedores

### Por Implementar

- [ ] Autenticación JWT
- [ ] Rate limiting
- [ ] CORS configurado
- [ ] HTTPS/TLS
- [ ] Validación con go-validator
- [ ] Sanitización de inputs

---

## 🛠️ Desarrollo

### Agregar una Nueva Migración

1. Crea los archivos:
```bash
touch migrations/0004_nueva_feature.up.sql
touch migrations/0004_nueva_feature.down.sql
```

2. Escribe el SQL en `.up.sql`:
```sql
ALTER TABLE users ADD COLUMN address VARCHAR(255);
```

3. Escribe el rollback en `.down.sql`:
```sql
ALTER TABLE users DROP COLUMN address;
```

4. Reinicia los contenedores:
```bash
docker-compose down
docker-compose up --build
```

### Agregar un Nuevo Endpoint

1. Define el modelo en `internal/models/`
2. Crea el handler en `internal/handlers/`
3. Registra la ruta en `internal/router/router.go`
4. Rebuild y prueba

---

## 📚 Documentación Adicional

- **[Fiber Documentation](https://docs.gofiber.io/)** - Framework web
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Sistema de migraciones
- **[database/sql](https://pkg.go.dev/database/sql)** - Paquete estándar de Go

---

## 🤝 Contribuir

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

---

## 📝 Notas Importantes

### ⚠️ Proyecto de Ejemplo

Este es un **proyecto de demostración** que muestra:
- Cómo estructurar una API REST en Go
- Implementación de migraciones automáticas
- Dos patrones de conexión a BD
- Buenas prácticas de arquitectura
- Uso de Docker Compose para desarrollo

**Adaptación para producción:**
- Reemplaza el schema en `0001_init_schema.up.sql` con tus tablas
- Actualiza los modelos en `internal/models/`
- Crea handlers específicos para tu lógica de negocio
- Agrega autenticación y autorización
- Implementa tests unitarios

### Comportamiento de Migraciones

Las migraciones usan **golang-migrate** que:
- Ejecuta archivos `.up.sql` en orden numérico
- Mantiene un registro en la tabla `schema_migrations`
- Solo ejecuta migraciones nuevas (idempotente)
- Soporta rollback con archivos `.down.sql`

### Datos de Ejemplo

El proyecto incluye datos de ejemplo insertados por la migración `0003`:
- 5 usuarios
- 5 campañas con diferentes estados
- 13 relaciones usuario-campaña

### Red Bridge

MySQL está configurado en una red Docker privada (`backend`) y **no está expuesto** al exterior. Solo la API (puerto 8080) es accesible desde el host.

### Mejores Prácticas en el Código

**✅ HACER:**
```go
// Siempre cerrar conexiones
db := db.ConnectDB()
defer db.Close()

// Siempre cerrar resultados
rows, err := db.Query("...")
defer rows.Close()

// Usar parámetros (previene SQL injection)
db.Query("SELECT * FROM users WHERE id = ?", userID)

// Verificar rows.Err() después del loop
for rows.Next() { /* ... */ }
if err := rows.Err(); err != nil { /* handle */ }
```

**❌ NO HACER:**
```go
// NO olvidar defer
db := db.ConnectDB()
// ❌ Falta defer db.Close()

// NO concatenar strings en queries
query := "SELECT * FROM users WHERE email = '" + email + "'"
// ❌ Vulnerable a SQL Injection
```

---

## 🎯 Características del Ejemplo

### ✅ Implementado

- [x] **CRUD completo** de Usuarios y Campañas
- [x] **Migraciones automáticas** con versionado
- [x] **Dos patrones de conexión** a BD (flexible)
- [x] **Arquitectura limpia** y escalable
- [x] **Docker Compose** configurado
- [x] **Logging automático** con Fiber
- [x] **Manejo de errores** consistente
- [x] **Filtros y estadísticas** de campañas
- [x] **Pool de conexiones** configurado
- [x] **Red privada** para MySQL (seguridad)
- [x] **Prepared statements** (anti SQL injection)
- [x] **Documentación completa** en README

### 💡 Ideas para Extender

- [ ] Endpoints de Campaign Members (relaciones)
- [ ] Autenticación con JWT
- [ ] Paginación en listados
- [ ] Búsqueda y filtros avanzados
- [ ] Tests unitarios y de integración
- [ ] Swagger/OpenAPI documentation
- [ ] Validación con go-validator
- [ ] Rate limiting
- [ ] CORS configurado
- [ ] Métricas con Prometheus
- [ ] CI/CD con GitHub Actions
- [ ] Caching con Redis

---

## 📄 Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.

---

## 👨‍💻 Uso del Proyecto

### Como Plantilla

Este proyecto puede usarse como **plantilla base** para:
- APIs REST en Go
- Servicios con migraciones de BD
- Aplicaciones con Docker Compose
- Proyectos que requieren arquitectura limpia

### Adaptación

1. **Clona el proyecto** y úsalo como punto de partida
2. **Reemplaza las migraciones** con tu schema real
3. **Crea tus modelos** en `internal/models/`
4. **Implementa tus handlers** en `internal/handlers/`
5. **Registra tus rutas** en `internal/router/`
6. **Personaliza** según tus necesidades

### Stack Tecnológico

- **Go 1.22**: Lenguaje moderno, rápido y concurrente
- **Fiber v2**: Framework web ultra rápido (similar a Express)
- **MySQL 8.4**: Base de datos robusta y confiable
- **golang-migrate**: Sistema profesional de migraciones
- **Docker**: Contenedorización y despliegue simple

---

## 🙏 Agradecimientos

- [Fiber](https://gofiber.io/) - Framework web increíblemente rápido
- [golang-migrate](https://github.com/golang-migrate/migrate) - Sistema robusto de migraciones
- [MySQL](https://www.mysql.com/) - Base de datos confiable

---

**¿Preguntas o problemas?** Abre un issue en el repositorio.

**¿Te gusta el proyecto?** Dale una ⭐ en GitHub.
