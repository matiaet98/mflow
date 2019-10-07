# Mflow

## Configuraciones

### Archivo config.json

Esta es la configuracion global de la aplicacion. Se puede customizar con las siguientes variables:
- max_process_concurrency: Cantidad de procesos que pueden estar corriendo en simultaneo
- check_new_config_interval: Cada cuantos segundos se revisa si los procesos ya terminaron para poder continuar con el DAG.
- log_directory: En que directorio crear los logs de los procesos y el log maestro

### Archivo de tareas tasks.json

En este archivo se define el grafo de procesos a ejecutar. Las variables son las siguientes para cada tarea:
id: El identificador unico de la tarea dentro de su grupo. Es importantisimo que no haya duplicados
name: Un nombre simbolico para la tarea. Sirve mas que nada para poder identificar el archivo de log
type: El tipo de tareas que se pueden correr. Hasta ahora solo acepta de tipo bash y oracle, pero se puede extender facilmente.
depends: Es un array con los ID de las tareas de las cuales depende esta tarea. Si las mismas no fueron completadas con SUCCESS, la misma no puede iniciar.
command: El comando que se desea ejecutar. Esta ligado al tipo de proceso.
Db: Identificador de la conexion en el cual se corre el proceso (para procesos oracle)

### Archivo oracle.json

Este archivo es autodescriptivo, contiene la informacion de los distintos datasources para oracle, identificados por un nombre.

### Configuracion especial para Oracle

- Se debe tener instalado OracleClient full, Oracle InstantClient u Oracle Database.
- Una vez realizado, se deben tener sus librerias en la variable de entorno LD_LIBRARY_PATH, por lo cual, habria que agregar en el /etc/bashrc (en el caso de que el cliente este en /opt):

  ```bash
  export LD_LIBRARY_PATH=/opt/instanclient:$LD_LIBRARY_PATH
  ```

- Si reciben el error **ORA-24408: could not generate unique server group name** es porque hay un mismatch en el hostname del equipo. Para arreglarlo, hay que agregar el nombre que nos provee el comando _hostname_ al archivo /etc/hosts

### TODO

- [ ] tomar archivo de tareas por parametro
- [ ] Crear plugin especifico para spark-submit que sepan interpretar la salida de cada comando. Un ejemplo de como hacerlo esta en el codigo de airflow para el hook de spark-submit.