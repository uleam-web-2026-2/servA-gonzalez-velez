package reservas

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/uleam-web-2026-2/servA-gonzalez-reina/internal/sesion"
)

// rolDePrueba es el rol con el que entran las pruebas de la semana 4:
// el que tiene permiso en MÁS rutas de su matriz.
const rolDePrueba = sesion.RolAgente

// tokensDePrueba firma con un secreto que solo existe en las pruebas.
var tokensDePrueba = sesion.NuevosTokens("secreto-solo-de-las-pruebas")

// enrutadorCon arma el servidor de prueba, SIN repositorio,
// con las mismas rutas y los mismos niveles que usa main.go.
// Con un rol, todas las peticiones entran con el token de un usuario de
// ese rol. Con "" entran sin token.
func enrutadorCon(rol string) http.Handler {
	r := chi.NewRouter()
	(&Manejador{}).Rutas(r, tokensDePrueba) // el mismo Rutas de main.go
	if rol == "" {
		return r
	}
	token, err := tokensDePrueba.Emitir(
		sesion.Usuario{ID: 1, Nombre: "prueba", Rol: rol})
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, p *http.Request) {
		p.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, p)
	})
}

// seCayo dice si el manejador llegó a usar la base, que en las pruebas
// no existe. Sirve para dar un mensaje claro en lugar de un pánico.
func seCayo(t *testing.T, f func()) (cayo bool) {
	t.Helper()
	defer func() {
		if recover() != nil {
			cayo = true
		}
	}()
	f()
	return false
}
