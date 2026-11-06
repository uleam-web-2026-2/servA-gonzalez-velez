-- name: BuscarReserva :one
SELECT id, cancha_id, cliente, inicio, horas, estado FROM reservas WHERE id = $1;

-- name: GuardarReserva :execrows
UPDATE reservas SET cancha_id = $2, cliente = $3, inicio = $4,
    horas = $5, estado = $6 WHERE id = $1;
