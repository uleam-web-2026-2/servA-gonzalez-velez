# Semana 11 · nota de Windows

Al redirigir la salida de `docker run` con la redirección de todos los flujos de PowerShell, las evidencias quedaron en UTF-16 e incluyeron el diagnóstico de error de PowerShell junto con el error esperado del contenedor. Para guardar texto limpio y UTF-8, capturamos `stdout` y `stderr` de Docker con `subprocess` y los escribimos explícitamente en UTF-8. Docker Desktop y el contenedor funcionaron correctamente.
