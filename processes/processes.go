package processes

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	_ "gopkg.in/goracle.v2" //se abstrae su uso con la libreria sql
	"os/exec"
)

// Process : Interface para correr procesos
type Process interface {
	Run() (string, error)
}

// BashProcess : Para procesos en bash
type BashProcess struct {
	Command string
}

// OracleProcess : Para procesos en oracle
type OracleProcess struct {
	Command          string
	ConnectionString string
	User             string
	Password         string
}

// Run : Corre un proceso bash
func (ps BashProcess) Run() (string, error) {
	cmd := exec.Command(ps.Command)
	out, err := cmd.CombinedOutput() //este chabon aparte de combinar stderr y stdout tambien hace el Run... poco intuitivo
	return string(out), err
}

// Run : Corre un proceso oracle
func (ps OracleProcess) Run() (string, error) {
	db, err := sql.Open("goracle", ps.User+"/"+ps.Password+"@"+ps.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	_, err = db.Exec(ps.Command)
	return "", err
}
