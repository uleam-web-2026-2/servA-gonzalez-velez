package reservas

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNoEncontrado = errors.New("registro no encontrado")

type Repositorio interface {
	ClienteExiste(id uint) (bool, error)
	VehiculoExiste(id uint) (bool, error)
	CrearReserva(reserva *Reserva) error
	ListarReservas(estado string, limit, offset int) ([]Reserva, error)
	BuscarReserva(id uint) (Reserva, error)
	GuardarReserva(reserva *Reserva) error
	BorrarReserva(id uint) (bool, error)
	ListarVehiculos() ([]Vehiculo, error)
	ListarSucursales() ([]Sucursal, error)
	ListarClientes() ([]Cliente, error)
	ListarUsuarios() ([]Usuario, error)
	ListarPagos() ([]Pago, error)
}

type RepositorioGorm struct {
	db *gorm.DB
}

func NuevoRepositorioGorm(db *gorm.DB) *RepositorioGorm {
	return &RepositorioGorm{db: db}
}

func (repo *RepositorioGorm) ClienteExiste(id uint) (bool, error) {
	var cliente Cliente
	err := repo.db.First(&cliente, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (repo *RepositorioGorm) VehiculoExiste(id uint) (bool, error) {
	var vehiculo Vehiculo
	err := repo.db.First(&vehiculo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (repo *RepositorioGorm) CrearReserva(reserva *Reserva) error {
	return repo.db.Debug().Create(reserva).Error
}

func (repo *RepositorioGorm) ListarReservas(estado string, limit, offset int) ([]Reserva, error) {
	var lista []Reserva
	query := repo.db.Debug()
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

func (repo *RepositorioGorm) BuscarReserva(id uint) (Reserva, error) {
	var reserva Reserva
	err := repo.db.First(&reserva, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Reserva{}, ErrNoEncontrado
	}
	return reserva, err
}

func (repo *RepositorioGorm) GuardarReserva(reserva *Reserva) error {
	return repo.db.Debug().Save(reserva).Error
}

func (repo *RepositorioGorm) BorrarReserva(id uint) (bool, error) {
	resultado := repo.db.Debug().Delete(&Reserva{}, id)
	return resultado.RowsAffected > 0, resultado.Error
}

func (repo *RepositorioGorm) ListarVehiculos() ([]Vehiculo, error) {
	var lista []Vehiculo
	err := repo.db.Debug().Preload("Reservas").Find(&lista).Error
	return lista, err
}

func (repo *RepositorioGorm) ListarSucursales() ([]Sucursal, error) {
	var lista []Sucursal
	err := repo.db.Debug().Find(&lista).Error
	return lista, err
}

func (repo *RepositorioGorm) ListarClientes() ([]Cliente, error) {
	var lista []Cliente
	err := repo.db.Debug().Find(&lista).Error
	return lista, err
}

func (repo *RepositorioGorm) ListarUsuarios() ([]Usuario, error) {
	var lista []Usuario
	err := repo.db.Debug().Find(&lista).Error
	return lista, err
}

func (repo *RepositorioGorm) ListarPagos() ([]Pago, error) {
	var lista []Pago
	err := repo.db.Debug().Find(&lista).Error
	return lista, err
}
