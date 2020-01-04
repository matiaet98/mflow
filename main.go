package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/matiaet98/mflow/config"
	"github.com/matiaet98/mflow/global"
	"github.com/matiaet98/mflow/tasks"
	log "github.com/sirupsen/logrus"
)

func init() {
	var err error
	flag.StringVar(&global.TaskFile, "taskfile", "./tasks.json", "Archivo json con el DAG de tareas a correr")
	flag.StringVar(&global.DatasourcesFile, "datasources", "./oracle.json", "Archivo json con los strings de conexion a bases de datos")
	flag.Parse()
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
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
	log.Infoln("Initializing log")
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
	err = tasks.ValidateTaskIds(config.Config.Tasks.Tasks)
	if err != nil {
		log.Fatalf(err.Error())
	}
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
	log.Infoln("All tasks finished")
}
