package reservas

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Prueba de integración: arma el servidor entero (rutas, guardia y
// manejadores) y le envía cada petición SIN token.
//
// Es la prueba que cuida la puerta: ninguna ruta protegida responde sin
// token. Si mañana alguien saca una ruta del grupo protegido, esta
// prueba lo dice antes de que el cambio llegue a los usuarios.
func TestSinTokenLasRutasProtegidasResponden401(t *testing.T) {
	// Una fila por cada ruta de la matriz que NO sea pública.
	protegidas := []struct{ metodo, ruta, cuerpo string }{
		{"GET", "/reservas", ""},
		{"POST", "/reservas", `{"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-04","Total":105,"Estado":"pendiente"}`},
		{"GET", "/reservas/1", ""},
		{"PATCH", "/reservas/1", `{"estado":"confirmada"}`},
		{"PUT", "/reservas/1", `{"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-04","Total":105,"Estado":"pendiente"}`},
		{"DELETE", "/reservas/1", ""},
		{"GET", "/sucursales", ""},
		{"GET", "/clientes", ""},
		{"GET", "/usuarios", ""},
		{"GET", "/pagos", ""},
	}

	servidor := enrutadorCon("") // "" = nadie inició sesión

	for _, p := range protegidas {
		peticion := httptest.NewRequest(p.metodo, p.ruta,
			strings.NewReader(p.cuerpo))
		grabadora := httptest.NewRecorder()

		cayo := seCayo(t, func() { servidor.ServeHTTP(grabadora, peticion) })

		if cayo {
			t.Errorf("%s %s sin token: el guardia no la detuvo y la "+
				"petición llegó hasta la base de datos", p.metodo, p.ruta)
			continue
		}
		if grabadora.Code != http.StatusUnauthorized {
			t.Errorf("%s %s sin token: se esperaba 401 y llegó %d",
				p.metodo, p.ruta, grabadora.Code)
		}
	}
}
