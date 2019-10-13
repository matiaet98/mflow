package tasks

import (
	log "github.com/sirupsen/logrus"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"io"
	"mflow/config"
	"mflow/global"
	"os"
	"os/exec"
	"strconv"
)

func runBash(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var err error
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.Name + ".log")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	mw := io.MultiWriter(os.Stdout, f1)
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableColors:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = mw
	cmd := exec.Command(task.Command)
	out, err := cmd.CombinedOutput() //este chabon aparte de combinar stderr y stdout tambien hace el Run... poco intuitivo
	logger.Println(out)
	logger.Infoln(out)
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		logger.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
}
