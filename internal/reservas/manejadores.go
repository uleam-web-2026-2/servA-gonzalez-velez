package reservas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
)

// Manejador guarda la conexión. Todas las rutas la leen con m.DB.
type Manejador struct{ DB *gorm.DB }

// Rutas registra las rutas de reservas y el listado de vehículos con Preload.
func (m *Manejador) Rutas(r chi.Router) {
	r.Post("/reservas", m.crear)
	r.Get("/reservas", m.listar)
	r.Get("/reservas/{id}", m.verUno)
	r.Put("/reservas/{id}", m.actualizar)
	r.Delete("/reservas/{id}", m.borrar)
	r.Get("/vehiculos", m.listarVehiculos)
}

func (m *Manejador) crear(w http.ResponseWriter, r *http.Request) {
	var reserva Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		respuesta.Error(w, http.StatusBadRequest,
			"json_invalido", "El cuerpo no es un JSON válido")
		return
	}
	reserva.ID = 0

	if !estadosValidos[reserva.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"estado_invalido", "Estado no válido")
		return
	}

	if err := validarReserva(reserva); err != nil {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"datos_invalidos", err.Error())
		return
	}

	var vehiculo Vehiculo
	if err := m.DB.First(&vehiculo, reserva.VehiculoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusUnprocessableEntity,
				"vehiculo_inexistente", "El vehículo no existe")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo validar el vehículo")
		return
	}

	if err := m.DB.Debug().Create(&reserva).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo guardar")
		return
	}
	respuesta.Exito(w, http.StatusCreated, reserva)
}

func (m *Manejador) listar(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")
	var lista []Reserva
	q := m.DB.Debug()
	if estado != "" {
		q = q.Where("estado = ?", estado)
	}
	if err := q.Find(&lista).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) verUno(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	var reserva Reserva
	err := m.DB.First(&reserva, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
		return
	}
	respuesta.Exito(w, http.StatusOK, reserva)
}

func (m *Manejador) actualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	var existente Reserva
	err := m.DB.First(&existente, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
		return
	}

	var entrada Reserva
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest,
			"json_invalido", "El cuerpo no es un JSON válido")
		return
	}

	if !estadosValidos[entrada.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"estado_invalido", "Estado no válido")
		return
	}

	entrada.ID = existente.ID
	entrada.VehiculoID = existente.VehiculoID // no se cambia la FK al actualizar
	if err := validarReserva(entrada); err != nil {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"datos_invalidos", err.Error())
		return
	}

	existente.Cliente = entrada.Cliente
	existente.FechaInicio = entrada.FechaInicio
	existente.FechaFin = entrada.FechaFin
	existente.Total = entrada.Total
	existente.Estado = entrada.Estado

	if err := m.DB.Debug().Save(&existente).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo actualizar")
		return
	}
	respuesta.Exito(w, http.StatusOK, existente)
}

func (m *Manejador) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	res := m.DB.Debug().Delete(&Reserva{}, id)
	if res.Error != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo borrar")
		return
	}
	if res.RowsAffected == 0 {
		respuesta.Error(w, http.StatusNotFound,
			"no_encontrado", "No existe una reserva con ese id")
		return
	}
	respuesta.Exito(w, http.StatusOK, map[string]string{
		"mensaje": "reserva eliminada",
	})
}

// listarVehiculos lista el lado del uno con sus reservas (Preload → 2 consultas).
func (m *Manejador) listarVehiculos(w http.ResponseWriter, r *http.Request) {
	var lista []Vehiculo
	if err := m.DB.Debug().Preload("Reservas").Find(&lista).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar vehículos")
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

// validarReserva aplica la regla extra de negocio (fase 2b).
func validarReserva(r Reserva) error {
	if strings.TrimSpace(r.Cliente) == "" {
		return errors.New("El cliente no puede estar vacío")
	}
	if r.VehiculoID == 0 {
		return errors.New("Debe indicar un VehiculoID válido")
	}
	inicio, err1 := time.Parse("2006-01-02", r.FechaInicio)
	fin, err2 := time.Parse("2006-01-02", r.FechaFin)
	if err1 != nil || err2 != nil {
		return errors.New("Las fechas deben tener formato YYYY-MM-DD")
	}
	if !fin.After(inicio) {
		return errors.New("La fecha de fin debe ser posterior a la de inicio")
	}
	if r.Total <= 0 {
		return errors.New("El total debe ser mayor que cero")
	}
	return nil
}
