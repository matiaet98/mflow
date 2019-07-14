package utils

import (
	"fmt"
	"siper/config"
	"siper/processes"
	"sync"
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

func runTask(task config.Task, wg *sync.WaitGroup) {
	defer wg.Done()
	var output string
	var ps processes.Process
	switch task.Type {
	case "bash":
		ps = processes.BashProcess{Command: task.Command}
	case "oracle":
		fmt.Println("No implementado aun") //TODO: implementar llamados a oracle
	}
	output = ps.Run()
	fmt.Println(output)
}

//RunTasks : Corre todas las tareas del slice que recibe
func RunTasks(AllTasks []config.Task, maxParallel int) {
	var wg sync.WaitGroup
	for index, task := range AllTasks {
		if index%maxParallel == 0 {
			wg.Wait() //espera a que terminen todas las goroutines que se mandaron en paralelo
		}
		wg.Add(1)
		go runTask(task, &wg)
	}

}
