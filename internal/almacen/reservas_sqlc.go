package almacen

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/almacen/filas"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/reservas"
)

type ReservasSqlc struct {
	*Reservas
	Conexion *sql.DB
	q        *filas.Queries
}

var _ reservas.Repositorio = (*ReservasSqlc)(nil)

func NuevoReservasSqlc(g *Reservas, c *sql.DB) *ReservasSqlc {
	return &ReservasSqlc{Reservas: g, Conexion: c, q: filas.New(c)}
}

func (a *ReservasSqlc) BuscarReserva(id uint) (reservas.Reserva, error) {
	f, err := a.q.BuscarReserva(context.Background(), int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return reservas.Reserva{}, reservas.ErrNoEncontrado
	}
	if err != nil {
		return reservas.Reserva{}, err
	}
	return reservas.Reserva{
		ID:          uint(f.ID),
		ClienteID:   uint(f.ClienteID),
		VehiculoID:  uint(f.VehiculoID),
		FechaInicio: f.FechaInicio,
		FechaFin:    f.FechaFin,
		Total:       f.Total,
		Estado:      f.Estado,
	}, nil
}

func (a *ReservasSqlc) GuardarReserva(r *reservas.Reserva) error {
	n, err := a.q.GuardarReserva(context.Background(), filas.GuardarReservaParams{
		ID:          int64(r.ID),
		ClienteID:   int64(r.ClienteID),
		VehiculoID:  int64(r.VehiculoID),
		FechaInicio: r.FechaInicio,
		FechaFin:    r.FechaFin,
		Total:       r.Total,
		Estado:      r.Estado,
	})
	if err == nil && n == 0 {
		return reservas.ErrNoEncontrado
	}
	return err
}
