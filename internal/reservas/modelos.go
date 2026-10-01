package reservas

import "time"

// Sucursal es la agencia física donde se retiran y devuelven los vehículos.
type Sucursal struct {
	ID        uint   `gorm:"primaryKey" json:"ID"`
	Nombre    string `gorm:"not null" json:"Nombre"`
	Ciudad    string `gorm:"not null" json:"Ciudad"`
	Direccion string `gorm:"not null" json:"Direccion"`
	Activa    bool   `gorm:"not null;default:true" json:"Activa"`
}

// Usuario es el personal interno (administrador o agente) de una sucursal.
type Usuario struct {
	ID         uint     `gorm:"primaryKey" json:"ID"`
	SucursalID uint     `gorm:"not null;index" json:"SucursalID"`
	Sucursal   Sucursal `json:"Sucursal,omitempty"`
	Nombre     string   `gorm:"not null" json:"Nombre"`
	Correo     string   `gorm:"not null" json:"Correo"`
	Rol        string   `gorm:"not null" json:"Rol"` // administrador | agente
}

// Cliente es quien solicita el alquiler.
type Cliente struct {
	ID       uint   `gorm:"primaryKey" json:"ID"`
	Nombre   string `gorm:"not null" json:"Nombre"`
	Correo   string `gorm:"not null" json:"Correo"`
	Telefono string `json:"Telefono"`
	Cedula   string `gorm:"not null" json:"Cedula"`
}

// Vehiculo es el lado del uno: un vehículo tiene muchas reservas.
// Tipo permitido: auto | camioneta.
type Vehiculo struct {
	ID          uint      `gorm:"primaryKey" json:"ID"`
	SucursalID  uint      `gorm:"not null;index" json:"SucursalID"`
	Sucursal    Sucursal  `json:"Sucursal,omitempty"`
	Marca       string    `gorm:"not null" json:"Marca"`
	Modelo      string    `gorm:"not null" json:"Modelo"`
	Tipo        string    `gorm:"not null" json:"Tipo"` // auto | camioneta
	PrecioDia   float64   `gorm:"not null" json:"PrecioDia"`
	EstadoFlota string    `gorm:"not null;default:disponible" json:"EstadoFlota"` // disponible | alquilado | mantenimiento
	Reservas    []Reserva `json:"Reservas,omitempty"`
}

// Reserva es el lado de los muchos y la entidad con estados (CRUD del taller).
type Reserva struct {
	ID          uint     `gorm:"primaryKey" json:"ID"`
	ClienteID   uint     `gorm:"not null;index" json:"ClienteID"`
	Cliente     Cliente  `json:"Cliente,omitempty"`
	VehiculoID  uint     `gorm:"not null;index" json:"VehiculoID"`
	Vehiculo    Vehiculo `json:"Vehiculo,omitempty"`
	FechaInicio string   `gorm:"not null" json:"FechaInicio"` // YYYY-MM-DD
	FechaFin    string   `gorm:"not null" json:"FechaFin"`
	Total       float64  `gorm:"not null" json:"Total"`
	Estado      string   `gorm:"not null" json:"Estado"`
	Pagos       []Pago   `json:"Pagos,omitempty"`
}

// Pago es el comprobante de transferencia o efectivo ligado a una reserva.
type Pago struct {
	ID         uint      `gorm:"primaryKey" json:"ID"`
	ReservaID  uint      `gorm:"not null;index" json:"ReservaID"`
	Monto      float64   `gorm:"not null" json:"Monto"`
	Metodo     string    `gorm:"not null" json:"Metodo"`     // transferencia | efectivo
	EstadoPago string    `gorm:"not null" json:"EstadoPago"` // por_verificar | verificado | rechazado
	Referencia string    `json:"Referencia"`
	Creado     time.Time `json:"Creado"`
}

var estadosValidos = map[string]bool{
	"pendiente":  true,
	"confirmada": true,
	"en_curso":   true,
	"finalizada": true,
	"cancelada":  true,
}

var tiposValidos = map[string]bool{
	"auto":      true,
	"camioneta": true,
}

// transicionesPermite[de][a] = true
var transicionesPermite = map[string]map[string]bool{
	"pendiente": {
		"confirmada": true,
		"cancelada":  true,
	},
	"confirmada": {
		"en_curso":  true,
		"cancelada": true,
	},
	"en_curso": {
		"finalizada": true,
	},
}

func transicionPermitida(de, a string) bool {
	if de == a {
		return true
	}
	destino, ok := transicionesPermite[de]
	if !ok {
		return false
	}
	return destino[a]
}
