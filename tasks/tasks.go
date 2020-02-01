package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"mflow/config"
	"mflow/global"
	"strconv"
	"time"

	_ "github.com/godror/godror" //se abstrae su uso con la libreria sql
	log "github.com/sirupsen/logrus"
)

const (
	runningStatus string = "RUNNING"
	noneStatus    string = "NONE"
	startedStatus string = "STARTED"
	failedStatus  string = "FAILED"
	successStatus string = "SUCCESS"
	endedStatus   string = "ENDED"
)

func runTask(task config.Task, sem chan bool) {
	switch task.Type {
	case "bash":
		runBash(task, sem)
		break
	case "oracle":
		runOracle(task, sem)
		break
	case "spark":
		runSparkSubmit(task, sem)
		break
	}
	log.Infoln("Finalizo la tarea " + task.ID)
}

//GetPendingTasks : Obtiene las tareas pendientes
func GetPendingTasks(AllTasks []config.Task) []config.Task {
	var PendingTasks []config.Task
	for _, task := range AllTasks {
		status, _, err := getTaskStatus(task.ID)
		if err != nil {
			log.Infoln(err) //no la agrego si hay error
			continue
		}
		if status == noneStatus {
			PendingTasks = append(PendingTasks, task)
		}
	}
	return PendingTasks
}

func dependenciesStatus(task config.Task) string {
	for _, dep := range task.Depends {
		status, _, err := getTaskStatus(dep)
		if err != nil {
			log.Warnln(err) //no la dejo correr si hay error
			return failedStatus
		}
		if status == failedStatus {
			return status
		}
		if status == noneStatus || status == runningStatus {
			return status
		}
		if status == successStatus {
			continue
		}
	}
	return successStatus
}

//RunTasks : Corre todas las tareas del slice que recibe
func RunTasks(Tasks []config.Task, maxParallel int) {
	sem := make(chan bool, maxParallel)
	for _, task := range Tasks {
		if dependenciesStatus(task) == successStatus {
			sem <- true
			log.Infoln("Iniciando la tarea " + task.ID)
			setTaskStatus(task.ID, runningStatus)
			go runTask(task, sem)
		} else if dependenciesStatus(task) == failedStatus {
			setTaskStatus(task.ID, failedStatus)
		}
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func getConnection(name string) config.OracleConn {
	var conn config.OracleConn
	for _, x := range config.Ora.Connections {
		if x.Name == name {
			conn = x
			break
		}
	}
	return conn
}

//CreateMaster : Crea el master de tareas para esta corrida
func CreateMaster() (err error) {
	conn := getConnection("mflow")

	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	tx.QueryRow(`select mflow.seq_tasks_master.nextval from dual`).Scan(&global.IDMaster)

	command := fmt.Sprintf(`
		insert into mflow.tasks_master(id,start_date,end_date,status)
		values(%v,sysdate,null,'%v')
	`, global.IDMaster, startedStatus)
	_, err = tx.Exec(command)
	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return
}

//EndMaster : Termina con el master de tareas para esta corrida
func EndMaster() {
	conn := getConnection("mflow")
	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Panicln(err)
	}
	command := fmt.Sprintf(`
		update tasks_master
		set status = '%v', end_date = sysdate
		where id = %v
	`, endedStatus, global.IDMaster)
	_, err = tx.Exec(command)
	if err != nil {
		log.Panicln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Panicln(err)
	}
	return
}

func createTask(taskID string) (err error) {
	conn := getConnection("mflow")

	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	command := fmt.Sprintf(`
		insert into mflow.tasks(id_master,id_task,start_date,end_date,status)
    	values (%v,'%v',sysdate,null,'%v')
	`, global.IDMaster, taskID, noneStatus)
	_, err = tx.Exec(command)
	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return
}

//ResetTasks : reset tasks
func ResetTasks() {
	conn := getConnection("mflow")
	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	const command string = `delete mflow.tasks where id_master = :id_master and status = :status`
	_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("status", runningStatus))
	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func getTaskStatus(IDTask string) (string, time.Time, error) {
	conn := getConnection("mflow")
	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	var status string
	var fecha time.Time
	const command string = `declare
		status varchar2(20);
		fecha date;
		begin
		mflow.pkg_taskman.get_status(:id_master,:id_task,:status,:fecha);
		end;
	`
	_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("id_task", IDTask), sql.Named("status", sql.Out{Dest: &status}), sql.Named("fecha", sql.Out{Dest: &fecha}))
	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return status, fecha, nil
}

func setTaskStatus(IDTask string, status string) (string, error) {
	conn := getConnection("mflow")
	db, err := sql.Open("godror", conn.User+"/"+conn.Password+"@"+conn.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	var endDate string
	if status == runningStatus {
		endDate = "null"
	} else {
		endDate = "sysdate"
	}
	command := fmt.Sprintf(`
	update tasks
    set status = '%s',
    end_date = %s
    where id_master = %s
    and id_task = '%s'
	`, status, endDate, strconv.Itoa(global.IDMaster), IDTask)
	_, err = tx.Exec(command)

	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return status, nil
}

//ValidateTaskIds Valida que no haya IDs repetidos en el archivo de tareas proporcionado
func ValidateTaskIds(AllTasks []config.Task) (err error) {
	dup := make(map[string]int)
	for _, task := range AllTasks {
		_, exists := dup[task.ID]
		if exists {
			err = errors.New("Hay tareas duplicadas, se sale")
			log.Error("La tarea " + task.ID + " esta duplicada")
		} else {
			dup[task.ID]++
		}
	}
	return
}

//ValidateTaskDependencies Valida que las dependencias declaradas existan como tareas
func ValidateTaskDependencies(AllTasks []config.Task) (err error) {
	var exist bool
	for _, task := range AllTasks {
		for _, dep := range task.Depends {
			exist = false
			for _, task2 := range AllTasks {
				if task2.ID == dep {
					exist = true
					return
				}
			}
			if !exist {
				err = errors.New("Existen errores de dependencias inexistentes, se sale")
				log.Error("La dependencia ID:" + dep + " de la tarea: " + task.ID + " no existe")
			}
		}
	}
	return
}

//ValidateTaskCiclicDependencies Valida que entre dependencias declaradas no haya circulares
func ValidateTaskCiclicDependencies(AllTasks []config.Task) (err error) {
	for _, task := range AllTasks {
		for _, dep := range task.Depends {
			for _, task2 := range AllTasks {
				if task2.ID == dep {
					for _, dep2 := range task2.Depends {
						if dep2 == task.ID {
							err = errors.New("Existen errores de dependencias circulares, se sale")
							log.Error("Error de dependencias circulares entre las tareas " + task.ID + " y " + task2.ID)
						}
					}
				}
			}
		}
	}
	return
}
