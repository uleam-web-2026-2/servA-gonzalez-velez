package reservas

import (
	"time"

	"gorm.io/gorm"
)

// Sembrar carga datos de ejemplo solo si no hay sucursales.
func Sembrar(db *gorm.DB) {
	var total int64
	db.Model(&Sucursal{}).Count(&total)
	if total > 0 {
		return
	}

	sucursal := Sucursal{
		Nombre:    "RentCar Manta Centro",
		Ciudad:    "Manta",
		Direccion: "Av. 24 y Calle 13",
		Activa:    true,
	}
	db.Create(&sucursal)

	usuarios := []Usuario{
		{SucursalID: sucursal.ID, Nombre: "Jipson Velez", Correo: "jipson@rentcar.ec", Rol: "administrador"},
		{SucursalID: sucursal.ID, Nombre: "Isaac Gonzalez", Correo: "isaac@rentcar.ec", Rol: "agente"},
	}
	db.Create(&usuarios)

	clientes := []Cliente{
		{Nombre: "Ana Pérez", Correo: "ana@example.com", Telefono: "0991111111", Cedula: "1310000001"},
		{Nombre: "Luis Mora", Correo: "luis@example.com", Telefono: "0992222222", Cedula: "1310000002"},
		{Nombre: "Carla Ruiz", Correo: "carla@example.com", Telefono: "0993333333", Cedula: "1310000003"},
	}
	db.Create(&clientes)

	vehiculos := []Vehiculo{
		{
			SucursalID:  sucursal.ID,
			Marca:       "Toyota",
			Modelo:      "Corolla",
			Tipo:        "auto",
			PrecioDia:   35,
			EstadoFlota: "disponible",
		},
		{
			SucursalID:  sucursal.ID,
			Marca:       "Kia",
			Modelo:      "Rio",
			Tipo:        "auto",
			PrecioDia:   28,
			EstadoFlota: "disponible",
		},
		{
			SucursalID:  sucursal.ID,
			Marca:       "Chevrolet",
			Modelo:      "Colorado",
			Tipo:        "camioneta",
			PrecioDia:   55,
			EstadoFlota: "disponible",
		},
	}
	db.Create(&vehiculos)

	reservas := []Reserva{
		{
			ClienteID:   clientes[0].ID,
			VehiculoID:  vehiculos[0].ID,
			FechaInicio: "2026-09-20",
			FechaFin:    "2026-09-23",
			Total:       105,
			Estado:      "confirmada",
		},
		{
			ClienteID:   clientes[1].ID,
			VehiculoID:  vehiculos[0].ID,
			FechaInicio: "2026-10-01",
			FechaFin:    "2026-10-03",
			Total:       70,
			Estado:      "pendiente",
		},
		{
			ClienteID:   clientes[2].ID,
			VehiculoID:  vehiculos[1].ID,
			FechaInicio: "2026-09-18",
			FechaFin:    "2026-09-19",
			Total:       28,
			Estado:      "en_curso",
		},
	}
	db.Create(&reservas)

	pagos := []Pago{
		{
			ReservaID:  reservas[0].ID,
			Monto:      105,
			Metodo:     "transferencia",
			EstadoPago: "verificado",
			Referencia: "TRX-1001",
			Creado:     time.Date(2026, 9, 19, 10, 0, 0, 0, time.UTC),
		},
		{
			ReservaID:  reservas[1].ID,
			Monto:      70,
			Metodo:     "transferencia",
			EstadoPago: "por_verificar",
			Referencia: "TRX-1002",
			Creado:     time.Date(2026, 9, 30, 15, 30, 0, 0, time.UTC),
		},
	}
	db.Create(&pagos)
}
