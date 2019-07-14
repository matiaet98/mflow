package processes

import (
	"fmt"
	"os/exec"
)

// Process : Interface para correr procesos
type Process interface {
	Run() string
}

// BashProcess : Para procesos en bash
type BashProcess struct {
	Command string
}

// Run : Corre un proceso bash
func (ps BashProcess) Run() string {
	cmd := exec.Command(ps.Command)
	out, err := cmd.CombinedOutput() //este chabon aparte de combinar stderr y stdout tambien hace el Run... poco intuitivo
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}
