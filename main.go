package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/middleware"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/respuesta"
	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/tickets"
)

func main() {
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

	almacen := tickets.NuevoAlmacen()
	r.Get("/tickets", almacen.Listar)
	r.Post("/tickets", almacen.Crear)
	r.Get("/tickets/{id}", almacen.Obtener)
	r.Get("/explotar", tickets.Explotar)

	log.Println("mesa de ayuda en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
