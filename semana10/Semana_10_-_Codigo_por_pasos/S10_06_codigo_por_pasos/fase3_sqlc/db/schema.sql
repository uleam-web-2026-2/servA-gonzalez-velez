CREATE TABLE reservas (
    id        bigserial PRIMARY KEY,
    cancha_id bigint NOT NULL,
    cliente   text NOT NULL,
    inicio    timestamptz NOT NULL,
    horas     bigint NOT NULL,
    estado    text NOT NULL
);
