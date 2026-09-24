# RentCar

**Hilo Servidor · ULEAM · Período 2026-2**
Aplicaciones Web II (TDI-610) · Aplicación para el Servidor Web (IS-503)

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Gonzalez Reina Isaac Mateo | ElRascuacho | Servidor Web A |
| Velez Briones Jipson Jordan | jipsonvb27 | Servidor Web A |

## El producto

**RentCar** — API de alquiler de vehículos. Solo servidor (Go + Chi + GORM + PostgreSQL).

Entidades: `Vehiculo` (1) — `Reserva` (N, con estados). El CRUD calificado es el de reservas.

## Cómo levantar el servidor

Requisitos: Go 1.25+ y PostgreSQL (contenedor Docker, ver abajo).

Si `go` no se reconoce en la terminal, el SDK está en `C:\Users\MSI BRAVO\go-sdk\bin` (ya se agregó al PATH de usuario: cierra y vuelve a abrir la terminal).

1. Copie `.env.example` a `.env` y complete la contraseña de PostgreSQL.
2. Arranque:

```bash
go run . -reset   # primera vez: crea tablas y semilla
go run .          # arranques siguientes
```

Guía completa de peticiones para Bruno/demo: [`docs/pruebas.md`](docs/pruebas.md).

El servidor responde en el puerto de `PUERTO` (por defecto `8080`).

Rutas principales:

- `POST/GET /reservas` · `GET/PUT/DELETE /reservas/{id}` · filtro `?estado=`
- `GET /vehiculos` — listado con `Preload` de reservas

## Base de datos

- Opción usada por la pareja: contenedor
- Con contenedor (puerto **5433** porque el 5432 local ya estaba ocupado):

```bash
docker run --name pg -e POSTGRES_PASSWORD=CAMBIE_ESTO -p 5433:5432 -d postgres:16
docker exec pg psql -U postgres -c "CREATE DATABASE rentcar;"
```

- La cadena de conexión (host, usuario, contraseña, `dbname`, puerto) va en `.env` — ver `.env.example`. No se sube al repositorio.

## Cómo correr las pruebas

```bash
go test ./...
```

Las pruebas corren también en la integración continua (pestaña Actions). Desde la semana 4, un entregable cuyas pruebas no pasan en la integración no se recibe.

## Convenciones del repositorio

- Un commit de cada integrante como mínimo por taller; el commit de cierre se hace en clase.
- Mensajes de commit: qué cambió y por qué, entendibles sin el autor presente.
- Uso de IA declarado en el cuerpo del commit: una línea con qué herramienta y para qué parte.
- Ningún secreto en el código ni en el historial: la configuración se externaliza (semana 4).

## Estructura

```
main.go                      → conexión, AutoMigrate, -reset, rutas
internal/reservas/           → modelos, manejadores, semilla
internal/respuesta/          → envoltura {"ok": …}
internal/middleware/         → registro y recuperación
docs/decisiones.md           → decisión D3 del taller S3
```
