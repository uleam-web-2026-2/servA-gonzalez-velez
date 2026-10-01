package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/config"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/middleware"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/reservas"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
)

func main() {
	reset := flag.Bool("reset", false, "borra las tablas y arranca con la base vacía")
	flag.Parse()

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal("configuración: ", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar: ", err)
	}

	if *reset {
		db.Migrator().DropTable(
			&reservas.Pago{},
			&reservas.Reserva{},
			&reservas.Vehiculo{},
			&reservas.Cliente{},
			&reservas.Usuario{},
			&reservas.Sucursal{},
		)
	}

	if err := db.Debug().AutoMigrate(
		&reservas.Sucursal{},
		&reservas.Usuario{},
		&reservas.Cliente{},
		&reservas.Vehiculo{},
		&reservas.Reserva{},
		&reservas.Pago{},
	); err != nil {
		log.Fatal("no se pudo migrar: ", err)
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

	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}
	log.Println("escuchando en el puerto", cfg.Puerto)
	log.Fatal(servidor.ListenAndServe())
}
