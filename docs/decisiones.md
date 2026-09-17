# Decisiones

## Semana 3 · RentCar

## D3 · ¿Qué campos se pueden cambiar al actualizar?

**Opciones:** permitir cambiar todos los campos (incluido VehiculoID), o fijar la FK y solo editar cliente, fechas, total y estado.

**Qué elegimos:** al actualizar una reserva no se puede cambiar el `VehiculoID`; solo cliente, fechas, total y estado.

**Por qué, en nuestro negocio:** la reserva queda ligada al vehículo que el cliente eligió. Si se pudiera cambiar el vehículo a mitad de camino, el total y la disponibilidad dejarían de cuadrar con lo acordado. Para alquilar otro carro se cancela (o finaliza) y se crea una reserva nueva.

**Qué pasaría con la otra opción:** el cliente podría “mover” la reserva a otro vehículo sin pasar por disponibilidad ni recalcular el precio, y el historial del vehículo original perdería coherencia.

---

## Semana 2 · Mesa de ayuda (referencia)

## D1 · Orden de los middleware

- **Decisión:** Registro afuera, Recuperación adentro (`r.Use(middleware.Registro); r.Use(middleware.Recuperacion)`).
- **Por qué:** Si hay un pánico, Recuperación lo atrapa, escribe el 500 con la envoltura y termina normal. Registro, que está por fuera, sigue vivo y anota `GET /explotar → 500` con la duración real.
- **Alternativa descartada:** Recuperación afuera y Registro adentro. El pánico salta por encima de Registro; esa petición no deja línea en el log.

## D2 · Código para JSON roto frente a datos inválidos

- **Decisión:** JSON malformado → **400** (`cuerpo_malformado` / `json_invalido`); JSON válido que viola reglas → **422**.
- **Por qué:** Separar “cómo envías” de “qué envías”.

## D3 · Espacios en el título antes de validar

- **Decisión:** Se aplica `strings.TrimSpace` al título antes de medir la longitud y de guardarlo.
- **Por qué:** Un título solo con espacios no debe pasar la regla de longitud mínima.
