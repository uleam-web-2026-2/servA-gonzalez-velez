package respuesta

import (
	"encoding/json"
	"net/http"
)

type detalleError struct {
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

type envoltura struct {
	OK    bool          `json:"ok"`
	Datos any           `json:"datos,omitempty"`
	Error *detalleError `json:"error,omitempty"`
}

func escribir(w http.ResponseWriter, estado int, cuerpo envoltura) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(cuerpo)
}

// Exito responde con el estado indicado y los datos dentro de la envoltura.
func Exito(w http.ResponseWriter, estado int, datos any) {
	escribir(w, estado, envoltura{OK: true, Datos: datos})
}

// JSON escribe directamente cualquier dato como JSON con el código HTTP indicado.
func JSON(w http.ResponseWriter, estado int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(datos)
}

// Error responde con el formato de error adecuado.
// Si recibe un solo texto, responde con {"error": mensaje} según el contrato del curso.
// Si recibe dos textos (codigo, mensaje), preserva la envoltura legacy de semanas anteriores.
func Error(w http.ResponseWriter, estado int, args ...string) {
	if len(args) == 1 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(estado)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": args[0]})
		return
	}
	codigo, mensaje := "", ""
	if len(args) > 0 {
		codigo = args[0]
	}
	if len(args) > 1 {
		mensaje = args[1]
	}
	escribir(w, estado, envoltura{OK: false,
		Error: &detalleError{Codigo: codigo, Mensaje: mensaje}})
}

