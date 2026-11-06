package reservas

import "testing"

type repositorioEnMemoria struct {
	clientes  map[uint]bool
	vehiculos map[uint]bool
	reservas  map[uint]Reserva
	siguiente uint
}

var _ Repositorio = (*repositorioEnMemoria)(nil)

func (r *repositorioEnMemoria) ClienteExiste(id uint) (bool, error) {
	return r.clientes[id], nil
}

func (r *repositorioEnMemoria) VehiculoExiste(id uint) (bool, error) {
	return r.vehiculos[id], nil
}

func (r *repositorioEnMemoria) CrearReserva(reserva *Reserva) error {
	r.siguiente++
	reserva.ID = r.siguiente
	r.reservas[reserva.ID] = *reserva
	return nil
}

func (r *repositorioEnMemoria) ListarReservas(_ string, _, _ int) ([]Reserva, error) {
	lista := make([]Reserva, 0, len(r.reservas))
	for _, reserva := range r.reservas {
		lista = append(lista, reserva)
	}
	return lista, nil
}

func (r *repositorioEnMemoria) BuscarReserva(id uint) (Reserva, error) {
	reserva, existe := r.reservas[id]
	if !existe {
		return Reserva{}, ErrNoEncontrado
	}
	return reserva, nil
}

func (r *repositorioEnMemoria) GuardarReserva(reserva *Reserva) error {
	r.reservas[reserva.ID] = *reserva
	return nil
}

func (r *repositorioEnMemoria) BorrarReserva(id uint) (bool, error) {
	if _, existe := r.reservas[id]; !existe {
		return false, nil
	}
	delete(r.reservas, id)
	return true, nil
}

func (r *repositorioEnMemoria) ListarVehiculos() ([]Vehiculo, error) {
	return nil, nil
}

func (r *repositorioEnMemoria) ListarSucursales() ([]Sucursal, error) {
	return nil, nil
}

func (r *repositorioEnMemoria) ListarClientes() ([]Cliente, error) {
	return nil, nil
}

func (r *repositorioEnMemoria) ListarUsuarios() ([]Usuario, error) {
	return nil, nil
}

func (r *repositorioEnMemoria) ListarPagos() ([]Pago, error) {
	return nil, nil
}

func TestServicioCrearFijaIDYEstado(t *testing.T) {
	repo := &repositorioEnMemoria{
		clientes:  map[uint]bool{1: true},
		vehiculos: map[uint]bool{1: true},
		reservas:  make(map[uint]Reserva),
		siguiente: 10,
	}
	servicio := NuevoServicio(repo)

	creada, err := servicio.Crear(Reserva{
		ID:          77,
		ClienteID:   1,
		VehiculoID:  1,
		FechaInicio: "2026-11-01",
		FechaFin:    "2026-11-04",
		Total:       105,
		Estado:      "confirmada",
	})
	if err != nil {
		t.Fatalf("Crear devolvió un error: %v", err)
	}
	if creada.ID != 11 || creada.Estado != "pendiente" {
		t.Fatalf("Crear devolvió ID=%d y Estado=%q; se esperaba ID=11 y Estado=pendiente",
			creada.ID, creada.Estado)
	}
}

func TestServicioActualizarPreservaVehiculoYEstado(t *testing.T) {
	repo := &repositorioEnMemoria{
		clientes:  map[uint]bool{1: true, 2: true},
		vehiculos: map[uint]bool{1: true},
		reservas: map[uint]Reserva{
			9: {
				ID:          9,
				ClienteID:   1,
				VehiculoID:  1,
				FechaInicio: "2026-11-01",
				FechaFin:    "2026-11-04",
				Total:       105,
				Estado:      "confirmada",
			},
		},
	}
	servicio := NuevoServicio(repo)

	actualizada, err := servicio.Actualizar(9, Reserva{
		ClienteID:   2,
		VehiculoID:  99,
		FechaInicio: "2026-11-05",
		FechaFin:    "2026-11-08",
		Total:       120,
		Estado:      "finalizada",
	})
	if err != nil {
		t.Fatalf("Actualizar devolvió un error: %v", err)
	}
	if actualizada.VehiculoID != 1 || actualizada.Estado != "confirmada" {
		t.Fatalf("Actualizar alteró campos protegidos: VehiculoID=%d Estado=%q",
			actualizada.VehiculoID, actualizada.Estado)
	}
	if actualizada.ClienteID != 2 || actualizada.FechaInicio != "2026-11-05" {
		t.Fatalf("Actualizar no guardó los cambios permitidos: %+v", actualizada)
	}
}
