package tasks

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"mflow/config"
	"mflow/global"
	"mflow/processes"
	"os"
	"strconv"
	"time"
)

const runningStatus string = "RUNNING"
const noneStatus string = "NONE"
const failedStatus string = "FAILED"
const successStatus string = "SUCCESS"

func runTask(task config.Task, sem chan bool) {
	switch task.Type {
	case "bash":
		runBash(task, sem)
		break
	case "oracle":
		runOracle(task, sem)
		break
	}
}

func runBash(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var output string
	var err error
	var ps processes.Process
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.Name + ".log")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = f1
	ps = processes.BashProcess{Command: task.Command}
	output, err = ps.Run()
	logger.Println(output)
	log.Infoln(output)
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		log.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
}

func runOracle(task config.Task, sem chan bool) {
	defer func() { <-sem }()
	var output string
	var err error
	var ps processes.Process
	f1, err := os.Create(config.Config.LogDirectory + "master_" + strconv.Itoa(global.IDMaster) + "_task_" + task.Name + ".log")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	logger.Out = f1
	conn := getConnection(task.Db)
	ps = processes.OracleProcess{
		User:             conn.User,
		Password:         conn.Password,
		ConnectionString: conn.ConnectionString,
		Command:          task.Command}
	output, err = ps.Run()
	logger.Println(output)
	log.Infoln(output)
	if err != nil {
		setTaskStatus(task.ID, failedStatus)
		log.Warnln(err)
		return
	}
	setTaskStatus(task.ID, successStatus)
	return
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
	fisco := getConnection("fisco")

	db, err := sql.Open("goracle", fisco.User+"/"+fisco.Password+"@"+fisco.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	var command string = `declare
		l_id number;
		begin
		mflow.pkg_taskman.create_master(:l_id);
		end;
	`
	_, err = tx.Exec(command, sql.Named("l_id", sql.Out{Dest: &global.IDMaster}))
	if err != nil {
		log.Fatalln(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	tx, err = db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	command = `declare
		l_id number;
		begin
	`
	for x := range config.Config.Tasks.Tasks {
		command = fmt.Sprintf("%s mflow.pkg_taskman.create_task(%v,%v);\n", command, global.IDMaster, config.Config.Tasks.Tasks[x].ID)
	}
	command = fmt.Sprintf("%s end;\n", command)
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
	fisco := getConnection("fisco")
	db, err := sql.Open("goracle", fisco.User+"/"+fisco.Password+"@"+fisco.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Panicln(err)
	}
	const command string = `
		begin
		mflow.pkg_taskman.end_master(:id_master);
		end;
	`
	_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster))
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
	fisco := getConnection("fisco")
	db, err := sql.Open("goracle", fisco.User+"/"+fisco.Password+"@"+fisco.ConnectionString)
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

func setTaskStatus(IDTask int, status string) (string, error) {
	fisco := getConnection("fisco")
	db, err := sql.Open("goracle", fisco.User+"/"+fisco.Password+"@"+fisco.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	command := `
	begin
	mflow.pkg_taskman.update_task(:id_master, :id_task,:status);
	end;
	`
	_, err = tx.Exec(command, sql.Named("id_master", global.IDMaster), sql.Named("id_task", IDTask), sql.Named("status", status))
	if err != nil {
		log.Fatalln(err)

	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return status, nil
}
