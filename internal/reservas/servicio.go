package reservas

import (
	"errors"
	"fmt"
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

type Servicio struct {
	repo Repositorio
}

func NuevoServicio(repo Repositorio) *Servicio {
	return &Servicio{repo: repo}
}

func (s *Servicio) repositorio() (Repositorio, error) {
	if s == nil || s.repo == nil {
		return nil, ErrRepositorioNoDisponible
	}
	return s.repo, nil
}

func (s *Servicio) Crear(reserva Reserva) (Reserva, error) {
	reserva.ID = 0
	reserva.Estado = "pendiente"
	if err := validarReserva(reserva); err != nil {
		return Reserva{}, fmt.Errorf("%w: %s", ErrDatosInvalidos, err)
	}

	repo, err := s.repositorio()
	if err != nil {
		return Reserva{}, err
	}
	clienteExiste, err := repo.ClienteExiste(reserva.ClienteID)
	if err != nil {
		return Reserva{}, err
	}
	if !clienteExiste {
		return Reserva{}, ErrClienteInexistente
	}
	vehiculoExiste, err := repo.VehiculoExiste(reserva.VehiculoID)
	if err != nil {
		return Reserva{}, err
	}
	if !vehiculoExiste {
		return Reserva{}, ErrVehiculoInexistente
	}
	if err := repo.CrearReserva(&reserva); err != nil {
		return Reserva{}, err
	}
	return reserva, nil
}

func (s *Servicio) Listar(estado string, limit, offset int) ([]Reserva, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarReservas(estado, limit, offset)
}

func (s *Servicio) Buscar(id uint) (Reserva, error) {
	repo, err := s.repositorio()
	if err != nil {
		return Reserva{}, err
	}
	return repo.BuscarReserva(id)
}

func (s *Servicio) Actualizar(id uint, cambios Reserva) (Reserva, error) {
	repo, err := s.repositorio()
	if err != nil {
		return Reserva{}, err
	}
	existente, err := repo.BuscarReserva(id)
	if err != nil {
		return Reserva{}, err
	}

	cambios.ID = existente.ID
	cambios.VehiculoID = existente.VehiculoID
	cambios.Estado = existente.Estado
	if err := validarReserva(cambios); err != nil {
		return Reserva{}, fmt.Errorf("%w: %s", ErrDatosInvalidos, err)
	}
	clienteExiste, err := repo.ClienteExiste(cambios.ClienteID)
	if err != nil {
		return Reserva{}, err
	}
	if !clienteExiste {
		return Reserva{}, ErrClienteInexistente
	}

	existente.ClienteID = cambios.ClienteID
	existente.FechaInicio = cambios.FechaInicio
	existente.FechaFin = cambios.FechaFin
	existente.Total = cambios.Total
	if err := repo.GuardarReserva(&existente); err != nil {
		return Reserva{}, err
	}
	return existente, nil
}

func (s *Servicio) CambiarEstado(id uint, nuevo string) (Reserva, error) {
	if !estadosValidos[nuevo] {
		return Reserva{}, fmt.Errorf("%w: %s", ErrEstadoInvalido, nuevo)
	}
	repo, err := s.repositorio()
	if err != nil {
		return Reserva{}, err
	}
	reserva, err := repo.BuscarReserva(id)
	if err != nil {
		return Reserva{}, err
	}
	if !puedePasar(reserva.Estado, nuevo) {
		return Reserva{}, fmt.Errorf("%w: de %s no se puede pasar a %s",
			ErrTransicionProhibida, reserva.Estado, nuevo)
	}
	reserva.Estado = nuevo
	if err := repo.GuardarReserva(&reserva); err != nil {
		return Reserva{}, err
	}
	return reserva, nil
}

func (s *Servicio) Borrar(id uint) (bool, error) {
	repo, err := s.repositorio()
	if err != nil {
		return false, err
	}
	return repo.BorrarReserva(id)
}

func (s *Servicio) ListarVehiculos() ([]Vehiculo, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarVehiculos()
}

func (s *Servicio) ListarSucursales() ([]Sucursal, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarSucursales()
}

func (s *Servicio) ListarClientes() ([]Cliente, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarClientes()
}

func (s *Servicio) ListarUsuarios() ([]Usuario, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarUsuarios()
}

func (s *Servicio) ListarPagos() ([]Pago, error) {
	repo, err := s.repositorio()
	if err != nil {
		return nil, err
	}
	return repo.ListarPagos()
}
