# Decisiones

## D1 · ¿Quién pone el estado al crear?

**Opciones:** que lo mande quien crea, o que el servidor ponga siempre el estado inicial.

**Qué elegimos:** quien crea la reserva envía el estado, y el servidor valida que pertenezca a la lista de estados válidos (`estadosValidos`).

**Por qué, en nuestro negocio:** en RentCar una reserva puede originarse de distintas formas: una reserva web estándar ingresa como `pendiente` a la espera de depósito, mientras que una reserva tomada presencialmente en mostrador con garantía pagada puede registrarse directamente como `confirmada`. Permitir enviar el estado evita forzar peticiones HTTP extras para actualizar el estado inmediatamente después de crearla.

**Qué pasaría con la otra opción:** si el servidor forzara siempre `pendiente`, los operadores en agencia o integraciones de pago tendrían que hacer obligatoriamente un segundo llamado (`PUT /reservas/{id}`) para confirmar la reserva, duplicando el tráfico de red y creando estados transitorios innecesarios.

---

## D2 · ¿Qué pasa si se borra un registro que tiene relacionados?

**Opciones:** impedir la eliminación respondiendo un error de conflicto, o eliminar en cascada el registro y todos sus relacionados.

**Qué elegimos:** se impide la eliminación del vehículo si tiene reservas asociadas; la base de datos rechaza la operación por restricción de clave foránea (`FOREIGN KEY CONSTRAINT`) y el servidor no permite eliminarlo.

**Por qué, en nuestro negocio:** los vehículos representan los activos de la empresa y las reservas constituyen contratos legales y transacciones financieras. Si un vehículo ya fue reservado, su eliminación destruiría el histórico de uso, los comprobantes de alquiler y la constancia de entrega de garantías de los clientes. Si un vehículo se retira de la flota o se vende, debe inactivarse operativamente, nunca borrarse físicamente si tiene registros asociados.

**Qué pasaría con la otra opción:** si se aplicara borrado en cascada, eliminar un vehículo borraría de golpe todas las reservas pasadas y futuras asociadas a él, perdiendo el respaldo de ingresos contables y cancelando alquileres futuros sin previo aviso.

---

## D3 · ¿Qué campos se pueden cambiar al actualizar?

**Opciones:** que se puedan cambiar todos los campos (incluido VehiculoID), o fijar la clave foránea y solo permitir cambiar cliente, fechas, total y estado.

**Qué elegimos:** al actualizar una reserva el `VehiculoID` queda fijo y no se puede modificar; únicamente se pueden actualizar el cliente, las fechas (inicio y fin), el total y el estado.

**Por qué, en nuestro negocio:** una reserva en RentCar representa un contrato de alquiler pactado para una unidad física específica (auto, moto o camioneta) con tarifa diaria y calendario de disponibilidad asignados. Si se permitiera cambiar el vehículo dentro de la misma reserva, se rompería la consistencia entre el vehículo asignado, el cálculo del total y el control de disponibilidad de la flota. Si el cliente necesita otro vehículo, la reserva debe cancelarse y crearse una nueva.

**Qué pasaría con la otra opción:** el cliente o usuario podría reasignar la reserva a un vehículo de mayor categoría o diferente tipo sin recalcular el importe, o causar sobreventa y conflictos de disponibilidad en el nuevo vehículo, dejando además registros inconsistentes en el historial de uso del vehículo originalmente contratado.
