package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
	UnknownAction          string `yaml:"unknown_action"`
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
	Day     int    `yaml:"day"`
}

// ReadConfig : Lee el archivo de configuracion.
func ReadConfig() (Conf, error) {
	var c Conf
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //obtengo el path del ejecutable
	yamlFile, err := ioutil.ReadFile(dir + "/config.yaml")
	if err != nil {
		yamlFile, err = ioutil.ReadFile("./config.yaml") //antes de salir con error pruebo en el directorio actual
		if err != nil {
			fmt.Printf("%v", err)
			return c, errors.New(err.Error())
		}
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("%v", err)
		return c, errors.New(err.Error())
	}
	return c, nil
}
