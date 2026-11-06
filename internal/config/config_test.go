package config

import (
	"strings"
	"testing"
)

func TestCargarSinDatabaseURLFalla(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("CLAVE_FIRMA", "clave_de_prueba_de_mas_de_32_caracteres_12345")
	t.Setenv("CLAVE_ANA", "ana-test")
	t.Setenv("CLAVE_LUIS", "luis-test")
	_, err := Cargar()
	if err == nil {
		t.Fatal("se esperaba un error por DATABASE_URL vacía y no hubo")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("el error no dice qué falta: %v", err)
	}
}

func TestCargarSinClaveFirmaFalla(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost user=postgres password=test dbname=test")
	t.Setenv("CLAVE_FIRMA", "")
	t.Setenv("CLAVE_ANA", "ana-test")
	t.Setenv("CLAVE_LUIS", "luis-test")
	_, err := Cargar()
	if err == nil {
		t.Fatal("se esperaba un error por CLAVE_FIRMA vacía y no hubo")
	}
	if !strings.Contains(err.Error(), "CLAVE_FIRMA") {
		t.Fatalf("el error no dice qué falta: %v", err)
	}
}

func TestCargarSinClaveAnaFalla(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost user=postgres password=test dbname=test")
	t.Setenv("CLAVE_FIRMA", "clave_de_prueba_de_mas_de_32_caracteres_12345")
	t.Setenv("CLAVE_ANA", "")
	t.Setenv("CLAVE_LUIS", "luis-test")

	_, err := Cargar()
	if err == nil {
		t.Fatal("se esperaba un error por CLAVE_ANA vacía y no hubo")
	}
	if !strings.Contains(err.Error(), "CLAVE_ANA") {
		t.Fatalf("el error no dice qué falta: %v", err)
	}
}

func TestCargarLeeConfiguracionCompleta(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost user=postgres password=test dbname=test")
	t.Setenv("CLAVE_FIRMA", "clave_de_prueba_de_mas_de_32_caracteres_12345")
	t.Setenv("CLAVE_ANA", "ana-test")
	t.Setenv("CLAVE_LUIS", "luis-test")
	t.Setenv("PUERTO", "9090")
	t.Setenv("TIEMPO_ESPERA_SEGUNDOS", "8")

	configuracion, err := Cargar()
	if err != nil {
		t.Fatalf("Cargar devolvió un error: %v", err)
	}
	if configuracion.Puerto != "9090" || configuracion.TiempoEspera.Seconds() != 8 {
		t.Fatalf("configuración de servidor inesperada: %+v", configuracion)
	}
	if configuracion.ClaveAna != "ana-test" || configuracion.ClaveLuis != "luis-test" {
		t.Fatal("no se cargaron las credenciales de prueba")
	}
}
