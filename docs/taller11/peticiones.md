# Semana 11 · las seis peticiones del hito 2

| Petición | Con go run | En el contenedor | Con Compose |
| --- | --- | --- | --- |
| 1 sin token | 401 | 401 | 401 |
| 2 login | 200 | 200 | 200 |
| 3 con token | 200 | 200 | 200 |
| 4 rol sin permiso | 403 | 403 | 403 |
| 5 cambio permitido | 200 | 200 | 200 |
| 6 transición prohibida | 422 | 422 | 422 |

Petición 3 con go run (respuesta completa):

```json
{"ok":true,"datos":[{"ID":1,"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-09-20","FechaFin":"2026-09-23","Total":105,"Estado":"confirmada"},{"ID":2,"ClienteID":2,"VehiculoID":1,"FechaInicio":"2026-10-01","FechaFin":"2026-10-03","Total":70,"Estado":"pendiente"},{"ID":3,"ClienteID":3,"VehiculoID":2,"FechaInicio":"2026-09-18","FechaFin":"2026-09-19","Total":28,"Estado":"en_curso"},{"ID":4,"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-05","Total":140,"Estado":"pendiente"}]}
```

La misma, en el contenedor (paso 2.5):

```json
{"ok":true,"datos":[{"ID":1,"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-09-20","FechaFin":"2026-09-23","Total":105,"Estado":"confirmada"},{"ID":2,"ClienteID":2,"VehiculoID":1,"FechaInicio":"2026-10-01","FechaFin":"2026-10-03","Total":70,"Estado":"pendiente"},{"ID":3,"ClienteID":3,"VehiculoID":2,"FechaInicio":"2026-09-18","FechaFin":"2026-09-19","Total":28,"Estado":"en_curso"},{"ID":4,"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-05","Total":140,"Estado":"confirmada"},{"ID":5,"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-05","Total":140,"Estado":"pendiente"}]}
```
