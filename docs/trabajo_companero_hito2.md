# Guía de Trabajo · Hito 2 (Entrega y Demostración)
**Responsable:** Integrante 2  
**Proyecto:** RentCar API (`servA-gonzalez-reina`)  
**Fecha límite:** Miércoles 28 de octubre de 2026 a las 23:59  

---

## Objetivo
Tu objetivo es completar toda la capa contractual, la colección de pruebas en Bruno, la hoja de cotejo oficial del hito y las capturas de respaldo requeridas por la rúbrica del docente.

---

## Tarea 1: Completar el Contrato (`docs/openapi.yaml` y `docs/api.html`)

### 1.1. Contexto
Actualmente el contrato solo declara `vehiculos`, `reservas` y `login`. El docente exige que **todas** las entidades del modelo del Hito 1 estén en el contrato. Aquellas que aún no tienen endpoints programados en Go deben figurar en OpenAPI marcadas al final del `summary` con `(plan)`.

### 1.2. Qué debes hacer:
1. Abrir `docs/openapi.yaml`.
2. Agregar los esquemas de las entidades faltantes bajo `components.schemas`:
   - `Sucursal`: `ID`, `Nombre`, `Ciudad`, `Direccion`, `Activa`.
   - `Usuario`: `ID`, `SucursalID`, `Nombre`, `Correo`, `Rol`.
   - `Cliente`: `ID`, `Nombre`, `Correo`, `Telefono`, `Cedula`.
   - `Pago`: `ID`, `ReservaID`, `Monto`, `Metodo`, `EstadoPago`, `Referencia`, `Creado`.
3. Agregar las rutas bajo `paths` con el sufijo `(plan)` en su `summary`. Por ejemplo:
   - `GET /sucursales` -> `summary: Lista las sucursales (plan)`
   - `GET /clientes` -> `summary: Lista los clientes (plan)`
   - `GET /pagos` -> `summary: Lista los pagos (plan)`
   - `GET /usuarios` -> `summary: Lista los usuarios del sistema (plan)`
4. Validar que no existan errores ni advertencias ejecutando:
   ```bash
   bun x @redocly/cli lint docs/openapi.yaml --format=stylish
   ```
5. Regenerar el archivo HTML navegable:
   ```bash
   bun x @redocly/cli build-docs docs/openapi.yaml -o docs/api.html
   ```

---

## Tarea 2: Crear la Petición `0` en Bruno (`docs/hito2_bruno/`)

### 2.1. Contexto
En `docs/hito2_bruno/` ya existen las peticiones de la 1 a la 6, pero **falta la petición 0 de preparación**.
El documento del hito explica que la petición 5 "gasta" el registro al cambiar su estado, por lo que antes de la demostración (o de cada ensayo) se necesita crear un registro fresco y usar su ID en las peticiones 4 y 5.

### 2.2. Qué debes hacer:
1. Crear el archivo `docs/hito2_bruno/0_preparacion.bru`.
2. Configurar la petición:
   - **Método:** `POST`
   - **URL:** `http://localhost:8080/reservas`
   - **Headers:** `Authorization: Bearer {{token_ana}}` (o token de cliente autenticado)
   - **Body (JSON):**
     ```json
     {
       "ClienteID": 1,
       "VehiculoID": 1,
       "FechaInicio": "2026-11-01",
       "FechaFin": "2026-11-05",
       "Total": 140.00,
       "Estado": "pendiente"
     }
     ```
3. Verificar que al enviarla responda `201 Created` y devuelva el `ID` de la nueva reserva.

---

## Tarea 3: Redactar la Hoja del Hito (`docs/hito2_hoja.md`)

Crear el archivo `docs/hito2_hoja.md` con las 4 secciones exactas que exige la rúbrica:

### Estructura exacta requerida:

```markdown
# Hoja del Hito 2 · RentCar API

## 1 · Usuarios de ejemplo

| Usuario | Rol | Contraseña de ejemplo |
| :--- | :--- | :--- |
| ana | cliente | ana-2026 |
| luis | agente | luis-2026 |

## 2 · Las seis peticiones

| # | Petición | Con qué credencial | Debe responder |
| :-: | :--- | :--- | :-: |
| 1 | GET /reservas | Ninguna | 401 |
| 2 | POST /login con ana | Ninguna | 200 y el token |
| 3 | GET /reservas | Token de Ana (cliente) | 200 |
| 4 | PATCH /reservas/{id} -> en_curso | Token de Ana (cliente) | 403 |
| 5 | PATCH /reservas/{id} -> confirmada | Token de Luis (agente) | 200 |
| 6 | PATCH /reservas/1 -> finalizada (estando pendiente) | Token de Luis (agente) | 422 |

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Cliente | Agente |
| :--- | :---: | :---: | :---: |
| POST /login | Sí | Sí | Sí |
| GET /vehiculos | Sí | Sí | Sí |
| GET /reservas | No (401) | Sí | Sí |
| POST /reservas | No (401) | Sí | Sí |
| GET /reservas/{id} | No (401) | Sí | Sí |
| PUT /reservas/{id} | No (401) | No (403) | Sí |
| PATCH /reservas/{id} | No (401) | No (403) | Sí |
| DELETE /reservas/{id} | No (401) | No (403) | Sí |

## 4 · Recuadro de cumplimiento

- [x] Los structs de todas las entidades están en el código y sus tablas existen
- [x] El validador termina sin errores ni avisos
- [x] Las seis peticiones están en la colección, en orden
- [x] Cada una responde lo que dice esta hoja
- [x] go test ./... en verde, con una prueba de integración
- [x] El código de la entidad principal está en capas y la estructura está explicada en docs/decisiones.md
- [x] Hay una captura de respaldo de cada paso
```

---

## Tarea 4: Generar las Capturas de Respaldo (`docs/hito2_capturas/`)

Crea la carpeta `docs/hito2_capturas/` y guarda capturas de pantalla legibles con los siguientes nombres:

1. `1_get_reservas_401.png` (Respuesta 401 sin token).
2. `2_post_login_200.png` (Respuesta 200 entregando token).
3. `3_get_reservas_con_token_200.png` (Respuesta 200 listando reservas).
4. `4_patch_reserva_403.png` (Respuesta 403 con rol no autorizado).
5. `5_patch_reserva_200.png` (Respuesta 200 con rol agente).
6. `6_patch_reserva_prohibida_422.png` (Respuesta 422 por transición prohibida).
7. `7_go_test.png` (Terminal ejecutando `go test ./...` con todos los paquetes en verde).

---

## Tarea 5: Preparación de la Defensa en Vivo (5 Minutos)

El docente sorteará quién maneja el teclado:
* **Bloque 1 (Demostraciones 1 y 2):** El Contrato y "Quién es" (Peticiones 1, 2 y 3).
* **Bloque 2 (Demostraciones 3, 4 y 5):** "Qué puede hacer", "Regla de negocio" y "Estructura en capas".

### Frases a decir (sin improvisar):
* **Demostración 1:** *"Este es nuestro contrato. El validador no encuentra nada. Aquí está security, y aquí el 401, el 403 y el 422."*
* **Demostración 2:** *"Sin token, GET /reservas responde 401. Iniciamos sesión como ana y recibimos el token. Con el token, responde 200."*
* **Demostración 3:** *"ana es cliente y no puede cambiar el estado: 403. luis es agente y sí puede: responde 200."*
* **Demostración 4:** *"Pasar de pendiente a finalizada no está permitido en nuestro negocio: responde 422. Y estas son nuestras pruebas unitarias y de integración, en verde."*
* **Demostración 5:** *"Organizamos el código en capas: aquí se atiende la petición HTTP y aquí está la regla del cambio de estado; no se mezclan."*

### Preguntas teóricas de defensa (estudiar ambos):
1. **¿En qué línea de su código se decide que esta ruta pide token?**
   - En `internal/reservas/manejadores.go`, donde se agrupa bajo `conSesion := r.With(tokens.Autenticar)`.
2. **¿Qué lleva su token por dentro? Si alguien le cambia el rol, ¿qué responde su servidor y por qué?**
   - Lleva el ID del usuario, su nombre y su rol (en los claims JWT). Si alguien altera el payload, la firma criptográfica (HMAC-SHA256) no coincide con `cfg.ClaveFirma` y el servidor responde `401 Unauthorized` (token inválido/manipulado).
3. **¿Por qué esta ruta es solo de [rol] en su negocio y no de los dos?**
   - `PATCH /reservas/{id}` es solo de `agente` porque representa inspección física, verificación de pago y entrega de llaves. Un cliente no puede auto-aprobarse un alquiler.
4. **Si mañana alguien saca esta ruta del grupo protegido, ¿cuál de sus pruebas avisa?**
   - `internal/reservas/sin_token_test.go` (o `sesion_ayuda_test.go`), que verifica explícitamente que la ruta responda 401 al llamarse sin cabecera de autorización.
5. **En la transición prohibida, ¿por qué 422 y no 400 ni 403? ¿Dónde se valida?**
   - Porque la sintaxis JSON y los tipos son válidos (no es 400) y el usuario tiene permisos de agente (no es 403), pero la regla de negocio semántica no permite saltar de `pendiente` a `finalizada` sin entrega previa (422). Se valida en `internal/reservas/estado.go` con `puedePasar()`.
6. **¿Por qué eligieron esta estructura y no otra? Si mañana cambian de base de datos, ¿qué archivos tocan y cuáles no?**
   - Porque separa la capa de transporte HTTP de la persistencia y del dominio. Si cambiamos PostgreSQL por MongoDB, solo se tocan los archivos de persistencia/repositorio; los handlers HTTP y las reglas de negocio (`estado.go`) permanecen intactos.
