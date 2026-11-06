CREATE TABLE reservas (
    id bigserial PRIMARY KEY,
    cliente_id bigint NOT NULL,
    vehiculo_id bigint NOT NULL,
    fecha_inicio text NOT NULL,
    fecha_fin text NOT NULL,
    total double precision NOT NULL,
    estado text NOT NULL
);
