# Auditoría de inyección SQL · Taller 12

Resultado de `git grep -n -E -f docs/taller12/patrones.txt -- "*.go" ":(exclude)*_test.go"` (27 líneas encontradas en `docs/taller12/auditoria.txt`).

| Línea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Qué hicimos |
| :--- | :--- | :--- | :--- | :--- |
| `internal/almacen/filas/consultas.sql.go:19` (`QueryRowContext(ctx, buscarReserva, id)`) | Sí | Sí, el `{id}` de la ruta | Sí: sqlc (`$1`) e `id` es `uint` (`strconv.Atoi`) | Nada: es código generado seguro. |
| `internal/almacen/filas/consultas.sql.go:50` (`ExecContext(ctx, guardarReserva, ...)`) | Sí | Sí, campos del cuerpo JSON | Sí: sqlc (`$1`..`$7`) | Nada: es código generado seguro. |
| `internal/almacen/reservas.go:23` (`First(&cliente, id)`) | Sí | Sí, `ClienteID` del cuerpo JSON | Sí: `id` es `uint`, parámetro de GORM | Nada: es segura. |
| `internal/almacen/reservas.go:32` (`First(&vehiculo, id)`) | Sí | Sí, `VehiculoID` del cuerpo JSON | Sí: `id` es `uint`, parámetro de GORM | Nada: es segura. |
| `internal/almacen/reservas.go:47` (`Where("estado = ?", estado)`) | Sí | Sí, `?estado=` de `r.URL.Query()` | Sí, con `?` en GORM | Nada: es segura. |
| `internal/almacen/reservas.go:55` (`query.Find(&lista)`) | Sí | Sí, `?estado=`, `?limit=`, `?offset=` | Sí: `estado` con `?`, `limit`/`offset` son `int` (`strconv.Atoi`) | Nada: es segura. |
| `internal/almacen/reservas.go:61` (`First(&reserva, id)`) | Sí | Sí, el `{id}` de la ruta | Sí: `id` es `uint` convertido con `strconv.Atoi` en `leerID` | Nada: es segura. |
| `internal/almacen/reservas.go:79, 85, 91, 97, 103` (`Find(&lista)`) | Sí | No (listados sin filtros del cliente) | — | Nada: es segura. |
| `internal/sesion/usuarios.go:55` (`Where("nombre = ?", e.nombre).Limit(1).Find(&u)`) | Sí | No (datos fijos de siembra) | Sí, con `?` en GORM | Nada: es segura. |
| `internal/sesion/usuarios.go:86` (`Where("nombre = ?", pedido.Usuario).Limit(1).Find(&u)`) | Sí | Sí, `usuario` del JSON de `POST /login` | Sí, con `?` en GORM | Nada: es segura. |
| `internal/almacen/filas/db.go:13, 15, 16` (`DBTX`), `internal/config/config.go:33, 39, 55`, `internal/reservas/estado.go:33, 53`, `internal/reservas/manejadores.go:96, 330, 331`, `internal/sesion/guardia.go:60`, `main.go:89` | No: firmas de interfaz, mensajes de error, lectura de `r.URL.Query()` y dirección del servidor | — | — | Nada. |

**Conclusión:** No se encontró ninguna consulta SQL vulnerable (todas las consultas con datos del cliente usan parámetros `?` en GORM, `$1..$N` en sqlc, o conversión numérica previa con `strconv.Atoi`). No fue necesario modificar código ni generar `auditoria_despues.txt`.
