# MFlow

Matiflow: Herramienta para generar DAGs de ejecucion de procesos. Inicialmente acepta procesos bash y oracle.

- [MFlow](#mflow)
  - [Configuraciones](#configuraciones)
    - [Archivo config.json](#archivo-configjson)
    - [Archivo de tareas tasks.json](#archivo-de-tareas-tasksjson)
    - [Archivo oracle.json](#archivo-oraclejson)
    - [Configuracion especial para Oracle](#configuracion-especial-para-oracle)
    - [Instalacion](#instalacion)
    - [Log](#log)
    - [Ejemplos de corrida](#ejemplos-de-corrida)
    - [TODO](#todo)

## Configuraciones

### Archivo config.json

Esta es la configuracion global de la aplicacion. Se puede customizar con las siguientes variables:
- max_process_concurrency: Cantidad de procesos que pueden estar corriendo en simultaneo
- check_new_config_interval: Cada cuantos segundos se revisa si los procesos ya terminaron para poder continuar con el DAG.
- log_directory: En que directorio crear los logs de los procesos y el log maestro

Ejemplo: **config.json**
~~~json
{
    "max_process_concurrency": 3,
    "check_new_config_interval": 1,
    "log_directory": "./logs/"
}
~~~

### Archivo de tareas tasks.json

En este archivo se define el grafo de procesos a ejecutar. Las variables son las siguientes para cada tarea:
id: El identificador unico de la tarea dentro de su grupo. Es importantisimo que no haya duplicados
name: Un nombre simbolico para la tarea. Sirve mas que nada para poder identificar el archivo de log
type: El tipo de tareas que se pueden correr. Hasta ahora solo acepta de tipo bash y oracle, pero se puede extender facilmente.
depends: Es un array con los ID de las tareas de las cuales depende esta tarea. Si las mismas no fueron completadas con SUCCESS, la misma no puede iniciar.
command: El comando que se desea ejecutar. Esta ligado al tipo de proceso.
Db: Identificador de la conexion en el cual se corre el proceso (para procesos oracle)
Opcionalmente este archivo se puede pasar como parametro usando la opcion **--taskfile**. Esto es muy util para disparar varias instancias de mflow cada uno con su grafo de procesos. De hecho un patron recomendado podria ser disparar un mflow padre donde sus tareas sean otras corridas de mflow, cada una con un archivo distinto de tareas).

Ejemplo: **tasks.json**
~~~json
{
    "tasks": [
        {
            "id": 1,
            "name": "tarea1",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea1.sh"
        },
        {
            "id": 2,
            "name": "tarea2",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea2.sh",
            "depends": [1]
        },
        {
            "id": 3,
            "name": "tarea3",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea3.sh",
            "depends": [1]
        },
        {
            "id": 4,
            "name": "tarea4",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea4.sh",
            "depends": [3,2]
        },
        {
            "id": 5,
            "name": "tarea5",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea5.sh",
            "depends": [2]
        },
        {
            "id": 6,
            "name": "tarea6",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea6.sh",
            "depends": [1,3]
        },
        {
            "id": 7,
            "name": "tarea7",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea7.sh"
        }
    ]
}
~~~
### Archivo oracle.json

Este archivo es autodescriptivo, contiene la informacion de los distintos datasources para oracle, identificados por un nombre.

### Configuracion especial para Oracle

- Se debe tener instalado OracleClient full, Oracle InstantClient u Oracle Database.
- Una vez realizado, se deben tener sus librerias en la variable de entorno LD_LIBRARY_PATH, por lo cual, habria que agregar en el /etc/bashrc (en el caso de que el cliente este en /opt):

  ```bash
  export LD_LIBRARY_PATH=/opt/instanclient:$LD_LIBRARY_PATH
  ```

- Si reciben el error **ORA-24408: could not generate unique server group name** es porque hay un mismatch en el hostname del equipo. Para arreglarlo, hay que agregar el nombre que nos provee el comando _hostname_ al archivo /etc/hosts

### Instalacion
- En una instancia Oracle, crear el usuario MFLOW (con roles resource y connect) y correr el script que se encuentra en db/db.sql
- Correr make build
- Si termino bien, copiar el ejecutable **mflow** y los archivos oracle.json y config.json al directorio deseado
- Agregar directoio al PATH, o agregar un link a mflow en la carpeta /usr/local/bin.

### Log
- MFlow cuenta con un log maestro donde se registran todas sus operaciones, y ademas genera un log para cada tarea que ejecuta.
- Ademas, tambien se puede dar un seguimiento a las tareas con los registros que van quedando en las tablas tasks en la base de datos.

### Ejemplos de corrida

Corrida normal
~~~bash
mflow
~~~

Corrida en background
~~~bash
mflow &> /dev/null &
~~~

Corrida con distintos archivos de tareas
~~~bash
mflow --taskfile /home/mestevez/fiscar/dag1.json &> /dev/null &
mflow --taskfile /home/mestevez/fiscar/dag2.json &> /dev/null &
time mflow --taskfile /home/mestevez/siper/c1.json &> /dev/null &
~~~

### TODO
- [ ] Scheduling: Hoy la herramienta de por si no maneja un scheduler interno ya que para eso existe CERO (La idea es que la corrida de mflow se dispare a travez de tivoli). No obstante, si fuese necesario, se pueden usar herramientas nativas de linux como anacron para lograr el scheduling.
- [ ] Retry: Agregar la opcion de reintentar una tarea cuando falla, y que sea configurable la cantidad de reintentos posibles
- [ ] Crear plugin especifico para spark-submit que sepan interpretar la salida de cada comando. Un ejemplo de como hacerlo esta en el codigo de airflow para el hook de spark-submit.
- [ ] Crear Dockerfile para containerizar el proceso
- [ ] Test units
- [ ] Boilerplate con la estructura necesaria para crear mas plugins
- [ ] Exponer un endpoint REST para consultar las tareas
- [ ] Crear uana webapp que consuma dicho endpoint para mostrar en un grafico las tareas. Idealmente podria ser un grafico tipo Sankey usando la libreria d3.js: http://bl.ocks.org/d3noob/5028304