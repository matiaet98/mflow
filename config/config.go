package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Conf : struct de configuracion.
type Conf struct {
	Global struct {
		MaxProcessConcurrency  int16 `yaml:"max_process_concurrency"`
		CheckNewConfigInterval int8  `yaml:"check_new_config_interval"`
	} `yaml:"global"`
	Oracle struct {
		EtlUser     string `yaml:"etl_user"`
		EtlPassword string `yaml:"etl_password"`
	} `yaml:"oracle"`
	Tasks []struct {
		ID      int16  `yaml:"id"`
		Type    string `yaml:"type"`
		Command string `yaml:"command"`
		Day     int8   `yaml:"day"`
	}
}

// ReadConfig : Lee el archivo de configuracion.
func ReadConfig() Conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	var c Conf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	return c
}
