# RentCar — Guía de pruebas (taller S3 + Hito 1)

Base URL: `http://localhost:8080`  
Herramienta recomendada: **Bruno** (o curl / Postman).

---

## 0. Antes de probar

1. **Go** en el PATH.
2. **PostgreSQL** en Docker (puerto host `5433`):

```bash
docker start pg
docker exec pg psql -U postgres -c "CREATE DATABASE rentcar;"
```

3. Copiar `.env.example` a `.env` y completar la contraseña.
4. Arrancar:

```bash
go run . -reset
```

---

## 1. Resumen de rutas

| Método | Ruta | Qué hace | Código OK |
| --- | --- | --- | --- |
| `POST` | `/reservas` | Crear reserva | `201` |
| `GET` | `/reservas` | Listar (`estado`, `limit`, `offset`) | `200` |
| `GET` | `/reservas/{id}` | Ver una | `200` |
| `PUT` | `/reservas/{id}` | Actualizar (sin cambiar VehiculoID) | `200` |
| `PATCH` | `/reservas/{id}` | Cambiar solo el estado | `200` |
| `DELETE` | `/reservas/{id}` | Borrar | `200` |
| `GET` | `/vehiculos` | Listar con reservas (`Preload`) | `200` |
| `GET` | `/sucursales` | Listar sucursales | `200` |
| `GET` | `/clientes` | Listar clientes | `200` |
| `GET` | `/usuarios` | Listar usuarios | `200` |
| `GET` | `/pagos` | Listar pagos | `200` |

**Estados de reserva:** `pendiente` · `confirmada` · `en_curso` · `finalizada` · `cancelada`

**Transiciones:** `pendiente`→`confirmada`|`cancelada` · `confirmada`→`en_curso`|`cancelada` · `en_curso`→`finalizada`. Terminales no vuelven atrás.

**Tipos de vehículo:** solo `auto` y `camioneta`.

**JSON:** claves con mayúscula del struct (`ClienteID`, `VehiculoID`, `FechaInicio`, `Estado`).

---

## 2. Crear reserva — `POST /reservas`

```json
{
  "ClienteID": 1,
  "VehiculoID": 1,
  "FechaInicio": "2026-11-01",
  "FechaFin": "2026-11-04",
  "Total": 105,
  "Estado": "pendiente"
}
```

**Esperado:** `201`.

## 3. Cambiar estado — `PATCH /reservas/{id}`

```json
{ "Estado": "confirmada" }
```

Transición ilegal (ej. `finalizada` → `pendiente`): **422** `transicion_prohibida`.

## 4. Validaciones

| Caso | Código |
| --- | --- |
| JSON roto | `400` `json_invalido` |
| Estado inventado | `422` `estado_invalido` |
| Cliente inexistente | `422` `cliente_inexistente` |
| Vehículo inexistente | `422` `vehiculo_inexistente` |
| Fecha fin ≤ inicio | `422` `datos_invalidos` |
| `limit=abc` | `400` `paginacion_invalida` |
| Transición prohibida | `422` `transicion_prohibida` |

## 5. Semilla (`-reset`)

| Entidad | Ejemplo |
| --- | --- |
| Sucursal | RentCar Manta Centro |
| Usuarios | administrador + agente |
| Clientes | Ana Pérez, Luis Mora, Carla Ruiz |
| Vehículos | Toyota Corolla (auto), Kia Rio (auto), Chevrolet Colorado (camioneta) |
| Reservas | confirmada, pendiente, en_curso |
| Pagos | verificado + por_verificar |

## 6. Pruebas automáticas

```bash
go test ./...
```
