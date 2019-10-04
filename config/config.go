package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//Config : Global donde se guarda la configuracion
var Config Conf

// Conf : struct de configuracion.
type Conf struct {
	Global `yaml:"global"`
	Oracle `yaml:"oracle"`
	Tasks  []Task `yaml:"tasks"`
}

// Global : Estructura de bloque global
type Global struct {
	MaxProcessConcurrency  int    `yaml:"max_process_concurrency"`
	CheckNewConfigInterval int    `yaml:"check_new_config_interval"`
	LogDirectory           string `yaml:"log_directory"`
}

// Oracle : Estructura de bloque oracle
type Oracle struct {
	EtlUser               string `yaml:"etl_user"`
	EtlPassword           string `yaml:"etl_password"`
	FiscoConnectionString string `yaml:"fisco_connection_string"`
	MisgConnectionString  string `yaml:"misg_connection_string"`
}

//Task : Estructura de una tarea
type Task struct {
	ID      int    `yaml:"id"`
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
	Depends []int  `yaml:"depends"`
}

func getConfigs(path string) error {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //obtengo el path del ejecutable
	yamlFile, err := ioutil.ReadFile(dir + "/" + path)
	if err != nil {
		yamlFile, err = ioutil.ReadFile("./" + path) //antes de salir con error pruebo en el directorio actual
		if err != nil {
			return err
		}
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return err
	}
	return err
}

// ReadConfig : Lee el archivo de configuracion.
func ReadConfig() (err error) {
	err = getConfigs("config.yaml")
	if err != nil {
		log.Panicln(err)
	}
	err = getConfigs("oracle.yaml")
	if err != nil {
		log.Panicln(err)
	}
	return
}
