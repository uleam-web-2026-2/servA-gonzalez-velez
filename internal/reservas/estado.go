package reservas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
	if !estadosValidos[pedido.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"estado no permitido: "+pedido.Estado)
		return
	}
	if m.Repo == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"no se pudo leer")
		return
	}
	reserva, err := m.Repo.BuscarReserva(uint(id))
	if errors.Is(err, ErrNoEncontrado) {
		respuesta.Error(w, http.StatusNotFound, "no existe la reserva")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo leer")
		return
	}
	// La regla del negocio: se compara el estado que TIENE con el que PIDEN.
	if !puedePasar(reserva.Estado, pedido.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"de "+reserva.Estado+" no se puede pasar a "+pedido.Estado)
		return
	}
	reserva.Estado = pedido.Estado
	if err := m.Repo.GuardarReserva(&reserva); err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo guardar")
		return
	}
	respuesta.JSON(w, http.StatusOK, reserva)
}
