package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"siper/config"
	"siper/global"
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
	f, err := os.OpenFile(config.Config.LogDirectory+"siper-service.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Println("Initializing log")
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
	global.IDMaster = tasks.CreateMaster()
	pendingTasks := tasks.GetPendingTasks(config.Config.Tasks)
	for len(pendingTasks) > 0 {
		tasks.RunTasks(pendingTasks, config.Config.Global.MaxProcessConcurrency)
		time.Sleep(time.Second * time.Duration(config.Config.Global.CheckNewConfigInterval))
		pendingTasks = tasks.GetPendingTasks(config.Config.Tasks)
	}
	log.Println("All tasks finished")
}
