package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/almacen"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/config"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/middleware"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/reservas"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/sesion"
)

func main() {
	reset := flag.Bool("reset", false, "borra las tablas y arranca con la base vacía (requiere usuario migrador)")
	migrar := flag.Bool("migrar", false, "ejecuta migraciones, siembra y permisos mínimos con DATABASE_URL_MIGRACION y termina")
	flag.Parse()

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal("configuración: ", err)
	}

	if *migrar || *reset {
		dsnMigracion := os.Getenv("DATABASE_URL_MIGRACION")
		if dsnMigracion == "" {
			dsnMigracion = cfg.DatabaseURL
		}
		dbMig, err := gorm.Open(postgres.Open(dsnMigracion), &gorm.Config{})
		if err != nil {
			log.Fatal("no se pudo conectar para migrar: ", err)
		}

		if *reset {
			dbMig.Migrator().DropTable(
				&sesion.Usuario{},
				&reservas.Pago{},
				&reservas.Reserva{},
				&reservas.Vehiculo{},
				&reservas.Cliente{},
				&reservas.Usuario{},
				&reservas.Sucursal{},
			)
		}

		if err := dbMig.AutoMigrate(
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

		almacen.Sembrar(dbMig)

		if err := sesion.Sembrar(dbMig); err != nil {
			log.Fatal("no se pudieron sembrar los usuarios: ", err)
		}

		apiUser := os.Getenv("DB_API_USER")
		if apiUser == "" {
			apiUser = "rentcar_api"
		}
		apiPass := os.Getenv("DB_API_PASSWORD")
		if apiPass != "" {
			if err := aplicarPermisosMinimos(dbMig, apiUser, apiPass); err != nil {
				log.Fatal("no se pudieron aplicar permisos mínimos: ", err)
			}
		}

		log.Println("migración, siembra y permisos mínimos completados")
		if *migrar {
			return
		}
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar: ", err)
	}

	r := chi.NewRouter()
	// CORS arriba de todo: responde las peticiones previas (OPTIONS), que
	// no tienen manejador. Mínimo privilegio: se nombra el front, no "*".
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         300, // segundos que el navegador recuerda la respuesta
	}))
	r.Use(chimiddleware.ClientIPFromRemoteAddr) // anota la IP real de la conexión
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
	// Pública, pero con límite: 5 intentos por minuto por IP. El sexto responde 429.
	r.With(httprate.LimitBy(5, time.Minute, porIP)).Post("/login", login.Login)

	conexion, err := db.DB()
	if err != nil {
		log.Fatal("no se pudo obtener la conexión: ", err)
	}
	repo := almacen.NuevoReservasSqlc(almacen.NuevoReservas(db), conexion)
	(&reservas.Manejador{Repo: repo}).Rutas(r, tokens)

	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}
	log.Println("escuchando en el puerto", cfg.Puerto)
	log.Fatal(servidor.ListenAndServe())
}

// porIP identifica al cliente por la IP de la conexión que anotó
// ClientIPFromRemoteAddr, sin confiar en cabeceras manipulables como X-Forwarded-For.
func porIP(r *http.Request) (string, error) {
	ip := chimiddleware.GetClientIP(r.Context())
	return httprate.CanonicalizeIP(ip), nil
}

// aplicarPermisosMinimos crea (o actualiza) el rol de mínimo privilegio para la API
// y le concede únicamente los permisos estrictamente necesarios tabla por tabla.
func aplicarPermisosMinimos(db *gorm.DB, rol, clave string) error {
	claveEsc := strings.ReplaceAll(clave, "'", "''")
	rolEsc := strings.ReplaceAll(rol, `"`, `""`)

	sentencias := []string{
		fmt.Sprintf(`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '%s') THEN
				CREATE ROLE "%s" LOGIN PASSWORD '%s' NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
			ELSE
				ALTER ROLE "%s" WITH LOGIN PASSWORD '%s' NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
			END IF;
		END $$;`, strings.ReplaceAll(rol, "'", "''"), rolEsc, claveEsc, rolEsc, claveEsc),
		fmt.Sprintf(`REVOKE ALL ON ALL TABLES IN SCHEMA public FROM "%s"`, rolEsc),
		fmt.Sprintf(`REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM "%s"`, rolEsc),
		fmt.Sprintf(`GRANT USAGE ON SCHEMA public TO "%s"`, rolEsc),
		fmt.Sprintf(`GRANT SELECT ON TABLE sucursals, usuarios, clientes, vehiculos, pagos, cuentas TO "%s"`, rolEsc),
		fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE reservas TO "%s"`, rolEsc),
		fmt.Sprintf(`GRANT USAGE, SELECT ON SEQUENCE reservas_id_seq TO "%s"`, rolEsc),
	}

	for _, sql := range sentencias {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
