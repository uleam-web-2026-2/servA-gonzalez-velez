package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/middleware"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/reservas"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
)

func main() {
	reset := flag.Bool("reset", false, "borra las tablas y arranca con la base vacía")
	flag.Parse()

	// Misma cadena del lab (Docker pg en 5433). Semana 4 se externaliza.
	dsn := "host=localhost port=5433 user=postgres password=taller2026 dbname=rentcar"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	if *reset {
		// primero la tabla del lado de los muchos, después la del uno
		db.Migrator().DropTable(&reservas.Reserva{}, &reservas.Vehiculo{})
	}

	err = db.Debug().AutoMigrate(&reservas.Vehiculo{}, &reservas.Reserva{})
	if err != nil {
		log.Fatal(err)
	}

	reservas.Sembrar(db)

	r := chi.NewRouter()
	r.Use(middleware.Registro)
	r.Use(middleware.Recuperacion)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		respuesta.Error(w, http.StatusNotFound, "ruta_inexistente",
			"la ruta no existe")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		respuesta.Error(w, http.StatusMethodNotAllowed, "metodo_no_permitido",
			"el método no está permitido en esta ruta")
	})

	(&reservas.Manejador{DB: db}).Rutas(r)

	log.Println("RentCar en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
