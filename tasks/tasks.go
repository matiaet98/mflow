package tasks

import (
	"database/sql"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"log"
	"siper/config"
	"siper/global"
	"siper/processes"
	"time"
)

const runningStatus string = "RUNNING"
const noneStatus string = "NONE"
const failedStatus string = "FAILED"
const successStatus string = "SUCCESS"

func runTask(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var output string
	var err error
	var ps processes.Process
	switch task.Type {
	case "bash":
		ps = processes.BashProcess{Command: task.Command}
	case "oracle":
		ps = processes.OracleProcess{
			User:             config.Config.EtlUser,
			Password:         config.Config.EtlPassword,
			ConnectionString: config.Config.FiscoConnectionString,
			Command:          task.Command}
	}
	setTaskStatus(task.ID, runningStatus)
	output, err = ps.Run()
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		log.Panicln(err)
	}
	setTaskStatus(task.ID, successStatus)
	log.Println(output)
}

//GetPendingTasks : Obtiene las tareas pendientes
func GetPendingTasks(AllTasks []config.Task) []config.Task {
	var PendingTasks []config.Task
	for _, task := range AllTasks {
		status, _, err := getTaskStatus(task.ID)
		if err != nil {
			log.Println(err) //no la agrego si hay error
			continue
		}
		if status == noneStatus {
			PendingTasks = append(PendingTasks, task)
		}
	}
	return PendingTasks
}

func dependenciesSucceded(task config.Task) bool {
	for _, dep := range task.Depends {
		status, _, err := getTaskStatus(dep)
		if err != nil {
			log.Println(err) //no la dejo correr si hay error
			return false
		}
		if status != successStatus {
			return false
		}
	}
	return true
}

//RunTasks : Corre todas las tareas del slice que recibe
func RunTasks(Tasks []config.Task, maxParallel int) {
	sem := make(chan bool, maxParallel)
	for _, task := range Tasks {
		if dependenciesSucceded(task) {
			sem <- true
			go runTask(task, sem)
		}
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

//CreateMaster : Crea el master de tareas para esta corrida
func CreateMaster() (ID int, err error) {
	db, err := sql.Open("goracle", config.Config.EtlUser+"/"+config.Config.EtlPassword+"@"+config.Config.FiscoConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Panicln(err)
	}
	const command string = `declare
		l_id number;
		begin
		:l_id := siper.pkg_taskman.create_master;
		end;
	`
	_, err = tx.Exec(command, sql.Named("l_id", sql.Out{Dest: &ID}))
	if err != nil {
		log.Panicln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Panicln(err)
	}
	return
}

func getTaskStatus(IDTask int) (string, time.Time, error) {
	db, err := sql.Open("goracle", config.Config.EtlUser+"/"+config.Config.EtlPassword+"@"+config.Config.FiscoConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Panicln(err)
	}
	var status string
	var fecha time.Time
	const command string = `declare
		status varchar2(20);
		fecha date;
		begin
		siper.pkg_taskman.get_status(:id_master,:id_task,:status,:fecha);
		end;
	`
	_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("id_task", IDTask), sql.Named("status", sql.Out{Dest: &status}), sql.Named("fecha", sql.Out{Dest: &fecha}))
	if err != nil {
		log.Panicln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Panicln(err)
	}
	return status, fecha, nil
}

func setTaskStatus(IDTask int, status string) (string, error) {
	db, err := sql.Open("goracle", config.Config.EtlUser+"/"+config.Config.EtlPassword+"@"+config.Config.FiscoConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Panicln(err)

	}
	var command string
	if status == runningStatus {
		command = `
		begin
		siper.pkg_taskman.start_task(:id_master, :id_task);
		end;
		`
		_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("id_task", IDTask))
	} else {
		command = `
		begin
		siper.pkg_taskman.update_task(:id_master, :id_task,:status);
		end;
		`
		_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("id_task", IDTask), sql.Named("status", status))
	}
	if err != nil {
		log.Panicln(err)

	}
	err = tx.Commit()
	if err != nil {
		log.Panicln(err)
	}
	return status, nil
}
