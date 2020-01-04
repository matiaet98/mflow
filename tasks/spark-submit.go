package tasks

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/matiaet98/mflow/config"
	"github.com/matiaet98/mflow/global"

	log "github.com/sirupsen/logrus"
	_ "github.com/godror/godror" //se abstrae su uso con la libreria sql
)

func runSparkSubmit(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var err error
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.ID + ".log")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	mw := io.MultiWriter(os.Stdout, f1)
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = mw
	command, err := buildCommand(task)
	if err != nil {
		log.Panicf("%v", err)
	}
	cmds := strings.Split(command, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput()
	logger.Infoln(string(out))
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		logger.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
}

func buildCommand(task config.Task) (command string, err error) {
	//variables de entorno primero
	if len(task.EnvVars) > 0 {
		for x := range task.EnvVars {
			command += fmt.Sprintf("%s=%s ", task.EnvVars[x].Key, task.EnvVars[x].Value)
		}
	}

	//spark-submit
	command += "spark-submit "

	//master
	if task.Master != "" {
		command += fmt.Sprintf("--master %s ", task.Master)
	}

	//deploy-mode
	if task.DeployMode != "" {
		command += fmt.Sprintf("--deploy-mode %s ", task.DeployMode)
	}

	//supervise y verbose
	command += "--supervise --verbose "

	//class
	if task.Class != "" {
		command += fmt.Sprintf("--class %s ", task.Class)
	}

	//driver memory
	if task.DriverMemory != "" {
		command += fmt.Sprintf("--driver-memory %s ", task.DriverMemory)
	}

	//executor memory
	if task.ExecutorMemory != "" {
		command += fmt.Sprintf("--executor-memory %s ", task.ExecutorMemory)
	}

	//executor cores
	if task.ExecutorCores != "" {
		command += fmt.Sprintf("--executor-cores %s ", task.ExecutorCores)
	}

	//executor cores
	if task.TotalExecutorCores != "" {
		command += fmt.Sprintf("--total-executor-cores %s ", task.TotalExecutorCores)
	}

	//configuraciones de spark
	if len(task.Confs) > 0 {
		for x := range task.Confs {
			command += fmt.Sprintf("--conf %s=%s ", task.Confs[x].Key, task.Confs[x].Value)
		}
	}

	//Script to execute
	command += fmt.Sprintf("%s ", task.IngestorFile)

	//Arguments for the script
	if len(task.Parameters) > 0 {
		for x := range task.Parameters {
			command += fmt.Sprintf("%s ", task.Parameters[x].Parameter)
		}
	}

	log.Infoln("Comando a ejecutar: " + command)
	return
}
