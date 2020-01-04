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
  - [Desarrollo](#desarrollo)
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
id: Un nombre identificador unico para la tarea.
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
            "id": "tarea1",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea1.sh"
        },
        {
            "id": "tarea2",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea2.sh",
            "depends": ["tarea1"]
        },
        {
            "id": "tarea3",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea3.sh",
            "depends": ["tarea2"]
        },
        {
            "id": "tarea4",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea4.sh",
            "depends": ["tarea1","tarea3"]
        },
        {
            "id": "sp1_tickets",
            "type": "spark",
            "master": "spark://hdp:7077",
            "deploy-mode": "client",
            "driver-memory": "1g",
            "executor-memory": "4g",
            "executor-cores": "5",
            "total-executor-cores": "10",
            "ingestor-file": "/opt/ingestion/spark/tickets.py",
            "confs" : [
                {
                    "key":"spark.driver.maxResultSize",
                    "value":"4g"
                }
            ]
        },
        {
            "id": "tarea6",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea6.sh",
            "depends": ["tarea5","tarea2"]
        },
        {
            "id": "tarea7",
            "type": "bash",
            "command": "/home/mestevez/tmp/tarea7.sh"
        },
        {
            "id": "sp1",
            "type": "spark",
            "master": "spark://hdp:7077",
            "deploy-mode": "client",
            "driver-memory": "1g",
            "executor-memory": "4g",
            "executor-cores": "5",
            "total-executor-cores": "10",
            "ingestor-file": "/opt/ingestion/spark/scala/fisca-ingestion.jar",
            "class": "ar.gob.afip.fiscalizacion.ingestion.batch.raw.ftContribuyentes",
            "confs" : [
                {
                    "key":"spark.driver.maxResultSize",
                    "value":"4g"
                }
            ]
        },
        {
            "id": "siper_desvio141",
            "type": "oracle",
            "command": "begin \n siper.test_data_gen.DESVIO_14_1; \n end;",
            "db": "siper_fisco"
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
- El archivo oracle.json puede pasarse tambien como parametro con el flag **--datasources**, pero hay que tener en cuenta que siempre tiene que estar el conector con el nombre "mflow", que es el que usa la aplicacion como backend para sincronizar las tareas. En caso de no pasarse por parametro se usa el que esta en el directorio actual.

Ejemplo: **oracle.json**
~~~json
{
    "connections":[
        {
            "name": "mflow",
            "connection_string": "(DESCRIPTION=(ADDRESS=(PROTOCOL=tcp)(HOST=10.30.205.127)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=fisco)))",
            "user": "MFLOW",
            "password": "MFLOW"
        },
	{
            "name": "siper_fisco",
            "connection_string": "10.30.205.127:1521/fisco",
            "user": "SIPER",
            "password": "SIPER"
        }
    ]
}
~~~


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

Corrida con distintos archivos de tareas y datasources
~~~bash
mflow --taskfile /home/mestevez/fiscar/dag1.json --datasources desa.json &> /dev/null &
mflow --taskfile /home/mestevez/fiscar/dag2.json &> /dev/null &
mflow --taskfile /home/mestevez/fiscar/dag2.json --datasources prod.json &> /dev/null &
time mflow --taskfile /home/mestevez/siper/c1.json &> /dev/null &
~~~

## Desarrollo

- Esta aplicacion fue compilada con go 1.13
- Es necesario tener make para poder buildear el proyecto
- Se pueden bajar las dependencias del proyecto corriendo make deps

### TODO
- [ ] Scheduling: Hoy la herramienta de por si no maneja un scheduler interno ya que para eso existe CERO (La idea es que la corrida de mflow se dispare a travez de tivoli). No obstante, si fuese necesario, se pueden usar herramientas nativas de linux como anacron para lograr el scheduling.
- [ ] Retry: Agregar la opcion de reintentar una tarea cuando falla, y que sea configurable la cantidad de reintentos posibles
- [ ] Crear Dockerfile para containerizar el proceso
- [ ] Test units
- [ ] Boilerplate con la estructura necesaria para crear mas plugins
- [ ] Exponer un endpoint REST para consultar las tareas