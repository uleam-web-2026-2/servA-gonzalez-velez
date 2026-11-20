# Las seis peticiones del hito 2 · Taller 12

| # | Petición | Esperado | Antes (taller 11) | Con las defensas | Con el usuario de la API (reto) |
| :---: | :--- | :---: | :---: | :---: | :---: |
| 1 | `GET /reservas` sin token | 401 | 401 | 401 | 401 |
| 2 | `POST /login` con `ana` | 200 | 200 | 200 | 200 |
| 3 | `GET /reservas` con token de `ana` | 200 | 200 | 200 | 200 |
| 4 | `PATCH /reservas/{id}` con token de `ana` (rol `cliente`) | 403 | 403 | 403 | 403 |
| 5 | `PATCH /reservas/{id}` con token de `luis` (`confirmada`) | 200 | 200 | 200 | 200 |
| 6 | `PATCH /reservas/{id}` con token de `luis` (transición prohibida a `finalizada`) | 422 | 422 | 422 | 422 |
