package reservas

import "encoding/json"

// ReservaPeticion es el DTO de entrada para crear o actualizar reservas.
// No expone ni permite fijar ID ni Estado desde el cliente HTTP.
type ReservaPeticion struct {
	ClienteID   uint    `json:"ClienteID"`
	VehiculoID  uint    `json:"VehiculoID"`
	FechaInicio string  `json:"FechaInicio"`
	FechaFin    string  `json:"FechaFin"`
	Total       float64 `json:"Total"`
}

// UnmarshalJSON admite tanto claves en PascalCase como snake_case.
func (p *ReservaPeticion) UnmarshalJSON(data []byte) error {
	var aux struct {
		ClienteID      uint    `json:"ClienteID"`
		ClienteIDAlt   uint    `json:"cliente_id"`
		VehiculoID     uint    `json:"VehiculoID"`
		VehiculoIDAlt  uint    `json:"vehiculo_id"`
		FechaInicio    string  `json:"FechaInicio"`
		FechaInicioAlt string  `json:"fecha_inicio"`
		FechaFin       string  `json:"FechaFin"`
		FechaFinAlt    string  `json:"fecha_fin"`
		Total          float64 `json:"Total"`
		TotalAlt       float64 `json:"total"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.ClienteID != 0 {
		p.ClienteID = aux.ClienteID
	} else {
		p.ClienteID = aux.ClienteIDAlt
	}

	if aux.VehiculoID != 0 {
		p.VehiculoID = aux.VehiculoID
	} else {
		p.VehiculoID = aux.VehiculoIDAlt
	}

	if aux.FechaInicio != "" {
		p.FechaInicio = aux.FechaInicio
	} else {
		p.FechaInicio = aux.FechaInicioAlt
	}

	if aux.FechaFin != "" {
		p.FechaFin = aux.FechaFin
	} else {
		p.FechaFin = aux.FechaFinAlt
	}

	if aux.Total != 0 {
		p.Total = aux.Total
	} else {
		p.Total = aux.TotalAlt
	}
	return nil
}

func (p ReservaPeticion) aEntidad() Reserva {
	return Reserva{
		ClienteID:   p.ClienteID,
		VehiculoID:  p.VehiculoID,
		FechaInicio: p.FechaInicio,
		FechaFin:    p.FechaFin,
		Total:       p.Total,
		Estado:      "pendiente",
	}
}

// ReservaRespuesta es el DTO de salida enviado a través de HTTP.
type ReservaRespuesta struct {
	ID          uint      `json:"ID"`
	ClienteID   uint      `json:"ClienteID"`
	VehiculoID  uint      `json:"VehiculoID"`
	FechaInicio string    `json:"FechaInicio"`
	FechaFin    string    `json:"FechaFin"`
	Total       float64   `json:"Total"`
	Estado      string    `json:"Estado"`
	Cliente     *Cliente  `json:"Cliente,omitempty"`
	Vehiculo    *Vehiculo `json:"Vehiculo,omitempty"`
	Pagos       []Pago    `json:"Pagos,omitempty"`
}

func NuevaRespuesta(r Reserva) ReservaRespuesta {
	var clientePtr *Cliente
	if r.Cliente.ID != 0 {
		c := r.Cliente
		clientePtr = &c
	}
	var vehiculoPtr *Vehiculo
	if r.Vehiculo.ID != 0 {
		v := r.Vehiculo
		vehiculoPtr = &v
	}
	return ReservaRespuesta{
		ID:          r.ID,
		ClienteID:   r.ClienteID,
		VehiculoID:  r.VehiculoID,
		FechaInicio: r.FechaInicio,
		FechaFin:    r.FechaFin,
		Total:       r.Total,
		Estado:      r.Estado,
		Cliente:     clientePtr,
		Vehiculo:    vehiculoPtr,
		Pagos:       r.Pagos,
	}
}

func NuevasRespuestas(lista []Reserva) []ReservaRespuesta {
	respuestas := make([]ReservaRespuesta, len(lista))
	for i, r := range lista {
		respuestas[i] = NuevaRespuesta(r)
	}
	return respuestas
}
