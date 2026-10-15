# Resumen Detallado · Fase 0 y Fase 1 (Taller Integrador Semana 7)

Este documento explica con total detalle qué se hizo, en qué archivos y por qué motivo en los dos primeros commits locales del taller integrador para el proyecto **RentCar**.

---

## 📌 Commit 1: Fase 0 · El Contrato de la API

- **Hash del commit:** `5c98fef`
- **Mensaje:** `S7: contrato de la API`
- **Declaración de IA:** `IA: Antigravity como guía técnica para el contrato de la API`

### ¿Cuál era el problema antes de este commit?
El servidor ya tenía endpoints funcionales, pero no existía una especificación formal estandarizada (OpenAPI) que definiera con rigor los esquemas de datos, tipos de parámetros, códigos de éxito y códigos de error que la API promete responder a clientes externos (frontend o Bruno).

### Archivos creados y qué hace cada uno:

#### 1. `redocly.yaml` (en la raíz del repositorio)
- **Ruta:** [redocly.yaml](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/redocly.yaml)
- **Qué contiene:** La configuración de reglas para el linter de documentación `@redocly/cli`.
- **Por qué está así:** Hereda de las reglas recomendadas (`recommended`) y mantiene temporalmente apagada la regla `security-defined: off` (se encenderá en la Fase 2 cuando se configure la seguridad JWT global).

#### 2. `docs/openapi.yaml`
- **Ruta:** [docs/openapi.yaml](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/docs/openapi.yaml)
- **Qué contiene:** El contrato formal en estándar OpenAPI 3.1 adaptado al negocio de RentCar:
  - **Servidor:** `http://localhost:8080`.
  - **Operaciones:**
    - `GET /reservas`: Listar reservas (`operationId: listarReservas`). Devuelve `200` con arreglo de `Reserva` o `400`.
    - `POST /reservas`: Crear reserva (`operationId: crearReserva`). Recibe `ReservaNueva` y responde `201` con la reserva creada o `400` / `422`.
    - `GET /reservas/{id}`: Detalle de una reserva (`operationId: verReserva`). Responde `200` o `400` / `404`.
    - `PATCH /reservas/{id}`: Cambio de estado de la reserva (`operationId: cambiarEstadoReserva`). Recibe `{ "estado": "..." }` y responde `200` con la reserva actualizada, o `400` (JSON/ID inválido), `404` (no encontrada) o `422` (transición prohibida o estado no permitido).
    - `PUT /reservas/{id}`: Modificación de datos generales (`operationId: modificarReserva`), no del estado.
    - `DELETE /reservas/{id}`: Eliminación técnica de reserva (`operationId: borrarReserva`).
    - `POST /login`: Inicio de sesión (`operationId: iniciarSesion`). Recibe `usuario` y `clave`, responde `200` con `token` o `401`.
  - **Esquemas (`components.schemas`):**
    - `EstadoReserva`: Enum con los 5 estados del negocio: `[pendiente, confirmada, en_curso, finalizada, cancelada]`.
    - `Reserva`: Mapeo exacto de los campos del struct de Go (`ID`, `ClienteID`, `VehiculoID`, `FechaInicio`, `FechaFin`, `Total`, `Estado`).
    - `ReservaNueva`: Campos obligatorios para crear una reserva sin el `ID`.
    - `Error`: Objeto estándar de error con la clave `error` de tipo string (`{ "error": "mensaje" }`).

#### 3. `docs/api.html`
- **Ruta:** [docs/api.html](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/docs/api.html)
- **Qué es:** La documentación web interactiva y visual autogenerada a partir del archivo YAML mediante el comando:
  ```powershell
  npx @redocly/cli build-docs docs/openapi.yaml -o docs/api.html
  ```
- **Resultado de la validación previa:**
  ```text
  validating docs\openapi.yaml...
  docs\openapi.yaml: validated in 82ms
  Woohoo! Your API description is valid. 🎉
  ```

---

## 📌 Commit 2: Fase 1 · El Cambio de Estado con 422 (40%)

- **Hash del commit:** `52de446`
- **Mensaje:** `S7: cambio de estado con 422`
- **Declaración de IA:** `IA: Antigravity como guía técnica para el cambio de estado con 422`

### ¿Cuál era el problema antes de este commit?
1. Cualquier cliente podía enviar un estado inventado o saltarse el flujo lógico del alquiler de autos (por ejemplo, pasar una reserva recién nacida a `finalizada` sin verificar el pago ni retirar el auto).
2. El método `PUT` permitía sobreescribir el estado directamente en el cuerpo JSON, violando la regla de que el estado **solo** debe modificarse mediante el endpoint dedicado `PATCH`.
3. El método `POST` permitía crear una reserva directamente en estados finales como `finalizada` o `cancelada`.

### Archivos modificados y qué se programó en cada uno:

#### 1. `internal/reservas/estado.go` (Archivo Nuevo)
- **Ruta:** [internal/reservas/estado.go](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/internal/reservas/estado.go)
- **Qué contiene:**
  - **La máquina de estados en memoria (`siguientes`):**
    ```go
    var siguientes = map[string][]string{
        "pendiente":  {"confirmada", "cancelada"},
        "confirmada": {"en_curso", "cancelada"},
        "en_curso":   {"finalizada"},
        "finalizada": {},
        "cancelada":  {},
    }
    ```
    - Todos los estados van como clave.
    - Los estados finales (`finalizada`, `cancelada`) tienen una lista vacía `{}` porque desde ellos ya no se puede transicionar a ningún otro estado.
  - **Función `puedePasar(de, a string) bool`:** Verifica con `slices.Contains(siguientes[de], a)` si el negocio permite el salto de un estado a otro.
  - **Manejador `(m *Manejador) cambiarEstado(w, r)`:**
    1. Extrae y valida el parámetro `id` de la URL (`chi.URLParam`). Si no es un entero positivo $\rightarrow$ responde `400 Bad Request` (`"el id debe ser un número entero positivo"`).
    2. Decodifica el cuerpo JSON en `var pedido struct { Estado string }`. Si el JSON está roto $\rightarrow$ responde `400 Bad Request` (`"el cuerpo no es JSON válido"`).
    3. Verifica si el estado pedido existe en el mapa `siguientes`. Si no existe $\rightarrow$ responde `422 Unprocessable Entity` (`"estado no permitido: " + pedido.Estado`).
    4. Consulta la reserva en la base de datos con `m.DB.First(&t, id)`. Si no existe $\rightarrow$ responde `404 Not Found` (`"no existe la reserva"`).
    5. **Aplica la regla de negocio:** Compara el estado actual `t.Estado` con el pedido `pedido.Estado`. Si `!puedePasar(t.Estado, pedido.Estado)` $\rightarrow$ responde `422 Unprocessable Entity` (`"de "+t.Estado+" no se puede pasar a "+pedido.Estado`).
    6. Asigna el nuevo estado, guarda en la base con `m.DB.Save(&t)` y responde `200 OK` con la reserva completa serializada en JSON mediante `respuesta.JSON(w, http.StatusOK, t)`.

#### 2. `internal/reservas/manejadores.go`
- **Ruta:** [internal/reservas/manejadores.go](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/internal/reservas/manejadores.go)
- **Cambios realizados:**
  - **Eliminación de colisión:** Se removió la función antigua `cambiarEstado` que estaba duplicada en este archivo para que rija exclusivamente la de `estado.go`.
  - **Cierre de puerta en `actualizar` (PUT):**
    Se eliminó la línea que reasignaba el estado:
    ```go
    // Antes: existente.Estado = entrada.Estado
    // Ahora:
    existente.ClienteID = entrada.ClienteID
    existente.FechaInicio = entrada.FechaInicio
    existente.FechaFin = entrada.FechaFin
    existente.Total = entrada.Total
    // el estado solo cambia con PATCH
    ```
    Si un cliente envía un estado distinto en un `PUT`, el servidor responde `200 OK` pero el estado de la reserva **no se modifica**.
  - **Cierre de puerta en `crear` (POST):**
    Se agregó una validación explícita para que toda reserva nueva empiece obligatoriamente en `pendiente`:
    ```go
    if reserva.Estado != "pendiente" {
        respuesta.Error(w, http.StatusUnprocessableEntity,
            "estado_invalido", "una reserva nueva nace en pendiente")
        return
    }
    ```

#### 3. `internal/respuesta/respuesta.go`
- **Ruta:** [internal/respuesta/respuesta.go](file:///C:/Users/Camila/Documents/Mateo/servA-gonzalez-velez/internal/respuesta/respuesta.go)
- **Cambios realizados:**
  - Se agregó la función `respuesta.JSON(w, estado, datos)` para enviar directamente estructuras Go serializadas como JSON sin envolturas adicionales, tal como lo exige el contrato del curso.
  - Se flexibilizó `respuesta.Error(w, estado, args ...string)` para que:
    - Si se llama con un solo mensaje (ej. desde `estado.go` o `guardia.go`), emita directamente `{"error": mensaje}` conforme a OpenAPI.
    - Si se llama con dos parámetros (`codigo`, `mensaje`), mantenga compatibilidad con las llamadas previas del proyecto.

---

## 🧪 Verificación Realizada y Estado de Pruebas

Se ejecutaron las pruebas automáticas del proyecto:
```powershell
go test ./internal/... -v
```

### Resultados obtenidos:
- `internal/config`: `PASS` (`TestCargarSinDatabaseURLFalla`)
- `internal/reservas`: `PASS`
  - `TestCrear/JSON_roto_responde_400`: `PASS`
  - `TestCrear/estado_inventado_responde_422`: `PASS`
  - `TestCrear/cliente_id_cero_responde_422`: `PASS`
  - `TestCrear/fecha_fin_no_posterior_responde_422`: `PASS`
  - `TestVerUnoConIDQueNoEsNumeroResponde400`: `PASS`
  - `TestListarConLimitNoNumericoResponde400`: `PASS`
  - `TestCambiarEstadoJSONRotoResponde400`: `PASS`
  - `TestCambiarEstadoInventadoResponde422`: `PASS`
  - `TestTransicionPermitida`: `PASS`

---

## 🎓 Preguntas Clave para la Defensa del Punto 1

1. **¿Cuál es la transición prohibida de su negocio y por qué?**
   - *Respuesta:* De `pendiente` a `finalizada` (o de `pendiente` a `en_curso`). Una reserva en estado `pendiente` aún no tiene el pago verificado ni ha pasado por la entrega física de la unidad en la sucursal. Pasar directamente a `finalizada` violaría los controles contables y de auditoría de la flota vehicular.
2. **¿Por qué responde 422 y no 400 ni 403?**
   - *Respuesta:* No es `400 Bad Request` porque la sintaxis HTTP y el JSON recibido están perfectamente formados; no es `403 Forbidden` porque el usuario tiene permisos; es `422 Unprocessable Entity` porque la instrucción tiene sentido semántico pero infringe una regla de negocio del dominio.
3. **¿En qué archivo y función se valida la máquina de estados?**
   - *Respuesta:* En `internal/reservas/estado.go`, dentro de la función `cambiarEstado`, mediante la condición:
     ```go
     if !puedePasar(t.Estado, pedido.Estado) { ... }
     ```
