package reservas

import "errors"

var ErrNoEncontrado = errors.New("registro no encontrado")

// Repositorio es la interfaz que define las operaciones de persistencia
// que necesita el dominio de reservas. No depende de ningún ORM ni base de datos concreta.
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
