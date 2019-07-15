package tasks

import (
	"database/sql"
	"fmt"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"siper/config"
	"siper/processes"
	"time"
)

func getDay() int {
	return time.Now().Day()
}

// GetTasksOfTheDay : Devuelve un slice con las tareas del dia
func GetTasksOfTheDay(AllTasks []config.Task) []config.Task {
	today := getDay()
	var tasks []config.Task

	for _, task := range AllTasks {
		if task.Day == today {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

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
	setTaskStatus(task.ID, "RUNNING")
	output, err = ps.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		setTaskStatus(task.ID, "FAILED")
		return
	}
	setTaskStatus(task.ID, "SUCCESS")
	fmt.Println(output)
}

//RunTasks : Corre todas las tareas del slice que recibe
func RunTasks(AllTasks []config.Task, maxParallel int) {
	var PendingTasks []config.Task
	for _, task := range AllTasks {
		status, fecha, err := getTaskStatus(task.ID)
		if err != nil {
			fmt.Println("Error: ", err) //no la agrego si hay error
			continue
		}
		if status == "RUNNING" {
			continue
		}
		if status == "NONE" {
			PendingTasks = append(PendingTasks, task)
			continue
		}
		year, month, day := time.Now().Date()
		if int(fecha.Sub(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))) < 0 { //si es anterior a hoy
			PendingTasks = append(PendingTasks, task)
			continue
		}
	}

	sem := make(chan bool, maxParallel)
	for _, task := range PendingTasks {
		sem <- true
		go runTask(task, sem)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func getTaskStatus(ID int) (string, time.Time, error) {
	db, err := sql.Open("goracle", config.Config.EtlUser+"/"+config.Config.EtlPassword+"@"+config.Config.FiscoConnectionString)
	if err != nil {
		fmt.Println(err)
		return "", time.Now(), err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return "", time.Now(), err
	}
	var status string
	var fecha time.Time
	const command string = `declare
		status varchar2(4000);
		fecha date;
		begin
		siper.pkg_taskman.get_status(:task_id,:status,:fecha);
		end;
	`
	_, err = tx.Exec(command, sql.Named("task_id", ID), sql.Named("status", sql.Out{Dest: &status}), sql.Named("fecha", sql.Out{Dest: &fecha}))
	if err != nil {
		fmt.Println(err)
		return "", time.Now(), err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return "", time.Now(), err
	}
	return status, fecha, nil
}

func setTaskStatus(ID int, status string) (string, error) {
	db, err := sql.Open("goracle", config.Config.EtlUser+"/"+config.Config.EtlPassword+"@"+config.Config.FiscoConnectionString)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var command string
	if status == "RUNNING" {
		command = `
		begin
		siper.pkg_taskman.start_task(:task_id);
		end;
		`
		_, err = tx.Exec(command, sql.Named("task_id", ID))
	} else {
		command = `
		begin
		siper.pkg_taskman.update_task(:task_id,:status);
		end;
		`
		_, err = tx.Exec(command, sql.Named("task_id", ID), sql.Named("status", status))
	}
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return status, nil
}
