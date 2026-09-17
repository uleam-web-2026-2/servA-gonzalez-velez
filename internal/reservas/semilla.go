package reservas

import "gorm.io/gorm"

// Sembrar carga vehículos y reservas de ejemplo solo si la tabla está vacía.
func Sembrar(db *gorm.DB) {
	var total int64
	db.Model(&Vehiculo{}).Count(&total)
	if total > 0 {
		return
	}

	datos := []Vehiculo{
		{
			Marca:     "Toyota",
			Modelo:    "Corolla",
			Tipo:      "auto",
			PrecioDia: 35,
			Reservas: []Reserva{
				{
					Cliente:     "Ana Pérez",
					FechaInicio: "2026-09-20",
					FechaFin:    "2026-09-23",
					Total:       105,
					Estado:      "confirmada",
				},
				{
					Cliente:     "Luis Mora",
					FechaInicio: "2026-10-01",
					FechaFin:    "2026-10-03",
					Total:       70,
					Estado:      "pendiente",
				},
			},
		},
		{
			Marca:     "Yamaha",
			Modelo:    "MT-07",
			Tipo:      "moto",
			PrecioDia: 25,
			Reservas: []Reserva{
				{
					Cliente:     "Carla Ruiz",
					FechaInicio: "2026-09-18",
					FechaFin:    "2026-09-19",
					Total:       25,
					Estado:      "en_curso",
				},
			},
		},
		{
			Marca:     "Chevrolet",
			Modelo:    "Colorado",
			Tipo:      "camioneta",
			PrecioDia: 55,
		},
	}
	db.Create(&datos)
}
