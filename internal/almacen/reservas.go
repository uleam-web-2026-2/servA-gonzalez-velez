package almacen

import (
	"errors"

	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/reservas"
)

type Reservas struct {
	DB *gorm.DB
}

var _ reservas.Repositorio = (*Reservas)(nil)

func NuevoReservas(db *gorm.DB) *Reservas {
	return &Reservas{DB: db}
}

func (a *Reservas) ClienteExiste(id uint) (bool, error) {
	var cliente reservas.Cliente
	err := a.DB.First(&cliente, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (a *Reservas) VehiculoExiste(id uint) (bool, error) {
	var vehiculo reservas.Vehiculo
	err := a.DB.First(&vehiculo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (a *Reservas) CrearReserva(reserva *reservas.Reserva) error {
	return a.DB.Debug().Create(reserva).Error
}

func (a *Reservas) ListarReservas(estado string, limit, offset int) ([]reservas.Reserva, error) {
	var lista []reservas.Reserva
	query := a.DB.Debug()
	if estado != "" {
		query = query.Where("estado = ?", estado)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&lista).Error
	return lista, err
}

func (a *Reservas) BuscarReserva(id uint) (reservas.Reserva, error) {
	var reserva reservas.Reserva
	err := a.DB.First(&reserva, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return reservas.Reserva{}, reservas.ErrNoEncontrado
	}
	return reserva, err
}

func (a *Reservas) GuardarReserva(reserva *reservas.Reserva) error {
	return a.DB.Debug().Save(reserva).Error
}

func (a *Reservas) BorrarReserva(id uint) (bool, error) {
	resultado := a.DB.Debug().Delete(&reservas.Reserva{}, id)
	return resultado.RowsAffected > 0, resultado.Error
}

func (a *Reservas) ListarVehiculos() ([]reservas.Vehiculo, error) {
	var lista []reservas.Vehiculo
	err := a.DB.Debug().Preload("Reservas").Find(&lista).Error
	return lista, err
}

func (a *Reservas) ListarSucursales() ([]reservas.Sucursal, error) {
	var lista []reservas.Sucursal
	err := a.DB.Debug().Find(&lista).Error
	return lista, err
}

func (a *Reservas) ListarClientes() ([]reservas.Cliente, error) {
	var lista []reservas.Cliente
	err := a.DB.Debug().Find(&lista).Error
	return lista, err
}

func (a *Reservas) ListarUsuarios() ([]reservas.Usuario, error) {
	var lista []reservas.Usuario
	err := a.DB.Debug().Find(&lista).Error
	return lista, err
}

func (a *Reservas) ListarPagos() ([]reservas.Pago, error) {
	var lista []reservas.Pago
	err := a.DB.Debug().Find(&lista).Error
	return lista, err
}
