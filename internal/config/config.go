// Package config lee los valores que cambian de máquina en máquina.
// El código no sabe cuánto valen: solo sabe cómo se llaman.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config es todo lo que el servidor necesita saber del entorno donde corre.
type Config struct {
	DatabaseURL  string        // secreto: lleva la contraseña
	Puerto       string        // configuración: cambia por máquina
	TiempoEspera time.Duration // configuración: cuánto se espera antes de cortar
	ClaveFirma   string        // secreto: firma los tokens JWT (mínimo 32 letras/números)
	ClaveAna     string        // contraseña de ejemplo del usuario de pruebas ana
	ClaveLuis    string        // contraseña de ejemplo del usuario de pruebas luis
}

// Cargar lee .env si existe, después las variables del sistema.
// Se detiene con un error claro si falta un valor sin valor por defecto.
func Cargar() (Config, error) {
	// Si no hay .env no es un error: en el servidor real las variables
	// ya vienen puestas por el sistema.
	_ = godotenv.Load()

	var c Config

	c.DatabaseURL = os.Getenv("DATABASE_URL")
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("falta DATABASE_URL " +
			"(copie .env.example a .env y complete la contraseña)")
	}

	c.ClaveFirma = os.Getenv("CLAVE_FIRMA")
	if c.ClaveFirma == "" {
		return c, fmt.Errorf("falta CLAVE_FIRMA " +
			"(copie .env.example a .env y defina una clave de al menos 32 caracteres)")
	}

	c.ClaveAna = os.Getenv("CLAVE_ANA")
	if c.ClaveAna == "" {
		return c, fmt.Errorf("falta CLAVE_ANA " +
			"(defina una contraseña de ejemplo para el usuario ana)")
	}

	c.ClaveLuis = os.Getenv("CLAVE_LUIS")
	if c.ClaveLuis == "" {
		return c, fmt.Errorf("falta CLAVE_LUIS " +
			"(defina una contraseña de ejemplo para el usuario luis)")
	}

	c.Puerto = os.Getenv("PUERTO")
	if c.Puerto == "" {
		c.Puerto = "8080" // valor por defecto: no es un secreto y casi nunca cambia
	}

	segundos := os.Getenv("TIEMPO_ESPERA_SEGUNDOS")
	if segundos == "" {
		segundos = "5"
	}
	n, err := strconv.Atoi(segundos)
	if err != nil || n <= 0 {
		return c, fmt.Errorf("TIEMPO_ESPERA_SEGUNDOS debe ser un entero "+
			"mayor que cero, llegó %q", segundos)
	}
	c.TiempoEspera = time.Duration(n) * time.Second

	return c, nil
}
