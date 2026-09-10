# Mesa de ayuda

> Reemplacen todo lo que está entre corchetes en su primer commit.
> Este README es la puerta de entrada del repositorio: en la semana 1 otra pareja debe poder levantar el servidor siguiendo solo lo que dice aquí, y desde la semana 4 es la base de la integración continua.

**Hilo Servidor · ULEAM · Período 2026-2**
Aplicaciones Web II (TDI-610) · Aplicación para el Servidor Web (IS-503)

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Gonzalez Reina Isaac Mateo | ElRascuacho | Servidor Web A |
| Velez Briones Jipson Jordan | jipsonvb27 | Servidor Web A |

## El producto

Dominio de trabajo de la semana 1–2: mesa de ayuda. La ficha del negocio propio se completa cuando pase la compuerta de la semana 3.

## Cómo levantar el servidor

Requisitos: Go 1.25+ y PostgreSQL (contenedor Docker, ver abajo).

```bash
go run .
```

El servidor responde en `http://localhost:8080`.

- `/` → `{"estado":"vivo"}`
- `/salud` → `{"bd":"ok","version":"..."}` (requiere la base levantada)

## Base de datos

- Opción usada por la pareja: contenedor
- Con contenedor (puerto **5433** porque el 5432 local ya estaba ocupado):

```bash
docker run --name pg -e POSTGRES_PASSWORD=taller2026 -p 5433:5432 -d postgres:16
docker exec pg psql -U postgres -c "create database mesa_ayuda;"
```

- Base de datos del proyecto: `mesa_ayuda`
- Usuario: `postgres` · Contraseña: `taller2026` · Puerto host: `5433`

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
docs/ficha_negocio.md   → ficha del negocio (se entrega antes del día A de la semana 3)
main.go                 → punto de entrada del servidor
```
