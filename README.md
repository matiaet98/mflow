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

### Oracle

- Se debe tener instalado OracleClient full, Oracle InstantClient u Oracle Database.
- Una vez realizado, se deben tener sus librerias en la variable de entorno LD_LIBRARY_PATH, por lo cual, habria que agregar en el /etc/bashrc (en el caso de que el cliente este en /opt):

  ```bash
  export LD_LIBRARY_PATH=/opt/instanclient:$LD_LIBRARY_PATH
  ```

- Si reciben el error **ORA-24408: could not generate unique server group name** es porque hay un mismatch en el hostname del equipo. Para arreglarlo, hay que agregar el nombre que nos provee el comando _hostname_ al archivo /etc/hosts

### TODO

- [ ] Toda la implementacion de Loggging
