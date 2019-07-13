package main

import (
	"fmt"
	"siper/config"
	"time"
)

func main() {
	for true {
		Conf := config.ReadConfig()
		time.Sleep(time.Second * time.Duration(Conf.Global.CheckNewConfigInterval))
		for _, value := range Conf.Tasks {
			println(value.ID)
			println(value.Command)
		}
		fmt.Println("Concurrency set to: ", Conf.Global.MaxProcessConcurrency)
	}
}
