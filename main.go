package main

import (
	"fmt"
	"siper/config"
	"time"
)

const intervalo time.Duration = time.Second * 2

func main() {
	for true {
		time.Sleep(intervalo)
		conf := config.ReadConfig()
		fmt.Println("The value is: ", conf.A)
	}
}
