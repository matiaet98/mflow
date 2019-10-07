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
	if err != nil {
		log.Panicln(err)
	}
	return string(out), nil
}

// Run : Corre un proceso oracle
func (ps OracleProcess) Run() (string, error) {
	db, err := sql.Open("goracle", ps.User+"/"+ps.Password+"@"+ps.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	var output string
	_, err = db.Exec(ps.Command, sql.Named("respuesta", sql.Out{Dest: &output}))
	if err != nil {
		log.Panicln(err)
	}
	return output, nil
}
