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
	"strings"
)

func runBash(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var err error
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.Name + ".log")
	if err != nil {
		log.Panicln(err)
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
	cmds := strings.Split(task.Command, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput() //este chabon aparte de combinar stderr y stdout tambien hace el Run... poco intuitivo
	logger.Infoln(string(out))
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		logger.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
}
