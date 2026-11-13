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

La imagen final `rentcar:latest` ocupa **38.9 MB** en disco (**11.7 MB** de contenido) y contiene el binario estático `/servidor` sobre la base mínima `distroless`. La etapa de compilación `rentcar:compilacion` ocupa **1.63 GB** (**389 MB** de contenido) porque conserva Go, las dependencias y el código fuente usados para compilar. El `.dockerignore` evita que `.env` entre al contexto de construcción. Sin ese archivo, `COPY . .` lo copiaría a la etapa de compilación; no llegaría a la imagen final porque solo copiamos `/servidor`, pero quedaría expuesto en esa etapa y en la caché de construcción.

### D2 · Dónde está la base, vista desde el contenedor

`localhost` dentro del contenedor apunta al propio contenedor, donde no corre PostgreSQL; por eso la conexión al puerto `5433` fue rechazada. En Docker Desktop usamos `host.docker.internal`, que permite llegar al host y al puerto `5433` publicado por el contenedor de PostgreSQL. En un servidor real evitaríamos depender de la máquina anfitriona: conectaríamos la API y PostgreSQL a una red privada y usaríamos el nombre del servicio de la base.

### D3 · Los puertos

Arrancamos con `-p 8080:8080`: el número de la izquierda (`8080`) es el puerto publicado en la máquina y al que se conecta Bruno; el de la derecha (`8080`) es el puerto dentro del contenedor. El número de la derecha debe coincidir con `PUERTO`, que aquí es `8080`. Si no coincide, Docker enviaría el tráfico a un puerto donde el servidor no escucha y las peticiones no llegarían a la API.

### D4 · La hora

Comparamos la respuesta completa de `GET /reservas` con `go run` y con el contenedor usando la misma base y el mismo estado: fueron idénticas. Los campos de fecha salen como `YYYY-MM-DD`, sin hora ni zona horaria, así que en esta respuesta no hay un instante que cambie de zona. El proceso dentro del contenedor usa UTC, mientras que la máquina de desarrollo está en la zona de Ecuador; la diferencia se observa en las horas de los logs. Si RentCar incorporara horas de retiro y devolución, habría que guardar instantes con zona explícita y mostrarlos en la zona del negocio para evitar confundir el horario o el día de entrega de un vehículo.
