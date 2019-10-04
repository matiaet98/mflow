package main

import (
	"log"
	"os"
	"os/signal"
	"siper/config"
	"siper/tasks"
	"syscall"
	"time"
)

func init() {
	var err error
	err = config.ReadConfig()
	if err != nil {
		log.Fatalln("Error Fatal: Revise la configuracion")
	}
}

func signalCatcher() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
	log.Fatal("Signal de terminacion capturada, se sale")
}

//Salir : Se ejecuta cuando se sale del programa

func main() {
	go signalCatcher()
	var err error
	pendingTasks := tasks.GetPendingTasks(config.Config.Tasks)
	for len(pendingTasks) > 0 {
		err = config.ReadConfig()
		if err != nil {
			log.Fatalln("Error: Revise la configuracion")
		}
		tasks.RunTasks(pendingTasks, config.Config.Global.MaxProcessConcurrency)
		time.Sleep(time.Second * time.Duration(config.Config.Global.CheckNewConfigInterval))
	}
	log.Println("All tasks finished")
}
