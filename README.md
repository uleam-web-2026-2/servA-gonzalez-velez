# servA-gonzalez-reina

Pareja del hilo Servidor A · ULEAM · Período 2026-2

## Integrantes

- Gonzalez (GitHub: ElRascuacho)
- Reina (completar usuario GitHub del segundo integrante)

## PostgreSQL

**Opción B — Contenedor Docker** (PostgreSQL 16).

El puerto `5432` ya estaba ocupado por una instalación nativa, así que el contenedor usa **5433** (mapeo `5433:5432`), como indica el taller.

```bash
docker run --name pg -e POSTGRES_PASSWORD=taller2026 -p 5433:5432 -d postgres:16
```

- Host: `localhost`
- Puerto: `5433`
- Usuario: `postgres`
- Contraseña: `taller2026`
- Base: `mesa_ayuda`

Crear la base (si el contenedor es nuevo):

```bash
docker exec pg psql -U postgres -c "create database mesa_ayuda;"
```

Verificar:

```bash
docker exec pg psql -U postgres -d mesa_ayuda -c "select version();"
```

## Cómo levantar el servidor

Requisitos: Go 1.25+ y el contenedor `pg` en marcha.

```bash
go run .
```

Probar:

- http://localhost:8080/ → `{"estado":"vivo"}`
- http://localhost:8080/salud → `{"bd":"ok","version":"..."}`

Para detener el servidor: `Ctrl+C`.

## Notas

La contraseña de laboratorio está en `main.go` solo para la prueba de montaje (semana 1). No usar contraseñas personales en commits.
