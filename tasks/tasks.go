package tasks

import (
	"fmt"
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
	output, err = ps.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(output)
}

//RunTasks : Corre todas las tareas del slice que recibe
func RunTasks(AllTasks []config.Task, maxParallel int) {
	sem := make(chan bool, maxParallel)
	for _, task := range AllTasks {
		sem <- true
		go runTask(task, sem)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}
