package reservas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
)

// Manejador guarda la conexión. Todas las rutas la leen con m.DB.
type Manejador struct{ DB *gorm.DB }

// Rutas registra el CRUD de reservas y los listados del modelo.
func (m *Manejador) Rutas(r chi.Router) {
	r.Post("/reservas", m.crear)
	r.Get("/reservas", m.listar)
	r.Get("/reservas/{id}", m.verUno)
	r.Put("/reservas/{id}", m.actualizar)
	r.Patch("/reservas/{id}", m.cambiarEstado)
	r.Delete("/reservas/{id}", m.borrar)

	r.Get("/vehiculos", m.listarVehiculos)
	r.Get("/sucursales", m.listarSucursales)
	r.Get("/clientes", m.listarClientes)
	r.Get("/usuarios", m.listarUsuarios)
	r.Get("/pagos", m.listarPagos)
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

	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo validar")
		return
	}

	var cliente Cliente
	if err := m.DB.First(&cliente, reserva.ClienteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusUnprocessableEntity,
				"cliente_inexistente", "El cliente no existe")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo validar el cliente")
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

	limit, offset, ok := leerPaginacion(w, r)
	if !ok {
		return
	}

	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}

	var lista []Reserva
	q := m.DB.Debug()
	if estado != "" {
		q = q.Where("estado = ?", estado)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
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
	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
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
	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
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
	if !transicionPermitida(existente.Estado, entrada.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"transicion_prohibida",
			"No se permite pasar de "+existente.Estado+" a "+entrada.Estado)
		return
	}

	entrada.ID = existente.ID
	entrada.VehiculoID = existente.VehiculoID // no se cambia la FK al actualizar
	if err := validarReserva(entrada); err != nil {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"datos_invalidos", err.Error())
		return
	}

	var cliente Cliente
	if err := m.DB.First(&cliente, entrada.ClienteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusUnprocessableEntity,
				"cliente_inexistente", "El cliente no existe")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo validar el cliente")
		return
	}

	existente.ClienteID = entrada.ClienteID
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

func (m *Manejador) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	var cuerpo struct {
		Estado string `json:"Estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cuerpo); err != nil {
		respuesta.Error(w, http.StatusBadRequest,
			"json_invalido", "El cuerpo no es un JSON válido")
		return
	}
	if !estadosValidos[cuerpo.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"estado_invalido", "Estado no válido")
		return
	}

	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo consultar")
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

	if !transicionPermitida(existente.Estado, cuerpo.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity,
			"transicion_prohibida",
			"No se permite pasar de "+existente.Estado+" a "+cuerpo.Estado)
		return
	}

	existente.Estado = cuerpo.Estado
	if err := m.DB.Debug().Save(&existente).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo actualizar el estado")
		return
	}
	respuesta.Exito(w, http.StatusOK, existente)
}

func (m *Manejador) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}
	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo borrar")
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

func (m *Manejador) listarVehiculos(w http.ResponseWriter, r *http.Request) {
	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar vehículos")
		return
	}
	var lista []Vehiculo
	if err := m.DB.Debug().Preload("Reservas").Find(&lista).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar vehículos")
		return
	}
	respuesta.Exito(w, http.StatusOK, lista)
}

func (m *Manejador) listarSucursales(w http.ResponseWriter, r *http.Request) {
	m.listarEntidad(w, &[]Sucursal{})
}

func (m *Manejador) listarClientes(w http.ResponseWriter, r *http.Request) {
	m.listarEntidad(w, &[]Cliente{})
}

func (m *Manejador) listarUsuarios(w http.ResponseWriter, r *http.Request) {
	m.listarEntidad(w, &[]Usuario{})
}

func (m *Manejador) listarPagos(w http.ResponseWriter, r *http.Request) {
	m.listarEntidad(w, &[]Pago{})
}

func (m *Manejador) listarEntidad(w http.ResponseWriter, dest any) {
	if m.DB == nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}
	if err := m.DB.Debug().Find(dest).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError,
			"error_base", "No se pudo listar")
		return
	}
	respuesta.Exito(w, http.StatusOK, dest)
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

// validarReserva aplica la regla extra de negocio (fechas y total).
func validarReserva(r Reserva) error {
	if r.ClienteID == 0 {
		return errors.New("Debe indicar un ClienteID válido")
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
