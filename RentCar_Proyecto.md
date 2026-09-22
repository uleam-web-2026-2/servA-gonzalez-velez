# RentCar — Sistema de alquiler de vehículos

## 1. Descripción general

RentCar será un sistema web de alquiler de vehículos desarrollado para la asignatura **Aplicaciones para el Servidor Web**.

La aplicación simulará una empresa donde los clientes podrán consultar motos, autos y camionetas disponibles, revisar precios y realizar reservas desde una página web. El administrador podrá gestionar vehículos, clientes y reservas.

## 2. Objetivo

Crear una aplicación web conectada a un servidor desarrollado en Go y una base de datos MySQL, aplicando:

- Servidor web.
- API REST.
- Métodos HTTP.
- CRUD.
- JSON.
- Validaciones.
- Middleware.
- Base de datos.
- Comunicación frontend-backend.
- Reglas de negocio.

## 3. Tecnologías

| Tecnología | Uso |
|---|---|
| Go | Servidor web y API REST |
| MySQL | Base de datos |
| HTML | Estructura de la página |
| CSS | Diseño visual y responsive |
| JavaScript | Comunicación con la API |
| Postman | Pruebas de endpoints |
| VS Code | Desarrollo |

## 4. Funcionalidades

### Página de inicio

- Logo de RentCar.
- Presentación del negocio.
- Botón “Explorar vehículos”.
- Vehículos destacados.

### Catálogo

- Motos, autos y camionetas.
- Marca y modelo.
- Tipo de vehículo.
- Precio por día.
- Imagen.
- Estado de disponibilidad.
- Botón “Alquilar ahora”.

### Formulario de reserva

- Nombre del cliente.
- Correo.
- Teléfono.
- Vehículo.
- Fecha de inicio.
- Fecha de devolución.
- Total automático.

### Panel administrativo

- Agregar vehículos.
- Editar vehículos.
- Eliminar vehículos.
- Consultar vehículos.
- Ver reservas.
- Cambiar estados.
- Consultar clientes.

## 5. Funcionamiento técnico

```text
Cliente abre RentCar
        |
        v
HTML + CSS + JavaScript
        |
        v
API REST en Go
        |
        v
Validar datos y reglas de negocio
        |
        v
Base de datos MySQL
        |
        v
Respuesta JSON
        |
        v
Mostrar resultado en la página
```

## 6. Endpoints

| Método | Endpoint | Acción |
|---|---|---|
| GET | `/api/vehiculos` | Listar vehículos |
| GET | `/api/vehiculos/{id}` | Consultar vehículo |
| POST | `/api/vehiculos` | Registrar vehículo |
| PATCH | `/api/vehiculos/{id}` | Editar vehículo |
| DELETE | `/api/vehiculos/{id}` | Eliminar vehículo |
| GET | `/api/clientes` | Listar clientes |
| POST | `/api/clientes` | Registrar cliente |
| GET | `/api/reservas` | Listar reservas |
| POST | `/api/reservas` | Crear reserva |
| PATCH | `/api/reservas/{id}` | Cambiar estado |

### Significado

- **GET:** consultar información.
- **POST:** crear registros.
- **PATCH:** modificar parcialmente.
- **DELETE:** eliminar registros.

## 7. Estructura propuesta

```text
RentCar/
├── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── middleware/
│   │   └── middleware.go
│   ├── respuesta/
│   │   └── respuesta.go
│   ├── vehiculos/
│   │   ├── vehiculos.go
│   │   ├── handler.go
│   │   └── repository.go
│   ├── clientes/
│   │   ├── clientes.go
│   │   └── handler.go
│   └── reservas/
│       ├── reservas.go
│       └── handler.go
├── config/
│   └── db.go
├── frontend/
│   ├── index.html
│   ├── estilos.css
│   └── app.js
└── docs/
    └── decisiones.md
```

## 8. Entidades

### Vehículos

- ID.
- Marca.
- Modelo.
- Tipo.
- Precio por día.
- Estado.
- Imagen.

### Clientes

- ID.
- Nombre.
- Correo.
- Teléfono.

### Reservas

- ID.
- Cliente.
- Vehículo.
- Fecha de inicio.
- Fecha de devolución.
- Total.
- Estado.

## 9. Base de datos MySQL

```sql
CREATE DATABASE rentcar;
USE rentcar;
```

### Tabla vehículos

```sql
CREATE TABLE vehiculos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    marca VARCHAR(50) NOT NULL,
    modelo VARCHAR(50) NOT NULL,
    tipo VARCHAR(30) NOT NULL,
    precio_dia DECIMAL(10,2) NOT NULL,
    estado VARCHAR(20) DEFAULT 'disponible',
    imagen VARCHAR(255)
);
```

### Tabla clientes

```sql
CREATE TABLE clientes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    correo VARCHAR(100) NOT NULL,
    telefono VARCHAR(20)
);
```

### Tabla reservas

```sql
CREATE TABLE reservas (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cliente_id INT NOT NULL,
    vehiculo_id INT NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    estado VARCHAR(20) DEFAULT 'pendiente',
    FOREIGN KEY (cliente_id) REFERENCES clientes(id),
    FOREIGN KEY (vehiculo_id) REFERENCES vehiculos(id)
);
```

## 10. Cálculo del alquiler

```text
Total = Precio por día × Número de días
```

Ejemplo:

```text
Precio diario: $35
Días: 3
Total: $105
```

Posibles mejoras:

- Descuento por alquiler prolongado.
- Recargo por entrega a domicilio.
- Precio según tipo de vehículo.
- Control de disponibilidad.
- Historial de reservas.
- Penalización por devolución tardía.

## 11. Reglas de negocio

1. Los campos obligatorios no pueden estar vacíos.
2. El precio debe ser mayor que cero.
3. La fecha de devolución debe ser posterior a la fecha de inicio.
4. El vehículo debe estar disponible.
5. El cliente debe existir.
6. El vehículo debe existir.
7. No se deben permitir reservas en fechas ocupadas.
8. El total debe calcularse en el servidor.
9. Los estados deben ser válidos.
10. Los errores deben devolverse en JSON.

### Estados de vehículos

- disponible
- alquilado
- mantenimiento

### Estados de reservas

- pendiente
- confirmada
- en_curso
- finalizada
- cancelada

## 12. Middleware

Un middleware se ejecuta durante el procesamiento de una solicitud HTTP.

### Middleware de registro

Registra:

- Método HTTP.
- Ruta.
- Código de respuesta.
- Tiempo de respuesta.

Ejemplo:

```text
GET /api/vehiculos → 200
POST /api/reservas → 201
GET /api/vehiculos/99 → 404
```

### Middleware de recuperación

Captura errores inesperados y devuelve una respuesta HTTP 500.

```json
{
  "ok": false,
  "error": {
    "codigo": "error_interno",
    "mensaje": "ocurrió un error inesperado"
  }
}
```

## 13. Respuestas JSON

### Éxito

```json
{
  "ok": true,
  "datos": {}
}
```

### Error

```json
{
  "ok": false,
  "error": {
    "codigo": "datos_invalidos",
    "mensaje": "el precio debe ser mayor que cero"
  }
}
```

## 14. Códigos HTTP

| Código | Significado | Ejemplo |
|---|---|---|
| 200 | OK | Consulta exitosa |
| 201 | Created | Registro creado |
| 400 | Bad Request | JSON mal formado |
| 404 | Not Found | Registro inexistente |
| 405 | Method Not Allowed | Método no permitido |
| 422 | Unprocessable Entity | Datos inválidos |
| 500 | Internal Server Error | Error inesperado |

## 15. Pruebas en Postman

### Listar vehículos

```http
GET http://localhost:8080/api/vehiculos
```

### Crear vehículo

```http
POST http://localhost:8080/api/vehiculos
```

Body:

```json
{
  "marca": "Toyota",
  "modelo": "Corolla",
  "tipo": "auto",
  "precio_dia": 35,
  "estado": "disponible",
  "imagen": "corolla.jpg"
}
```

### Consultar vehículo

```http
GET http://localhost:8080/api/vehiculos/1
```

### Crear cliente

```http
POST http://localhost:8080/api/clientes
```

Body:

```json
{
  "nombre": "Jipson González",
  "correo": "jipson@example.com",
  "telefono": "0999999999"
}
```

### Crear reserva

```http
POST http://localhost:8080/api/reservas
```

Body:

```json
{
  "cliente_id": 1,
  "vehiculo_id": 1,
  "fecha_inicio": "2026-09-20",
  "fecha_fin": "2026-09-23"
}
```

### Consultar reservas

```http
GET http://localhost:8080/api/reservas
```

### Actualizar reserva

```http
PATCH http://localhost:8080/api/reservas/1
```

Body:

```json
{
  "estado": "confirmada"
}
```

### Eliminar vehículo

```http
DELETE http://localhost:8080/api/vehiculos/1
```

## 16. Fases de desarrollo

### Fase 1: Preparación

- Crear carpeta RentCar.
- Configurar Go.
- Crear servidor.
- Instalar dependencias.
- Crear base de datos.
- Configurar conexión.

### Fase 2: API REST

- CRUD de vehículos.
- CRUD de clientes.
- Crear reservas.
- Consultar reservas.
- Validaciones.
- Middleware.
- Respuestas JSON.

### Fase 3: Página web

- Diseñar inicio.
- Crear catálogo.
- Mostrar vehículos desde la API.
- Crear formulario.
- Calcular total.
- Crear panel administrativo.

### Fase 4: Pruebas

- Probar endpoints.
- Probar datos válidos e inválidos.
- Probar registros inexistentes.
- Probar reservas duplicadas.
- Probar errores.

### Fase 5: Presentación

- Explicar objetivo.
- Mostrar página web.
- Mostrar servidor.
- Ejecutar consultas en Postman.
- Explicar base de datos.
- Mostrar el flujo de una reserva.

## 17. ¿Qué quiere enseñar?

RentCar busca enseñar:

- Cómo funciona un servidor web.
- Cómo crear una API REST.
- Cómo aplicar GET, POST, PATCH y DELETE.
- Cómo enviar y recibir JSON.
- Cómo conectar un backend con MySQL.
- Cómo utilizar middleware.
- Cómo validar información.
- Cómo aplicar reglas de negocio.
- Cómo conectar frontend y backend.
- Cómo desarrollar una aplicación web funcional.

## 18. Explicación para el profesor

Mi proyecto se denomina **RentCar**, un sistema web de alquiler de vehículos desarrollado para la asignatura Aplicaciones para el Servidor Web.

La finalidad del sistema es permitir que los clientes consulten vehículos disponibles, revisen sus precios y realicen reservas mediante una página web. El administrador podrá gestionar vehículos, clientes y reservas.

El backend estará desarrollado en Go y utilizará una API REST para procesar solicitudes HTTP. La información se almacenará en MySQL. El frontend estará desarrollado con HTML, CSS y JavaScript, y se comunicará con el servidor mediante solicitudes GET, POST, PATCH y DELETE.

El sistema incorporará validaciones, middleware, respuestas JSON y reglas de negocio, como el cálculo automático del alquiler y la verificación de disponibilidad.

En conclusión, RentCar permitirá aplicar de forma práctica los fundamentos del desarrollo backend, la comunicación entre cliente y servidor, el manejo de bases de datos y la creación de aplicaciones web funcionales.

## 19. Identidad del proyecto

**Nombre:** RentCar

**Subtítulo:** Sistema web de alquiler y gestión de vehículos

**Descripción corta:** Plataforma web que permite consultar vehículos, gestionar clientes y realizar reservas de alquiler mediante una API REST desarrollada en Go y una base de datos MySQL.

**Frase sugerida:** “Tu próximo viaje comienza aquí.”

## 20. Recomendación final

La primera versión debe concentrarse en:

1. Vehículos.
2. Clientes.
3. Reservas.
4. Cálculo de precios.
5. Disponibilidad.
6. API REST.
7. Página web.
8. Panel administrativo.

No es necesario comenzar con pagos reales, autenticación avanzada o inteligencia artificial. Es mejor desarrollar primero una versión funcional y después añadir mejoras.

La combinación recomendada es:

```text
Go + MySQL + HTML + CSS + JavaScript + Postman
```
