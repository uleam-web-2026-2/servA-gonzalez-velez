package reservas

import (
	"errors"
	"slices"
	"time"
)

var siguientes = map[string][]string{
	"pendiente":  {"confirmada", "cancelada"},
	"confirmada": {"en_curso", "cancelada"},
	"en_curso":   {"finalizada"},
	"finalizada": {},
	"cancelada":  {},
}

func puedePasar(de, a string) bool {
	return de == a || slices.Contains(siguientes[de], a)
}

func validarReserva(reserva Reserva) error {
	if reserva.ClienteID == 0 {
		return errors.New("Debe indicar un ClienteID válido")
	}
	if reserva.VehiculoID == 0 {
		return errors.New("Debe indicar un VehiculoID válido")
	}
	inicio, errInicio := time.Parse("2006-01-02", reserva.FechaInicio)
	fin, errFin := time.Parse("2006-01-02", reserva.FechaFin)
	if errInicio != nil || errFin != nil {
		return errors.New("Las fechas deben tener formato YYYY-MM-DD")
	}
	if !fin.After(inicio) {
		return errors.New("La fecha de fin debe ser posterior a la de inicio")
	}
	if reserva.Total <= 0 {
		return errors.New("El total debe ser mayor que cero")
	}
	return nil
}
