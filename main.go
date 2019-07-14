package main

import (
	"fmt"
	"os"
	"os/signal"
	"siper/config"
	"siper/utils"
	"syscall"
	"time"
)

func init() {
	var err error
	err = config.ReadConfig()
	if err != nil {
		fmt.Println("Error: Revise la configuracion")
		os.Exit(1) //Salgo con error porque no pude ni leer la config
	}
}

func signalCatcher() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
	salir()
}

//Salir : Se ejecuta cuando se sale del programa
func salir() {
	os.Exit(0)
}

func main() {
	go signalCatcher()
	var err error
	for {
		err = config.ReadConfig()
		if err != nil {
			fmt.Println("Error: Revise la configuracion")
			os.Exit(1)
		}
		TasksOfTheDay := utils.GetTasksOfTheDay(config.Config.Tasks)
		utils.RunTasks(TasksOfTheDay, config.Config.Global.MaxProcessConcurrency)
		time.Sleep(time.Second * time.Duration(config.Config.Global.CheckNewConfigInterval))
	}
}
