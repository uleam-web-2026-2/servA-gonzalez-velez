## D1 · Orden de los middleware

- **Decisión:** Registro afuera, Recuperación adentro (`r.Use(middleware.Registro); r.Use(middleware.Recuperacion)`).
- **Por qué:** Si hay un pánico, Recuperación lo atrapa, escribe el 500 con envoltura y termina normal. Registro, que está por fuera, sigue vivo y anota `GET /explotar → 500` con la duración real. Así el log refleja también los fallos internos, no solo los éxitos.
- **Alternativa descartada:** Recuperación afuera y Registro adentro. El pánico salta por encima de Registro (no tiene `recover`), Recuperación responde 500, pero esa petición no deja línea en el log: el servidor respondió y el registro no se enteró.
- **Cómo lo comprobamos:** `GET /explotar` → cuerpo `{"ok":false,"error":{"codigo":"error_interno",...}}` con estado 500; en la terminal del servidor aparecen `PÁNICO en GET /explotar: ...` y `GET /explotar → 500 (...)`.

## D2 · Código para JSON roto frente a datos inválidos

- **Decisión:** JSON malformado → **400** (`cuerpo_malformado`); JSON válido que viola reglas de negocio → **422** (`datos_invalidos`).
- **Por qué:** Un 400 le dice al cliente «arregla *cómo* envías» (sintaxis/estructura). Un 422 le dice «te entendí, arregla *qué* envías» (título corto, prioridad fuera de baja/media/alta). Separarlos evita que el cliente trate un typo de JSON igual que un dato de negocio incorrecto.
- **Alternativa descartada:** Usar 400 para ambos. Es defendible y más simple, pero pierde la distinción entre fallo de formato y fallo de regla, que es justamente lo que el contrato de la mesa de ayuda quiere hacer visible.
- **Cómo lo comprobamos:** `POST` con `{"titulo": "x` → 400 `cuerpo_malformado`; `POST` con título `"abc"` o prioridad `"urgente"` → 422 `datos_invalidos`.

## D3 · Espacios en el título antes de validar

- **Decisión:** Se aplica `strings.TrimSpace` al título antes de medir la longitud y de guardarlo.
- **Por qué:** Un título como `"     "` (solo espacios) no debería pasar la regla de «al menos 5 caracteres». Recortar también evita tickets que lucen vacíos o con padding accidental en la lista.
- **Alternativa descartada:** Validar `len(entrada.Titulo)` sin recortar. Aceptaría cuerpos con espacios de relleno que engañan la regla de longitud mínima.
- **Cómo lo comprobamos:** `POST` con `{"titulo":"  abc  ","prioridad":"alta"}` → 422 `datos_invalidos`; `POST` con `{"titulo":"  Teclado dañado  ","prioridad":"media"}` → 201 y el título guardado sin espacios extremos.
