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
| 4 | PATCH /reservas/{reservaId} → en_curso | Token de Ana (cliente) | 403 |
| 5 | PATCH /reservas/{reservaId} → confirmada | Token de Luis (agente) | 200 |
| 6 | PATCH /reservas/1 → finalizada (estando confirmada) | Token de Luis (agente) | 422 |

La petición 0 crea la reserva nueva para las peticiones 4 y 5. La reserva de ejemplo
con ID 1 empieza en `confirmada`; por eso pasarla directamente a `finalizada` en
la petición 6 responde 422.

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Cliente | Agente |
| :--- | :---: | :---: | :---: |
| POST /login | Sí | Sí | Sí |
| GET /vehiculos | Sí | Sí | Sí |
| GET /reservas | No (401) | Sí | Sí |
| POST /reservas | No (401) | Sí | Sí |
| GET /reservas/{id} | No (401) | Sí | Sí |
| GET /sucursales | No (401) | Sí | Sí |
| GET /clientes | No (401) | Sí | Sí |
| GET /usuarios | No (401) | Sí | Sí |
| GET /pagos | No (401) | Sí | Sí |
| PUT /reservas/{id} | No (401) | No (403) | Sí |
| PATCH /reservas/{id} | No (401) | No (403) | Sí |
| DELETE /reservas/{id} | No (401) | No (403) | Sí |

## 4 · Recuadro de cumplimiento

- [x] Los structs de todas las entidades están en el código y sus tablas se crean al iniciar.
- [x] El validador termina sin errores ni avisos.
- [x] La colección incluye la petición 0 y las seis peticiones en orden.
- [x] Las siete peticiones de demostración/preparación respondieron lo que dice esta hoja en la instancia local.
- [x] `go test ./...` está en verde y hay una prueba de integración.
- [x] El código de reservas está organizado en capas y la estructura está explicada en `docs/decisiones.md`.
- [x] Hay una captura de respaldo de cada respuesta y de las pruebas.

Antes del turno, enviar la petición 2 para obtener un token vigente de Ana y otro inicio
de sesión con Luis; guardar su token en la variable activa de Bruno `tokenLuis`. Enviar
la petición 0 con el token de Ana; Bruno guarda el nuevo ID en `reservaId` para las
peticiones 4 y 5.
