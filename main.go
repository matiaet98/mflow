package main

import (
	"flag"
	"io"
	"mflow/config"
	"mflow/global"
	"mflow/tasks"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	readFlags()
	config.ReadCoreConfig()
	logSetup()
	envLoad()
	config.ReadDatabaseConfig()
	config.ReadTasksConfig()
}

func signalCatcher() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
	log.Fatal("Signal de terminacion capturada, se sale")
}

func logSetup() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	_ = os.Mkdir(config.Config.LogDirectory, os.ModePerm)
	f, err := os.OpenFile(config.Config.LogDirectory+"mflow.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Infoln("Initializing log")
	return
}

func envLoad() {
	for env := range config.Config.EnvVars {
		os.Setenv(config.Config.EnvVars[env].Key, os.ExpandEnv(config.Config.EnvVars[env].Value))
		println(os.Getenv(config.Config.EnvVars[env].Key))
		log.Infof("Setting env var %s value %s", config.Config.EnvVars[env].Key, config.Config.EnvVars[env].Value)
	}
}

func readFlags() {
	flag.StringVar(&global.ConfigFile, "config", "./config.json", "Archivo json con la configuracion central")
	flag.StringVar(&global.TaskFile, "taskfile", "./tasks.json", "Archivo json con el DAG de tareas a correr")
	flag.StringVar(&global.DatasourcesFile, "datasources", "./oracle.json", "Archivo json con los strings de conexion a bases de datos")
	flag.StringVar(&global.ResumeTask, "resume", "", "Master ID de tarea a resumir")
	flag.StringVar(&global.OneShotTask, "oneshot", "", "Para cuando se quiere correr solo una tarea del archivo de tareas. No valida dependencias")
	flag.Parse()
}

func taskValidations() {
	var err error
	err = tasks.ValidateTaskIds(config.Config.Tasks.Tasks)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = tasks.ValidateTaskDependencies(config.Config.Tasks.Tasks)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = tasks.ValidateTaskCiclicDependencies(config.Config.Tasks.Tasks)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func checkConfigChanges() {
	for {
		time.Sleep(time.Second * time.Duration(config.Config.CheckNewConfigInterval))
		config.ReadCoreConfig()
	}
}

func main() {
	go signalCatcher()
	taskValidations()
	var err error
	if global.ResumeTask != "" {
		global.IDMaster, err = strconv.Atoi(global.ResumeTask)
		if err != nil {
			log.Fatal("El numero de tarea a resumir fue mal ingresado")
		}
		tasks.ResetTasks()
	} else {
		err = tasks.CreateMaster()
	}
	if err != nil {
		log.Fatalf("No puedo crear el master: " + err.Error())
	}
	pendingTasks := tasks.GetPendingTasks(config.Config.Tasks.Tasks)
	go checkConfigChanges()
	for len(pendingTasks) > 0 {
		tasks.RunTasks(pendingTasks, config.Config.MaxProcessConcurrency)
		pendingTasks = tasks.GetPendingTasks(config.Config.Tasks.Tasks)
	}
	tasks.EndMaster()
	log.Infoln("All tasks finished")
}
