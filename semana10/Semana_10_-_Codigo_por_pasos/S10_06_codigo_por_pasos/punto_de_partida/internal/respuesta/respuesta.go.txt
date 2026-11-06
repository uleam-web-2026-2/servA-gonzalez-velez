// Package respuesta define la única forma de cuerpo que devuelve el servidor
// (contrato del curso, semana 2): los datos en éxito y {"error": "..."} en fallo.
package respuesta

import (
	"encoding/json"
	"net/http"
)

// JSON escribe datos con el código indicado.
func JSON(w http.ResponseWriter, codigo int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(datos)
}

// Error escribe {"error": mensaje} con el código indicado.
func Error(w http.ResponseWriter, codigo int, mensaje string) {
	JSON(w, codigo, map[string]string{"error": mensaje})
}
