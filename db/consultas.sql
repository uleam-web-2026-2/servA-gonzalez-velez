-- name: BuscarReserva :one
SELECT id, cliente_id, vehiculo_id, fecha_inicio, fecha_fin, total, estado
FROM reservas
WHERE id = $1;

-- name: GuardarReserva :execrows
UPDATE reservas
SET cliente_id = $2, vehiculo_id = $3, fecha_inicio = $4, fecha_fin = $5, total = $6, estado = $7
WHERE id = $1;
