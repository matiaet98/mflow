package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf : struct de configuracion.
type Conf struct {
	A int64
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
