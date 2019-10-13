package tasks

import (
	log "github.com/sirupsen/logrus"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"io"
	"mflow/config"
	"mflow/global"
	"mflow/processes"
	"os"
	"strconv"
)

func runOracle(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var output string
	var err error
	var ps processes.Process
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.Name + ".log")
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f1)
	defer f1.Close()
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableColors:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = mw
	conn := getConnection(task.Db)
	ps = processes.OracleProcess{
		User:             conn.User,
		Password:         conn.Password,
		ConnectionString: conn.ConnectionString,
		Command:          task.Command}
	output, err = ps.Run()
	logger.Println(output)
	logger.Infoln(output)
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		logger.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
}
