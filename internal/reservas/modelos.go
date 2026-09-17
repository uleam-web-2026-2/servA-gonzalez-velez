package reservas

// Vehiculo es el lado del uno: un vehículo tiene muchas reservas.
type Vehiculo struct {
	ID        uint      `json:"ID"`
	Marca     string    `json:"Marca"`
	Modelo    string    `json:"Modelo"`
	Tipo      string    `json:"Tipo"` // moto | auto | camioneta
	PrecioDia float64   `json:"PrecioDia"`
	Reservas  []Reserva `json:"Reservas,omitempty"`
}

// Reserva es el lado de los muchos y la entidad con estados (CRUD del taller).
type Reserva struct {
	ID          uint    `json:"ID"`
	VehiculoID  uint    `json:"VehiculoID"`
	Cliente     string  `json:"Cliente"`
	FechaInicio string  `json:"FechaInicio"` // YYYY-MM-DD
	FechaFin    string  `json:"FechaFin"`
	Total       float64 `json:"Total"`
	Estado      string  `json:"Estado"`
}

// estadosValidos es la lista cerrada de estados de una reserva (kit S3 §5).
var estadosValidos = map[string]bool{
	"pendiente":   true,
	"confirmada":  true,
	"en_curso":    true,
	"finalizada":  true,
	"cancelada":   true,
}
