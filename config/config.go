package config

import (
	"encoding/json"
	"io/ioutil"
	"mflow/global"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

//Config : Global donde se guarda la configuracion
var Config Global

//Ora : Global donde se guarda la configuracion de oracle
var Ora Oracle

// Global : struct de configuracion.
type Global struct {
	MaxProcessConcurrency  int       `json:"max_process_concurrency"`
	CheckNewConfigInterval int       `json:"check_new_config_interval"`
	LogDirectory           string    `json:"log_directory"`
	EnvVars                []EnvVars `json:"env,omitempty"`
	Tasks                  Tasks
}

// Oracle : Estructura de bloque oracle
type Oracle struct {
	Connections []OracleConn `json:"connections"`
}

//Tasks : Grupo de tareas
type Tasks struct {
	Tasks []Task `json:"tasks"`
}

// OracleConn : Estructura de bloque oracle
type OracleConn struct {
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	User             string `json:"user"`
	Password         string `json:"password"`
}

//Task : Estructura de una tarea
type Task struct {
	ID                 string       `json:"id"`
	Type               string       `json:"type"`
	Command            string       `json:"command,omitempty"`
	Db                 string       `json:"db,omitempty"`
	Depends            []string     `json:"depends,omitempty"`
	Master             string       `json:"master,omitempty"`
	DeployMode         string       `json:"deploy-mode,omitempty"`
	Name               string       `json:"name,omitempty"`
	TotalExecutorCores string       `json:"total-executor-cores,omitempty"`
	ExecutorCores      string       `json:"executor-cores,omitempty"`
	ExecutorMemory     string       `json:"executor-memory,omitempty"`
	NumExecutors       string       `json:"num-executors,omitempty"`
	DriverMemory       string       `json:"driver-memory,omitempty"`
	DriverCores        string       `json:"driver-cores,omitempty"`
	Verbose            string       `json:"verbose,omitempty"`
	Supervise          string       `json:"supervise,omitempty"`
	IngestorFile       string       `json:"ingestor-file,omitempty"`
	Class              string       `json:"class,omitempty"`
	Parameters         []Parameters `json:"parameters,omitempty"`
	Confs              []SparkConf  `json:"confs,omitempty"`
	EnvVars            []EnvVars    `json:"env,omitempty"`
}

//SparkConf configuraciones de spark que se desean pasar a un proceso
type SparkConf struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//EnvVars Variables de entorno que se desean setear antes de ejecutar un comando
type EnvVars struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Parameters : Parametros para las tareas que lanzamos
type Parameters struct {
	Parameter string `json:"parameter"`
}

func getConfigs(path string, conf interface{}) error {
	jfile, err := ioutil.ReadFile(path) //por si viene como parametro
	if err != nil {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //obtengo el path del ejecutable
		jfile, err = ioutil.ReadFile(dir + "/" + path)
		if err != nil {
			jfile, err = ioutil.ReadFile("./" + path) //antes de salir con error pruebo en el directorio actual
			if err != nil {
				return err
			}
		}
	}
	err = json.Unmarshal(jfile, &conf)
	if err != nil {
		return err
	}
	return err
}

// ReadCoreConfig : Lee el archivo de configuracion core.
func ReadCoreConfig() {
	err := getConfigs(global.ConfigFile, &Config)
	if err != nil {
		log.Error("No se pudo leer la configuracion core")
		log.Fatalln(err)
	}
	return
}

// ReadDatabaseConfig : Lee el archivo de configuracion de datasources.
func ReadDatabaseConfig() {
	err := getConfigs(global.DatasourcesFile, &Ora)
	if err != nil {
		log.Error("No se pudo leer la configuracion de la base")
		log.Fatalln(err)
	}
	return
}

// ReadTasksConfig : Lee el archivo de configuracion de tareas.
func ReadTasksConfig() {
	err := getConfigs(global.TaskFile, &Config.Tasks)
	if err != nil {
		log.Error("No se pudo leer la configuracion de tareas")
		log.Fatalln(err)
	}
	if global.OneShotTask != "" {
		var OSTask []Task
		for _, task := range Config.Tasks.Tasks {
			if global.OneShotTask == task.ID {
				task.Depends = nil
				OSTask = append(OSTask, task)
				break
			}
		}
		Config.Tasks.Tasks = OSTask
	}
	return
}
