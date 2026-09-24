# RentCar — Guía de pruebas (taller S3)

Base URL: `http://localhost:8080`  
Herramienta recomendada: **Bruno** (o curl / Postman).

---

## 0. Antes de probar

### Requisitos

1. **Go** en el PATH (si `go` no se reconoce, cierra y vuelve a abrir la terminal; el binario está en `C:\Users\MSI BRAVO\go-sdk\bin`).
2. **PostgreSQL** en Docker (puerto host `5433`):

```bash
docker start pg
# si el contenedor no existe:
docker run --name pg -e POSTGRES_PASSWORD=CAMBIE_ESTO -p 5433:5432 -d postgres:16
docker exec pg psql -U postgres -c "CREATE DATABASE rentcar;"
```

3. Copiar `.env.example` a `.env` y completar la contraseña real de PostgreSQL.

4. Arrancar el servidor desde la carpeta del repo:

```bash
go run . -reset
```

En la terminal deben verse dos `CREATE TABLE` y el segundo con `CONSTRAINT "fk_vehiculos_reservas" FOREIGN KEY`.  
Mensaje final: `escuchando en el puerto 8080`.

Para arranques siguientes (sin borrar datos):

```bash
go run .
```

---

## 1. Resumen de rutas

| Método | Ruta | Qué hace | Código OK |
| --- | --- | --- | --- |
| `POST` | `/reservas` | Crear reserva | `201` |
| `GET` | `/reservas` | Listar todas | `200` |
| `GET` | `/reservas?estado=pendiente` | Filtrar por estado | `200` |
| `GET` | `/reservas/{id}` | Ver una | `200` |
| `PUT` | `/reservas/{id}` | Actualizar | `200` |
| `DELETE` | `/reservas/{id}` | Borrar | `200` |
| `GET` | `/vehiculos` | Listar vehículos **con** sus reservas (`Preload`) | `200` |

**Estados válidos:** `pendiente` · `confirmada` · `en_curso` · `finalizada` · `cancelada`

**Importante:** en el JSON las claves van con mayúscula exacta del struct (`VehiculoID`, `Cliente`, `FechaInicio`, …). Si escribes `vehiculo_id`, Go lo ignora y el ID llega en 0.

---

## 2. Punto 1 — Lo que mostrar en la primera pasada

### 2.1 Arranque con tablas

```bash
go run . -reset
```

Mostrar en la terminal los dos `CREATE TABLE` y la FK.

### 2.2 Crear una reserva — `POST /reservas`

**Bruno:** Body → JSON

```json
{
  "VehiculoID": 1,
  "Cliente": "Prueba Taller",
  "FechaInicio": "2026-11-01",
  "FechaFin": "2026-11-04",
  "Total": 105,
  "Estado": "pendiente"
}
```

**curl (Git Bash):**

```bash
curl -i -X POST http://localhost:8080/reservas \
  -H "Content-Type: application/json" \
  -d '{"VehiculoID":1,"Cliente":"Prueba Taller","FechaInicio":"2026-11-01","FechaFin":"2026-11-04","Total":105,"Estado":"pendiente"}'
```

**Esperado:** `201` y envoltura `{"ok":true,"datos":{...}}` con `ID` asignado.

### 2.3 Listar — `GET /reservas`

```bash
curl -i http://localhost:8080/reservas
```

**Esperado:** `200` y la lista incluyendo la reserva recién creada.

---

## 3. Punto 2 — Lo que mostrar en la segunda pasada

### 3.1 Ver uno — `GET /reservas/{id}`

```bash
curl -i http://localhost:8080/reservas/1
```

**ID inexistente (404):**

```bash
curl -i http://localhost:8080/reservas/999
```

### 3.2 Actualizar — `PUT /reservas/{id}`

```json
{
  "Cliente": "Prueba Taller",
  "FechaInicio": "2026-11-01",
  "FechaFin": "2026-11-05",
  "Total": 140,
  "Estado": "confirmada"
}
```

```bash
curl -i -X PUT http://localhost:8080/reservas/1 \
  -H "Content-Type: application/json" \
  -d '{"Cliente":"Prueba Taller","FechaInicio":"2026-11-01","FechaFin":"2026-11-05","Total":140,"Estado":"confirmada"}'
```

**Nota:** el servidor **no cambia** el `VehiculoID` (decisión D3).

**ID inexistente → 404.**

### 3.3 Borrar — `DELETE /reservas/{id}`

```bash
curl -i -X DELETE http://localhost:8080/reservas/4
```

**Esperado:** `200` con mensaje de eliminada. Si el id no existe → `404`.

### 3.4 Validaciones

| Caso | Petición | Código |
| --- | --- | --- |
| JSON roto | Body `{roto` | `400` `json_invalido` |
| Estado inventado | `"Estado":"inventado"` | `422` `estado_invalido` |
| Vehículo que no existe | `"VehiculoID":99` | `422` `vehiculo_inexistente` |
| Cliente vacío | `"Cliente":""` | `422` `datos_invalidos` |
| Fecha fin ≤ inicio | fechas al revés | `422` `datos_invalidos` |

**Estado inventado:**

```json
{
  "VehiculoID": 1,
  "Cliente": "X",
  "FechaInicio": "2026-11-01",
  "FechaFin": "2026-11-04",
  "Total": 10,
  "Estado": "inventado"
}
```

**JSON roto:**

```bash
curl -i -X POST http://localhost:8080/reservas \
  -H "Content-Type: application/json" \
  -d '{roto'
```

### 3.5 Preload (sin N+1) — `GET /vehiculos`

```bash
curl -i http://localhost:8080/vehiculos
```

En la **terminal del servidor** (después del arranque) deben verse **exactamente 2 consultas** SQL: una a `vehiculos` y otra a `reservas`. Cada vehículo trae su lista `Reservas` en el JSON.

### 3.6 Filtro seguro — `GET /reservas?estado=`

Solo pendientes:

```bash
curl -i "http://localhost:8080/reservas?estado=pendiente"
```

Intento de inyección (no debe devolver todos):

```bash
curl -i "http://localhost:8080/reservas?estado=x'%20OR%20'1'='1"
```

**Esperado:** lista vacía `[]` (o sin esos registros), **no** todas las reservas.

### 3.7 Decisión

Abrir `docs/decisiones.md` → sección **D3** (qué campos se pueden cambiar al actualizar) y explicarla sin leer el código.

---

## 4. Datos de la semilla (después de `-reset`)

| ID vehículo | Marca / modelo | Tipo | Precio/día |
| --- | --- | --- | --- |
| 1 | Toyota Corolla | auto | 35 |
| 2 | Yamaha MT-07 | moto | 25 |
| 3 | Chevrolet Colorado | camioneta | 55 |

Reservas de ejemplo: Ana Pérez (`confirmada`), Luis Mora (`pendiente`), Carla Ruiz (`en_curso`).

---

## 5. Orden sugerido para la demo en clase

1. `go run . -reset` → mostrar `CREATE TABLE` + FK.  
2. `POST /reservas` → `201`.  
3. `GET /reservas` → aparece lo creado.  
4. `GET /reservas/999` → `404`.  
5. `PUT` + `DELETE` de una reserva.  
6. `POST` con estado inventado → `422`; JSON roto → `400`.  
7. `GET /vehiculos` → mirar las **2 consultas** en la terminal.  
8. `GET /reservas?estado=pendiente` y luego el intento de inyección.  
9. Abrir `docs/decisiones.md` y explicar D3.

---

## 6. Si algo falla

| Síntoma | Qué revisar |
| --- | --- |
| `go` no se reconoce | Cerrar terminal / Cursor y abrir de nuevo. Comprobar: `go version`. |
| `connection refused` | `docker start pg` |
| `database "rentcar" does not exist` | `docker exec pg psql -U postgres -c "CREATE DATABASE rentcar;"` |
| `address already in use` | Otro proceso en `:8080`; ciérralo. |
| Campo en 0 / vacío al crear | Claves JSON con mayúsculas (`VehiculoID`, no `vehiculo_id`). |
