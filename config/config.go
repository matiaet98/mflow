package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mflow/global"
	"os"
	"path/filepath"
)

//Config : Global donde se guarda la configuracion
var Config Global

//Ora : Global donde se guarda la configuracion de oracle
var Ora Oracle

// Global : struct de configuracion.
type Global struct {
	MaxProcessConcurrency  int    `json:"max_process_concurrency"`
	CheckNewConfigInterval int    `json:"check_new_config_interval"`
	LogDirectory           string `json:"log_directory"`
	Tasks                  Tasks
}

// Oracle : Estructura de bloque oracle
type Oracle struct {
	Connections []OracleConn `json:"connections"`
}

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
	ID                 int          `json:"id"`
	Type               string       `json:"type"`
	Name               string       `json:"name"`
	Command            string       `json:"command,omitempty"`
	Db                 string       `json:"db,omitempty"`
	Depends            []int        `json:"depends,omitempty"`
	DeployMode         string       `json:"deploy-mode,omitempty"`
	Master             string       `json:"master,omitempty"`
	TotalExecutorCores string       `json:"total-executor-cores,omitempty"`
	ExecutorCores      string       `json:"executor-cores,omitempty"`
	ExecutorMemory     string       `json:"executor-memory,omitempty"`
	DriverMemory       string       `json:"driver-memory,omitempty"`
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

// ReadConfig : Lee el archivo de configuracion.
func ReadConfig() (err error) {
	err = getConfigs("config.json", &Config)
	if err != nil {
		log.Panicln(err)
	}
	err = getConfigs(global.DatasourcesFile, &Ora)
	if err != nil {
		log.Panicln(err)
	}
	err = getConfigs(global.TaskFile, &Config.Tasks)
	if err != nil {
		log.Panicln(err)
	}
	return
}
