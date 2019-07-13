package main

import (
	"fmt"
	"os"
	"os/signal"
	"siper/config"
	"syscall"
	"time"
)

//Config : Global donde se guarda la configuracion
var Config config.Conf

const pidFilename string = "siper.pid"

func init() {
	var err error
	Config, err = config.ReadConfig()
	if err != nil {
		fmt.Println("Error: Revise la configuracion")
		os.Exit(1) //Salgo con error porque no pude ni leer la config
	}
}

func salir() {
	os.Exit(0)
}

func signalCatcher() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
	salir()
}

func main() {
	go signalCatcher()
	for {
		Conf, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error: Revise la configuracion")
			os.Exit(1)
		}
		time.Sleep(time.Second * time.Duration(Conf.Global.CheckNewConfigInterval))
		for _, value := range Conf.Tasks {
			fmt.Println(value.ID)
			fmt.Println(value.Command)
		}
	}
}
