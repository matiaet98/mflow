package tasks

import (
	"context"
	"database/sql"
	"io"
	"os"
	"strconv"

	"mflow/config"
	"mflow/global"

	_ "github.com/godror/godror" //se abstrae su uso con la libreria sql
	log "github.com/sirupsen/logrus"
)

func runOracle(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var err error
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.ID + ".log")
	if err != nil {
		log.Panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f1)
	defer f1.Close()
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = mw
	conn := getConnection(task.Db)
	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		log.Panicln(err)
	}
	_, err = tx.Exec(task.Command)
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		logger.Warnln(err)
		tx.Rollback()
		return
	}
	setTaskStatus(task.ID, successStatus)
	tx.Commit()
	return
}
