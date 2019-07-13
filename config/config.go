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
	Global struct {
		MaxProcessConcurrency  uint16 `yaml:"max_process_concurrency"`
		CheckNewConfigInterval uint   `yaml:"check_new_config_interval"`
		UnknownAction          string `yaml:"unknown_action"`
		LogDirectory           string `yaml:"log_directory"`
	} `yaml:"global"`
	Oracle struct {
		EtlUser               string `yaml:"etl_user"`
		EtlPassword           string `yaml:"etl_password"`
		FiscoConnectionString string `yaml:"fisco_connection_string"`
		MisgConnectionString  string `yaml:"misg_connection_string"`
	} `yaml:"oracle"`
	Tasks []struct {
		ID      uint16 `yaml:"id"`
		Type    string `yaml:"type"`
		Command string `yaml:"command"`
		Day     uint8  `yaml:"day"`
	}
}

// ReadConfig : Lee el archivo de configuracion.
func ReadConfig() (Conf, error) {
	var c Conf
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //obtengo el path del ejecutable
	yamlFile, err := ioutil.ReadFile(dir + "/config.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
		return c, errors.New(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
		return c, errors.New(err.Error())
	}
	return c, nil
}
