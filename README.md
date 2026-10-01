# RentCar

**Hilo Servidor · ULEAM · Período 2026-2**
Aplicaciones Web II (TDI-610) · Aplicación para el Servidor Web (IS-503)

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Gonzalez Reina Isaac Mateo | ElRascuacho | Servidor Web A |
| Velez Briones Jipson Jordan | jipsonvb27 | Servidor Web A |

## El producto

**RentCar** — API de alquiler de **autos y camionetas**. Solo servidor (Go + Chi + GORM + PostgreSQL).

Entidades: `Sucursal`, `Usuario`, `Cliente`, `Vehiculo`, `Reserva` (con estados), `Pago`. El CRUD calificado es el de reservas.

## Cómo levantar el servidor

Requisitos: Go 1.25+ y PostgreSQL (contenedor Docker, ver abajo).

Si `go` no se reconoce en la terminal, el SDK está en `C:\Users\MSI BRAVO\go-sdk\bin` (ya se agregó al PATH de usuario: cierra y vuelve a abrir la terminal).

1. Copie `.env.example` a `.env` y complete la contraseña de PostgreSQL.
2. Arranque:

```bash
go run . -reset   # primera vez: crea tablas y semilla
go run .          # arranques siguientes
```

Guía completa de peticiones: [`docs/pruebas.md`](docs/pruebas.md).

Hito 1: [`DOC/hito1_ficha_del_negocio.md`](DOC/hito1_ficha_del_negocio.md) · [`DOC/hito1_addendum.md`](DOC/hito1_addendum.md).

El servidor responde en el puerto de `PUERTO` (por defecto `8080`).

Rutas principales:

- `POST/GET /reservas` · `GET/PUT/PATCH/DELETE /reservas/{id}` · filtros `?estado=` · `?limit=` · `?offset=`
- `GET /vehiculos` — listado con `Preload` de reservas
- `GET /sucursales` · `GET /clientes` · `GET /usuarios` · `GET /pagos`

## Base de datos

- Opción usada por la pareja: contenedor
- Con contenedor (puerto **5433** porque el 5432 local ya estaba ocupado):

```bash
docker run --name pg -e POSTGRES_PASSWORD=CAMBIE_ESTO -p 5433:5432 -d postgres:16
docker exec pg psql -U postgres -c "CREATE DATABASE rentcar;"
```

- La cadena de conexión va en `.env` — ver `.env.example`. No se sube al repositorio.

## Cómo correr las pruebas

```bash
go test ./...
```

Las pruebas corren también en la integración continua (pestaña Actions).

## Convenciones del repositorio

- Un commit de cada integrante como mínimo por taller; el commit de cierre se hace en clase.
- Mensajes de commit: qué cambió y por qué, entendibles sin el autor presente.
- Uso de IA declarado en el cuerpo del commit.
- Ningún secreto en el código ni en el historial.

## Estructura

```
main.go                      → conexión, AutoMigrate, -reset, rutas
internal/reservas/           → modelos (6 entidades), manejadores, semilla
internal/config/             → lectura de .env
internal/respuesta/          → envoltura {"ok": …}
internal/middleware/         → registro y recuperación
docs/                        → decisiones, pruebas
DOC/                         → Hito 1 (ficha y addendum)
```
