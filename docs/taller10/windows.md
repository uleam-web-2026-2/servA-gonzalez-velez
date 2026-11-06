# Ajustes y resolución de incidencias en Windows

## 1. Caracteres especiales en rutas y compilación de Go
- **Falla:** Al ejecutar `go test ./...`, el compilador de Go arrojaba `malformed import path ".../Semana_10_-_Código_por_pasos/...": invalid char 'ó'`.
- **Causa:** Go no admite caracteres no ASCII (tildes) en los nombres de directorios que contienen código fuente Go bajo un módulo. Además, los archivos de ejemplo en `fase2_dto` y `fase3_sqlc` tenían extensión `.go` e importaban paquetes ajenos (`canchas/...`), lo que intentaba compilarse al correr los tests del proyecto.
- **Solución:** Se normalizó el nombre de la carpeta a caracteres ASCII (`Semana_10_-_Codigo_por_pasos`) y se renombraron los archivos `.go` de ejemplo a `.txt` (siguiendo el mismo criterio del profesor en `punto_de_partida`), permitiendo que `go test ./...` analice y ejecute únicamente el código del proyecto.

## 2. Instalación y disponibilidad de `sqlc` en Windows
- **Falla:** El comando `sqlc` no era reconocido en PowerShell/cmd de Windows por defecto.
- **Causa:** La compilación mediante `go install` en Windows de `sqlc` v1.31.1 descarga múltiples dependencias complejas y puede demorar o carecer de configuración de CGO.
- **Solución:** Se descargó el binario oficial precompilado de Windows (`sqlc_1.31.1_windows_amd64.zip`) directamente desde los releases oficiales de `sqlc-dev/sqlc`, se extrajo `sqlc.exe` en `$env:USERPROFILE\go\bin` y se aseguró su inclusión en la variable de entorno `PATH`.
