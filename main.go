package main

import (
	"flag"
	"io"
	"log"
	"mflow/config"
	"mflow/global"
	"mflow/tasks"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	var err error
	flag.StringVar(&global.TaskFile, "taskfile", "./tasks.json", "Archivo json con el DAG de tareas a correr")
	flag.Parse()
	err = config.ReadConfig()
	if err != nil {
		log.Fatalln("Error Fatal: Revise la configuracion")
	}
	_ = os.Mkdir(config.Config.LogDirectory, os.ModePerm)
	f, err := os.OpenFile(config.Config.LogDirectory+"mflow.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
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
	err = tasks.CreateMaster()
	pendingTasks := tasks.GetPendingTasks(config.Config.Tasks.Tasks)
	if err != nil {
		log.Fatalf("No puedo crear el master: " + err.Error())
	}
	for len(pendingTasks) > 0 {
		tasks.RunTasks(pendingTasks, config.Config.MaxProcessConcurrency)
		time.Sleep(time.Second * time.Duration(config.Config.CheckNewConfigInterval))
		pendingTasks = tasks.GetPendingTasks(config.Config.Tasks.Tasks)
	}
	tasks.EndMaster()
	log.Println("All tasks finished")
}
