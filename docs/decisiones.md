# Decisiones

## D1 · ¿Quién pone el estado al crear?

**Opciones:** que lo mande quien crea, o que el servidor ponga siempre el estado inicial.

**Qué elegimos:** quien crea la reserva envía el estado, y el servidor valida que pertenezca a la lista de estados válidos (`estadosValidos`).

**Por qué, en nuestro negocio:** en RentCar una reserva puede originarse de distintas formas: una reserva web estándar ingresa como `pendiente` a la espera del comprobante de transferencia, mientras que una reserva tomada presencialmente en mostrador con garantía pagada puede registrarse directamente como `confirmada`. Permitir enviar el estado evita forzar peticiones HTTP extras para actualizar el estado inmediatamente después de crearla.

**Qué pasaría con la otra opción:** si el servidor forzara siempre `pendiente`, los operadores en agencia tendrían que hacer obligatoriamente un segundo llamado (`PATCH /reservas/{id}`) para confirmar la reserva, duplicando el tráfico y creando estados transitorios innecesarios.

---

## D2 · ¿Qué pasa si se borra un registro que tiene relacionados?

**Opciones:** impedir la eliminación respondiendo un error de conflicto, o eliminar en cascada el registro y todos sus relacionados.

**Qué elegimos:** se impide la eliminación del vehículo si tiene reservas asociadas; la base de datos rechaza la operación por restricción de clave foránea y el servidor no permite eliminarlo.

**Por qué, en nuestro negocio:** los vehículos son activos de la empresa y las reservas son contratos. Si un vehículo ya fue reservado, borrarlo destruiría el histórico y los comprobantes (`Pago`) asociados.

**Qué pasaría con la otra opción:** el borrado en cascada eliminaría reservas y pagos, perdiendo respaldo contable.

---

## D3 · ¿Qué campos se pueden cambiar al actualizar?

**Opciones:** que se puedan cambiar todos los campos (incluido VehiculoID), o fijar la clave foránea y solo permitir cambiar cliente, fechas, total y estado.

**Qué elegimos:** al actualizar una reserva el `VehiculoID` queda fijo; únicamente se pueden actualizar `ClienteID`, fechas, total y estado. El cambio de estado también pasa por la máquina de transiciones (igual que el `PATCH`).

**Por qué, en nuestro negocio:** una reserva representa un contrato pactado para una unidad específica (auto o camioneta). Si el cliente necesita otro vehículo, se cancela y se crea una nueva reserva.

**Qué pasaría con la otra opción:** se podría reasignar a otra categoría sin recalcular el total y romper disponibilidad.

---

## Semana 4 · D3 · La prueba que no escribimos

**Prueba que no escribimos:** la fila ★ de vehículo inexistente (`VehiculoID` que no está en la base → 422 `vehiculo_inexistente`).

**Dónde empieza la base:** `internal/reservas/manejadores.go` — `m.DB.First(&vehiculo, reserva.VehiculoID)`.

**Qué haría falta para probarla:** una base de datos de prueba (o un doble de `*gorm.DB`) que responda `ErrRecordNotFound` sin tocar PostgreSQL real; hoy las pruebas unitarias usan `DB: nil` y no pueden cruzar esa línea.

---

## Semana 7 · D1 · Qué nivel lleva cada ruta

| Ruta | Nivel | Por qué |
| --- | --- | --- |
| `POST /login` | Pública | Es la puerta por la que se consigue el token; cualquier usuario no autenticado debe poder enviar credenciales para autenticarse. |
| `GET /vehiculos` | Pública | Permite que potenciales clientes y visitantes exploren el catálogo de vehículos disponibles en la flota antes de registrarse o iniciar sesión. |
| `GET /reservas` | Con sesión | Permite consultar reservas registradas en el sistema; solo usuarios autenticados deben tener acceso a las operaciones del negocio. |
| `POST /reservas` | Con sesión | La creación de reservas requiere un usuario autenticado (cliente o agente) para vincular la solicitud con un usuario registrado y trazable. |
| `GET /reservas/{id}` | Con sesión | Ver el detalle de una reserva específica requiere autenticación previa para salvaguardar la privacidad de las reservas. |
| `PATCH /reservas/{id}` | Solo agente | El avance del flujo de negocio (confirmación de pago, entrega de vehículo, cierre de alquiler) requiere inspección física y financiera que solo el agente de sucursal está facultado a autorizar. |
| `PUT /reservas/{id}` | Solo agente | La modificación integral de una reserva (fechas, cliente, monto) es una potestad operativa reservada a los agentes para evitar fraudes en tarifas pactadas. |
| `DELETE /reservas/{id}` | Solo agente | La anulación o eliminación técnica de una reserva es de exclusiva competencia del personal de gestión de RentCar. |

## Semana 7 · D2 · La ruta de un solo rol

**Decisión:** `PATCH /reservas/{id}` es de acceso exclusivo para el rol `agente`.

**Por qué:** En RentCar, el cambio de estado de una reserva refleja eventos de negocio críticos en la sucursal: verificación del pago o depósito bancario (`confirmada`), entrega física de las llaves del vehículo con contrato firmado (`en_curso`), y recepción y liquidación del vehículo devuelto (`finalizada`). Si un `cliente` pudiera cambiar el estado, podría marcar su propia reserva como confirmada sin haber pagado, o marcarla en curso sin autorización, evadiendo la verificación de pago y el control de la flota.

**Qué descartamos:** Permitir que el `cliente` pueda cambiar su propio estado a `cancelada` directamente mediante `PATCH`. Lo descartamos para este taller para mantener la ruta con un único rol autorizado (`agente`), dejando que el cliente solicite la cancelación al agente o canalice la cancelación en una evolución futura mediante una ruta específica de autoservicio.

## Semana 7 · D3 · La transición prohibida

**Transición:** de `pendiente` a `finalizada` (y de `pendiente` a `en_curso`).

**Por qué el negocio no la permite:** En RentCar no se puede marcar una reserva como `finalizada` directamente desde `pendiente`. Una reserva en estado `pendiente` aún no ha tenido su pago verificado ni ha pasado por la entrega física de la unidad en sucursal (`en_curso`). Permitir finalizar una reserva sin entrega previa corrompería el histórico de cobro, los contratos y los registros de kilometraje y devolución de la flota.

**Qué responde y dónde se valida:** Responde código HTTP `422 Unprocessable Entity` con el mensaje `"de pendiente no se puede pasar a finalizada"` (o el estado solicitado). La transición se decide con `puedePasar` en `internal/reservas/servicio.go`; el handler HTTP de `internal/reservas/estado.go` convierte el resultado en la respuesta.

## Semana 8 · estructura

La entidad principal, `Reserva`, se organiza en tres responsabilidades dentro de
`internal/reservas/`:

- `manejadores.go` y `estado.go` reciben las peticiones HTTP, validan su formato y
  traducen los resultados a respuestas HTTP. No consultan GORM ni la base de datos.
- `servicio.go` contiene reglas independientes del transporte, como la validación
  de fechas y total y la máquina de transiciones de estado (`puedePasar`).
- `repositorio.go` concentra las consultas y escrituras de PostgreSQL mediante GORM.
  `Repositorio` es la interfaz usada por los manejadores y `RepositorioGorm` su
  implementación actual.

Elegimos esta separación porque permite probar las reglas de negocio sin HTTP y
mantener las operaciones SQL en un único lugar. Si cambia la base de datos, se
reemplaza la implementación del repositorio y se conserva la lógica de negocio y
el manejo HTTP.

## Semana 11 · el contenedor

### D1 · Qué entra en la imagen
La imagen final (`rentcar:latest`, `38.9MB` en disco / `11.7MB` de contenido) parte de `gcr.io/distroless/static-debian13:nonroot` y contiene únicamente el binario estático `/servidor` (compilado con `CGO_ENABLED=0`), certificados raíz, zonas horarias base, `/tmp` y el usuario sin privilegios `65532`. Todo el código fuente Go, el compilador `golang:1.27` y los módulos descargados se quedan en la etapa `compilacion` (`rentcar:compilacion`, `1.63GB` en disco / `389MB` de contenido). Si no existiera `.dockerignore`, el archivo `.env` sí entraría al contexto de construcción (`/src`) en la etapa de compilación, aunque no llegaría a la imagen final porque la segunda etapa solo copia `/servidor` con `COPY --from=compilacion`.

### D2 · Dónde está la base, vista desde el contenedor
`localhost` falló con `connection refused` en `127.0.0.1:5433` porque dentro del contenedor `localhost` apunta al propio contenedor de la API, donde no hay ningún servidor PostgreSQL escuchando. En `.env.docker` se reemplazó por `host.docker.internal` manteniendo el puerto `5433` para alcanzar el contenedor `pg` publicado en la máquina anfitriona. Cuando la API corra en un servidor real, `DATABASE_URL` deberá apuntar al nombre del servicio o dirección IP/host de la base de datos en esa red o proveedor (y nunca a `host.docker.internal`).

### D3 · Los puertos
Se usó `-p 8080:8080`: el primer `8080` es el puerto publicado en la máquina anfitriona (al que se conectan Bruno y el navegador) y el segundo `8080` es el puerto interno donde escucha el contenedor. El número de la derecha debe coincidir con la variable `PUERTO` (`8080`); si no coinciden, Docker reenvía el tráfico a un puerto donde el proceso no está escuchando y las peticiones fallan con conexión reiniciada o rechazada.

### D4 · La hora
En RentCar, la petición 3 devuelve `FechaInicio` y `FechaFin` como cadenas de fecha simple (`YYYY-MM-DD`, por ejemplo `"2026-09-20"`), por lo que la respuesta JSON fue idéntica entre `go run` y el contenedor; sin embargo, en los logs del contenedor se observa la diferencia de zona horaria (`2026/11/13 03:06` en UTC frente a `2026/11/12 22:06` en `-05:00` local). Si la API manejara `time.Time` con hora y zona horaria (por ejemplo, para registrar la hora exacta de entrega/devolución de vehículos o vencimiento de reservas por hora), en el contenedor saldrían en UTC (`Z`) representando el mismo instante, pero cualquier cálculo o corte por día local sin fijar la zona horaria (`America/Guayaquil` o `TZ`) podría desplazar reservas nocturnas al día siguiente. Decidimos mantener `YYYY-MM-DD` para las fechas de reserva y tener en cuenta UTC en los logs y futuros timestamps.

## Semana 11 · reto: la API y la base, juntas

### D5 · Cómo encuentra la API a la base
La API se conecta usando el nombre del servicio de PostgreSQL definido en `compose.yaml` (`host=db`) y el puerto interno `port=5432`. Docker Compose crea una red privada por defecto para el proyecto donde cada servicio se resuelve por su nombre DNS (`db`), y dentro de esa red PostgreSQL escucha en su puerto estándar (`5432`), no en el `5433` que solo se usaba para publicar el contenedor `pg` hacia Windows.

### D6 · El orden de arranque
Con la base recién creada, `depends_on` simple solo espera a que el contenedor `db` inicie, no a que PostgreSQL termine de inicializar la base de datos; si la API arranca antes, falla al conectar. Lo resolvimos configurando un `healthcheck` en `db` con `pg_isready -U postgres -d rentcar` y haciendo que `api` dependa de `db` con `condition: service_healthy`. Así se probó `docker compose down -v` seguido de `docker compose up -d --build` tres veces seguidas y las tres arrancaron sin errores.

### D7 · Dónde quedan los datos
Los datos de PostgreSQL se guardan en el volumen nombrado `pgdata` montado en `/var/lib/postgresql/data` dentro del servicio `db`. Al ejecutar `docker compose down` (sin `-v`), los contenedores y la red se eliminan pero el volumen `pgdata` se conserva, por lo que al volver a subir los registros siguen ahí. Solo `docker compose down -v` (o `docker volume rm`) borra el volumen y reinicia la base desde cero.

### D8 · Los secretos
`compose.yaml` no contiene contraseñas ni claves escritas: utiliza interpolación de variables de entorno (`${POSTGRES_PASSWORD:?...}`, `${CLAVE_FIRMA:?...}`), que Docker Compose lee automáticamente del archivo `.env` local (o del entorno del sistema/servidor). De este modo `compose.yaml` se puede subir al repositorio sin exponer secretos.

### D9 · Lo que se publica en la máquina
`compose.yaml` publica únicamente el puerto de la API (`${PUERTO:-8080}:${PUERTO:-8080}`) para que Bruno y los clientes lleguen en el puerto habitual sin chocar con el contenedor `pg` de la semana 3 (que sigue usando el `5433`). El servicio `db` no publica ningún puerto hacia la máquina anfitriona; cuando necesitemos consultarlo con `psql`, podemos ejecutarlo directamente dentro del contenedor de Compose con `docker compose exec db psql -U postgres -d rentcar`.


## Semana 12 · la seguridad de la API

### D1 · Lo que encontramos y cerramos

La búsqueda con `git grep -n -E -f docs/taller12/patrones.txt -- "*.go" ":(exclude)*_test.go"` devolvió **27 líneas** (guardadas en `docs/taller12/auditoria.txt` y revisadas una por una en `docs/taller12/auditoria.md`). De esas 27 líneas, **8 arman SQL con datos que vienen del cliente**:

- `internal/almacen/filas/consultas.sql.go:19` y `:50` (código generado por `sqlc`): usan consultas fijas con parámetros posicionales `$1`..`$7` (`QueryRowContext` y `ExecContext`), e `id` llega convertido a `uint` por `strconv.Atoi`.
- `internal/almacen/reservas.go:23`, `:32` y `:61` (`First` de cliente, vehículo y reserva): reciben un `id` tipado como `uint` (ya convertido con `strconv.Atoi` en `leerID` o decodificado como entero desde JSON), por lo que GORM lo envía como parámetro numérico y nunca como condición de texto.
- `internal/almacen/reservas.go:47` y `:55` (`Where("estado = ?", estado)` y `query.Find(&lista)`): el filtro `?estado=` viaja aparte mediante el marcador `?` de GORM, y `limit`/`offset` se validan antes con `strconv.Atoi` como `int`.
- `internal/sesion/usuarios.go:86` (`Where("nombre = ?", pedido.Usuario)`): el usuario de `POST /login` también se pasa aparte con `?`.

Las demás líneas encontradas corresponden a listados sin entrada del cliente (`Find` de catálogos y `Sembrar`), firmas de la interfaz `DBTX`, mensajes de error o lectura de `r.URL.Query()`. **No hubo consultas vulnerables que cerrar**, por lo que no fue necesario modificar código ni generar `auditoria_despues.txt`.

En la prueba del paso 1.4 (`docs/taller12/inyeccion.txt`), `GET /reservas?estado=pendiente` respondió `200 OK` con **2 reservas** pendientes (`ID` 2 y 7); al enviar el ataque `GET /reservas?estado=' OR 1=1 --`, respondió `200 OK` con `{"ok":true,"datos":[]}` (0 registros, demostrando que el texto se trató como dato y no como SQL); y al enviar `GET /reservas/1%20OR%201=1`, `leerID` lo rechazó antes de tocar la base con `400 Bad Request` (`id_invalido`).

### D2 · El origen que puede leer nuestras respuestas

Configuramos `AllowedOrigins: []string{"http://localhost:5173"}` (sin barra final), que es el origen local donde corre el front con Vite. No usamos `"*"` porque aplicamos principio de mínimo privilegio: con `"*"` cualquier página maliciosa en otro dominio podría leer las respuestas de nuestra API desde el navegador de un usuario, además de impedir un control estricto sobre qué cliente web está autorizado.

El middleware `cors.Handler(...)` se colocó **arriba de todo**, inmediatamente después de `r := chi.NewRouter()` y antes de `middleware.Registro`, el guardia de autenticación (`tokens.Autenticar`) y las rutas. Esto es necesario porque las peticiones previas (*preflight* `OPTIONS`) que envía el navegador no llevan cabecera `Authorization` ni tienen manejador propio: si pasaran antes por el guardia recibirían `401 Unauthorized`, y si se aplicaran solo con `r.With(...)` recibirían `405 Method Not Allowed`.

Al probar con `curl` un `PATCH /reservas/8` enviando `Origin: http://malo.example` y un token válido de `luis` (`docs/taller12/cors_curl_ejecuta.txt`), el servidor respondió `200 OK` y cambió la reserva a `confirmada`. Esto confirma que **CORS no es un muro en el servidor sino una política que hace cumplir el navegador**: `curl` y Bruno ignoran las cabeceras `Access-Control-*`, por lo que la seguridad real de esa operación depende del token JWT y del control de rol (`SoloRol`).

Cuando publiquemos el front en producción, deberemos agregar o reemplazar `http://localhost:5173` por el origen público exacto (esquema `https://`, dominio y puerto si aplica, sin barra final), idealmente leyéndolo desde una variable de entorno (por ejemplo `ORIGEN_FRONT`) en lugar de dejarlo fijo en el código.

### D3 · El límite del login

Configuramos el límite de intentos en **`5` peticiones por ventana de `1 minuto` por IP** mediante `httprate.LimitBy(5, time.Minute, porIP)`. Elegimos 5 intentos por minuto porque es suficiente margen para un usuario legítimo que se equivoque al escribir su contraseña un par de veces (o para iniciar sesión con dos cuentas distintas durante pruebas como `ana` y `luis`), pero frena ataques de fuerza bruta y relleno de credenciales automatizados.

El middleware se montó **solo en `POST /login`** usando `r.With(...)` y no de forma global con `r.Use(...)`: ponerlo en todo el router haría que las operaciones normales autenticadas (o las 6 peticiones consecutivas del flujo del hito 2) chocaran contra el tope de 5 por minuto y devolvieran `429`.

En la prueba de ráfaga de 7 intentos consecutivos (`docs/taller12/rafaga.txt`), las peticiones 1 a 5 devolvieron `401 Unauthorized` con `X-Ratelimit-Limit: 5` y `X-Ratelimit-Remaining` descendiendo de `4` a `0`. El primer **`429 Too Many Requests`** llegó en la **6.ª petición** (y se repitió en la 7.ª), incluyendo la cabecera **`Retry-After: 60`** y `X-Ratelimit-Remaining: 0`.

### D4 · Detrás de un proxy

El middleware `chimiddleware.ClientIPFromRemoteAddr` obtiene la IP directamente de **`r.RemoteAddr`**, es decir, del socket TCP de la conexión que llega al servidor Go, sin leer cabeceras HTTP.

Si publicáramos la API detrás de un proxy reverso (como Nginx, Caddy, Cloudflare o un balanceador de carga en la nube) sin ajustar esto, todas las conexiones TCP hacia nuestra API provendrían de la misma IP interna del proxy. En consecuencia, `porIP` agruparía a todos los usuarios bajo la misma clave y bastaría con que una sola persona hiciera 5 intentos en `/login` para bloquear el inicio de sesión de todos los demás usuarios durante un minuto.

Aun así, no se puede confiar ciegamente en cabeceras como `X-Forwarded-For` o `X-Real-IP`, porque cualquier cliente puede enviar esas cabeceras falsificadas con IPs distintas en cada petición y evadir por completo el límite de tasa. Cuando haya un proxy reverso en producción, la solución segura es aceptar `X-Forwarded-For` **únicamente** si la conexión proviene de la IP confiable del proxy (y tomar la IP que agregó el proxy, no la que pudo inyectar el cliente).
