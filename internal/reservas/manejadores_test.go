package reservas // el MISMO paquete que manejadores.go

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// enrutador arma el servidor de prueba SIN base de datos (DB: nil).
// Sirve porque todas las pruebas de este archivo responden ANTES
// de la línea "desde aquí hace falta la base de datos".
func enrutador() http.Handler {
	r := chi.NewRouter()
	(&Manejador{DB: nil}).Rutas(r)
	return r
}

func TestCrear(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		codigo int
	}{
		{
			nombre: "JSON roto responde 400",
			cuerpo: `{"ClienteID": 1`,
			codigo: http.StatusBadRequest,
		},
		{
			nombre: "estado inventado responde 422",
			cuerpo: `{"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-04","Total":10,"Estado":"inventado"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "cliente id cero responde 422",
			cuerpo: `{"ClienteID":0,"VehiculoID":1,"FechaInicio":"2026-11-01","FechaFin":"2026-11-04","Total":10,"Estado":"pendiente"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "fecha fin no posterior responde 422",
			cuerpo: `{"ClienteID":1,"VehiculoID":1,"FechaInicio":"2026-11-04","FechaFin":"2026-11-01","Total":10,"Estado":"pendiente"}`,
			codigo: http.StatusUnprocessableEntity,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(http.MethodPost, "/reservas",
				strings.NewReader(caso.cuerpo))
			grabadora := httptest.NewRecorder()
			enrutador().ServeHTTP(grabadora, peticion)
			if grabadora.Code != caso.codigo {
				t.Fatalf("se esperaba %d y llegó %d con cuerpo %s",
					caso.codigo, grabadora.Code, grabadora.Body.String())
			}
		})
	}
}

func TestVerUnoConIDQueNoEsNumeroResponde400(t *testing.T) {
	peticion := httptest.NewRequest(http.MethodGet, "/reservas/abc", nil)
	grabadora := httptest.NewRecorder()
	enrutador().ServeHTTP(grabadora, peticion)
	if grabadora.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400 y llegó %d", grabadora.Code)
	}
}

func TestListarConLimitNoNumericoResponde400(t *testing.T) {
	peticion := httptest.NewRequest(http.MethodGet, "/reservas?limit=abc", nil)
	grabadora := httptest.NewRecorder()
	enrutador().ServeHTTP(grabadora, peticion)
	if grabadora.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400 y llegó %d", grabadora.Code)
	}
}

func TestCambiarEstadoJSONRotoResponde400(t *testing.T) {
	peticion := httptest.NewRequest(http.MethodPatch, "/reservas/1",
		strings.NewReader(`{roto`))
	grabadora := httptest.NewRecorder()
	enrutador().ServeHTTP(grabadora, peticion)
	if grabadora.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400 y llegó %d", grabadora.Code)
	}
}

func TestCambiarEstadoInventadoResponde422(t *testing.T) {
	peticion := httptest.NewRequest(http.MethodPatch, "/reservas/1",
		strings.NewReader(`{"Estado":"inventado"}`))
	grabadora := httptest.NewRecorder()
	enrutador().ServeHTTP(grabadora, peticion)
	if grabadora.Code != http.StatusUnprocessableEntity {
		t.Fatalf("se esperaba 422 y llegó %d", grabadora.Code)
	}
}

func TestTransicionPermitida(t *testing.T) {
	casos := []struct {
		de, a string
		ok    bool
	}{
		{"pendiente", "confirmada", true},
		{"pendiente", "cancelada", true},
		{"pendiente", "en_curso", false},
		{"confirmada", "en_curso", true},
		{"finalizada", "pendiente", false},
		{"cancelada", "confirmada", false},
		{"en_curso", "finalizada", true},
		{"pendiente", "pendiente", true},
	}
	for _, c := range casos {
		if got := transicionPermitida(c.de, c.a); got != c.ok {
			t.Fatalf("%s -> %s: se esperaba %v y llegó %v", c.de, c.a, c.ok, got)
		}
	}
}
