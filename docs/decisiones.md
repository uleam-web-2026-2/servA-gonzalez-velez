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
