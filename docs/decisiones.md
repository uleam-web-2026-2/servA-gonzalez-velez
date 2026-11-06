# Decisiones

## D1 · ¿Quién pone el estado al crear?

**Opciones:** que lo mande quien crea, o que el servidor ponga siempre el estado inicial.

**Qué elegimos:** el servidor asigna siempre el estado inicial `pendiente`. El DTO de entrada no acepta `ID` ni `Estado`; los cambios de estado se hacen aparte mediante `PATCH` y se validan con la máquina de transiciones.

**Por qué, en nuestro negocio:** en RentCar una reserva nueva debe esperar la verificación del pago antes de confirmarse. Fijar `pendiente` evita que un cliente se atribuya una confirmación o se salte las etapas del alquiler; un agente puede confirmarla después mediante `PATCH`.

**Qué pasaría con la otra opción:** si el cliente pudiera enviar el estado, podría crear reservas ya confirmadas o en etapas posteriores sin que se hubieran verificado las condiciones del negocio.

---

## D2 · ¿Qué pasa si se borra un registro que tiene relacionados?

**Opciones:** impedir la eliminación respondiendo un error de conflicto, o eliminar en cascada el registro y todos sus relacionados.

**Qué elegimos:** se impide la eliminación del vehículo si tiene reservas asociadas; la base de datos rechaza la operación por restricción de clave foránea y el servidor no permite eliminarlo.

**Por qué, en nuestro negocio:** los vehículos son activos de la empresa y las reservas son contratos. Si un vehículo ya fue reservado, borrarlo destruiría el histórico y los comprobantes (`Pago`) asociados.

**Qué pasaría con la otra opción:** el borrado en cascada eliminaría reservas y pagos, perdiendo respaldo contable.

---

## D3 · ¿Qué campos se pueden cambiar al actualizar?

**Opciones:** que se puedan cambiar todos los campos (incluido VehiculoID), o fijar la clave foránea y solo permitir cambiar cliente, fechas, total y estado.

**Qué elegimos:** al actualizar una reserva el `VehiculoID` y el estado quedan fijos; mediante `PUT` únicamente se pueden actualizar `ClienteID`, fechas y total. El cambio de estado se realiza por `PATCH` y pasa por la máquina de transiciones.

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
- `repositorio.go` declara la interfaz `Repositorio`; los adaptadores concretos están
  en `internal/almacen`.

Elegimos esta separación porque permite probar las reglas de negocio sin HTTP y
mantener las operaciones SQL en un único lugar. Si cambia la base de datos, se
reemplaza la implementación del repositorio y se conserva la lógica de negocio y
el manejo HTTP.

## Semana 10 · Fase 6 · transacciones

### D1 · El caso

El caso analizado es confirmar una reserva después de registrar su pago. En la
base de datos se escribirían el nuevo estado `confirmada` en `Reserva` y el
registro correspondiente en `Pago`. Si se guarda primero la reserva y falla la
creación del pago, el sistema mostraría una reserva confirmada sin el pago que
justifica esa confirmación. Una transacción debe hacer que ambas escrituras se
confirmen juntas o que ninguna quede aplicada.

En RentCar existen las entidades `Reserva` y `Pago`, pero el contrato actual no
incluye el registro de pagos ni una operación que combine ambas escrituras: el
cambio de estado guarda la reserva y `GET /pagos` solo consulta pagos; los datos
de ejemplo se cargan durante el sembrado. Por eso este caso se analiza como
diseño futuro y no se agrega como funcionalidad nueva en esta fase.

### D2 · Opciones que consideramos

1. **Unidad de Trabajo / Transactor en el repositorio.** El puerto de
   persistencia expone una operación como `DentroDeTransaccion(fn)`, que entrega
   al callback un repositorio ligado a la transacción. El servicio depende solo
   de interfaces propias del dominio; no importa `gorm` ni `database/sql`.
   Para las pruebas, un doble en memoria puede implementar el puerto y simular
   éxito o error sin una base real. La implementación concreta se ocupa de
   iniciar, confirmar o revertir la transacción.
2. **Transacción gestionada desde el servicio con un contexto o ejecutor
   transaccional.** El servicio coordina el inicio y las operaciones usando una
   abstracción que represente el contexto/ejecutor. Si esa abstracción expone
   directamente `*gorm.DB` o `*sql.Tx`, el servicio queda acoplado a
   infraestructura y se incumple la separación requerida. Si se define como
   interfaz propia, se puede probar con un doble, pero el servicio debe propagar
   ese ejecutor a todas las operaciones y las pruebas deben simular también su
   ciclo de vida; un `context.Context` por sí solo no hace atómicas las
   escrituras.

### D3 · La que elegimos y por qué

Elegimos la Unidad de Trabajo en el límite del repositorio para una futura
operación que requiera varias escrituras. Mantiene el servicio y sus pruebas
independientes de GORM, `database/sql` y una base real, y concentra el ciclo de
vida de la transacción en el adaptador de persistencia. El callback debe recibir
un repositorio ligado a esa misma transacción y todas las escrituras del caso
deben usarlo.

Esta decisión es de diseño; no se agrega el puerto ni la operación de confirmar
y cobrar ahora porque el caso no existe en el comportamiento actual de la API.

### D4 · La que descartamos y por qué

Descartamos que el servicio gestione directamente la transacción. Si recibe
`*gorm.DB` o `*sql.Tx`, se filtran detalles de infraestructura a la capa de
negocio y las pruebas requieren esos tipos concretos. Una abstracción propia
podría evitar esa fuga, pero obligaría a pasar el ejecutor a cada operación y a
simular su ciclo de vida en las pruebas; el límite transaccional seguiría
repartido entre el servicio y el adaptador. La Unidad de Trabajo mantiene ese
límite en persistencia.

### D5 · Con sqlc

En una implementación basada en `database/sql`, el adaptador obtiene un
`*sql.Tx` desde la conexión `*sql.DB` (por ejemplo, con `BeginTx`) y crea un
conjunto de consultas ligado a esa transacción con `q.WithTx(tx)`. Las consultas
de ambas escrituras se ejecutan mediante ese conjunto; ante un error se revierte
la transacción y, si ambas tienen éxito, se confirma. No se debe ejecutar una de
las escrituras mediante `q` ligado a `*sql.DB` mientras la otra usa `q.WithTx`,
porque entonces no pertenecerían a la misma unidad atómica.

### D6 · Archivos

Si se implementara el caso en el código, se revisarían o modificarían:

- `internal/reservas/repositorio.go`, para declarar el puerto de Unidad de
  Trabajo sin tipos de infraestructura.
- `internal/reservas/servicio.go`, para orquestar el caso de uso usando el
  repositorio transaccional recibido.
- `internal/reservas/servicio_test.go`, para probar éxito y fallo con un doble
  en memoria, sin base de datos real.
- `internal/almacen/reservas.go` y `internal/almacen/reservas_sqlc.go`, para
  implementar la unidad de trabajo en el adaptador y ligar sus operaciones a la
  misma transacción; si SQLC necesitara consultas nuevas, también
  `db/consultas.sql`, `db/schema.sql` y el código generado en
  `internal/almacen/filas/`.

En este análisis no se modifican los modelos, la base de datos, las rutas ni
OpenAPI: `Pago` ya existe y no hace falta inventar tablas para explicar la
atomicidad. Si el caso de uso se añadiera en el futuro, también habría que
actualizar sus pruebas de integración y el contrato HTTP correspondiente.
