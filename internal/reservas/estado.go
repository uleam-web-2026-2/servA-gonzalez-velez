package reservas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
)

// cambiarEstado: PATCH /reservas/{id} con el cuerpo {"estado": "..."}.
// 400 id o JSON · 422 estado que no existe · 404 no existe
// · 422 transición prohibida · 200 cambiado.
func (m *Manejador) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		respuesta.Error(w, http.StatusBadRequest,
			"el id debe ser un número entero positivo")
		return
	}
	var pedido struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "el cuerpo no es JSON válido")
		return
	}
	reserva, err := m.Servicio.CambiarEstado(uint(id), pedido.Estado)
	switch {
	case errors.Is(err, ErrEstadoInvalido):
		mensaje := strings.TrimPrefix(err.Error(), ErrEstadoInvalido.Error()+": ")
		respuesta.Error(w, http.StatusUnprocessableEntity, mensaje)
		return
	case errors.Is(err, ErrNoEncontrado):
		respuesta.Error(w, http.StatusNotFound, "no existe la reserva")
		return
	case errors.Is(err, ErrTransicionProhibida):
		mensaje := strings.TrimPrefix(err.Error(), ErrTransicionProhibida.Error()+": ")
		respuesta.Error(w, http.StatusUnprocessableEntity, mensaje)
		return
	case err != nil:
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo cambiar el estado")
		return
	}
	respuesta.JSON(w, http.StatusOK, NuevaRespuesta(reserva))
}
