# SIPER

## Configuraciones de entorno

### Archivo service.sh

Se deben configurar los valores de las siguientes variables:

```bash
#path absoluto hacia ejecutable de siper, ej: /usr/local/siper/siper
SIPER="/usr/local/siper/siper"
#path absoluto hacia el pid de siper, ej: /usr/local/siper/siper.pid
PIDFILE="/usr/local/siper/siper.pid"
```

### Archivo config.yaml

Este es el fichero maestro de la aplicacion, cada parametro tiene comentario sobre su significado.
Puede ser modificado mientras el proceso corre y el mismo tomara los cambios automaticamente.

_Warning: El unico parametro que **NO** se debe modificar es el ID de cada proceso_
