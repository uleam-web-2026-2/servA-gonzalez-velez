package reservas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/sesion"
)

type Manejador struct{ Servicio *Servicio }

// Rutas registra el CRUD de reservas y los listados del modelo agrupados por nivel.
func (m *Manejador) Rutas(r chi.Router, tokens *sesion.Tokens) {
	// Nivel 1 · Pública: cualquiera, sin token.
	r.Get("/vehiculos", m.listarVehiculos)

	// Nivel 2 · Con sesión: antes del manejador corre el guardia.
	conSesion := r.With(tokens.Autenticar)
	conSesion.Get("/reservas", m.listar)
	conSesion.Post("/reservas", m.crear)
	conSesion.Get("/reservas/{id}", m.verUno)
	conSesion.Get("/sucursales", m.listarSucursales)
	conSesion.Get("/clientes", m.listarClientes)
	conSesion.Get("/usuarios", m.listarUsuarios)
	conSesion.Get("/pagos", m.listarPagos)

	// Nivel 3 · Solo agentes: después del guardia se revisa el rol.
	soloAgente := conSesion.With(sesion.SoloRol(sesion.RolAgente))
	soloAgente.Put("/reservas/{id}", m.actualizar)
	soloAgente.Patch("/reservas/{id}", m.cambiarEstado)
	soloAgente.Delete("/reservas/{id}", m.borrar)
}

func (m *Manejador) crear(w http.ResponseWriter, r *http.Request) {
	var pedido ReservaPeticion
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest,
			"json_invalido", "El cuerpo no es un JSON válido")
		return
	}

	reserva, err := m.Servicio.Crear(pedido.aEntidad())
	switch {
	case errors.Is(err, ErrDatosInvalidos):
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"datos_invalidos",
			strings.TrimPrefix(err.Error(), ErrDatosInvalidos.Error()+": "))
		return
	case errors.Is(err, ErrClienteInexistente):
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"cliente_inexistente", "El cliente no existe")
		return
	case errors.Is(err, ErrVehiculoInexistente):
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"vehiculo_inexistente", "El vehículo no existe")
		return
	case err != nil:
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo guardar la reserva")
		return
	}
	respuesta.Exito(w, http.StatusCreated, NuevaRespuesta(reserva))
}

func (m *Manejador) listar(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")

	limit, offset, ok := leerPaginacion(w, r)
	if !ok {
		return
	}

	lista, err := m.Servicio.Listar(estado, limit, offset)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, NuevasRespuestas(lista))
}

func (m *Manejador) verUno(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	reserva, err := m.Servicio.Buscar(id)
	if errors.Is(err, ErrNoEncontrado) {
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
		return
	}
	respuesta.Exito(w, http.StatusOK, NuevaRespuesta(reserva))
}

func (m *Manejador) actualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	var pedido ReservaPeticion
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest,
			"json_invalido", "El cuerpo no es un JSON válido")
		return
	}

	actualizada, err := m.Servicio.Actualizar(id, pedido.aEntidad())
	switch {
	case errors.Is(err, ErrNoEncontrado):
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	case errors.Is(err, ErrDatosInvalidos):
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"datos_invalidos",
			strings.TrimPrefix(err.Error(), ErrDatosInvalidos.Error()+": "))
		return
	case errors.Is(err, ErrClienteInexistente):
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"cliente_inexistente", "El cliente no existe")
		return
	case err != nil:
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo actualizar")
		return
	}
	respuesta.Exito(w, http.StatusOK, NuevaRespuesta(actualizada))
}

func (m *Manejador) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	eliminada, err := m.Servicio.Borrar(id)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo borrar")
		return
	}
	if !eliminada {
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	}
	respuesta.Exito(w, http.StatusOK, map[string]string{
		"mensaje": "reserva eliminada",
	})
}

func (m *Manejador) listarVehiculos(w http.ResponseWriter, r *http.Request) {
	lista, err := m.Servicio.ListarVehiculos()
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar vehículos")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) listarSucursales(w http.ResponseWriter, r *http.Request) {
	m.listarSucursalesRespuesta(w)
}

func (m *Manejador) listarClientes(w http.ResponseWriter, r *http.Request) {
	m.listarClientesRespuesta(w)
}

func (m *Manejador) listarUsuarios(w http.ResponseWriter, r *http.Request) {
	m.listarUsuariosRespuesta(w)
}

func (m *Manejador) listarPagos(w http.ResponseWriter, r *http.Request) {
	m.listarPagosRespuesta(w)
}

func (m *Manejador) listarSucursalesRespuesta(w http.ResponseWriter) {
	lista, err := m.Servicio.ListarSucursales()
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) listarClientesRespuesta(w http.ResponseWriter) {
	lista, err := m.Servicio.ListarClientes()
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) listarUsuariosRespuesta(w http.ResponseWriter) {
	lista, err := m.Servicio.ListarUsuarios()
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) listarPagosRespuesta(w http.ResponseWriter) {
	lista, err := m.Servicio.ListarPagos()
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func leerID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	n, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || n <= 0 {
		respuesta.Error(w, http.StatusBadRequest,
			"id_invalido", "El id debe ser un número positivo")
		return 0, false
	}
	return uint(n), true
}

func leerPaginacion(w http.ResponseWriter, r *http.Request) (limit, offset int, ok bool) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n < 0 {
			respuesta.Error(w, http.StatusBadRequest,
				"paginacion_invalida", "limit debe ser un número entero no negativo")
			return 0, 0, false
		}
		limit = n
	}
	if offsetStr != "" {
		n, err := strconv.Atoi(offsetStr)
		if err != nil || n < 0 {
			respuesta.Error(w, http.StatusBadRequest,
				"paginacion_invalida", "offset debe ser un número entero no negativo")
			return 0, 0, false
		}
		offset = n
	}
	return limit, offset, true
}
