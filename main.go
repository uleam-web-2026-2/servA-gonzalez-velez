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
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/sesion"
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
			&sesion.Usuario{},
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
		&sesion.Usuario{},
	); err != nil {
		log.Fatal("no se pudo migrar: ", err)
	}

	reservas.Sembrar(db)

	if err := sesion.Sembrar(db); err != nil {
		log.Fatal("no se pudieron sembrar los usuarios: ", err)
	}

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

	tokens := sesion.NuevosTokens(cfg.ClaveFirma)
	login := &sesion.Manejador{DB: db, Tokens: tokens}
	r.Post("/login", login.Login) // pública: es la puerta para conseguir el token

	(&reservas.Manejador{Repo: reservas.NuevoRepositorioGorm(db)}).Rutas(r, tokens)

	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}
	log.Println("escuchando en el puerto", cfg.Puerto)
	log.Fatal(servidor.ListenAndServe())
}
