package main

import (
	"fmt"
	"github.com/facebookgo/pidfile"
	"os"
	"os/signal"
	"siper/config"
	"syscall"
	"time"
)

//Config : Global donde se guarda la configuracion
var Config config.Conf

const pidFilename string = "siper.pid"

func init() {
	Config, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error: Revise la configuracion")
		os.Exit(1) //Salgo con error porque no pude ni leer la config
	}
	pidfile.SetPidfilePath(Config.Global.PidDirectory + pidFilename)
	pid, err := pidfile.Read()
	if err == nil { //Significa que ya existe un PID
		fmt.Println("Error: El proceso ya esta corriendo (ya existe un PID)")
		fmt.Println("PID: ", pid)
		os.Exit(1)
	}
	err = pidfile.Write()
	if err != nil {
		fmt.Println("Error al crear PID: " + err.Error())
		os.Exit(1)
	}
}

func salir() {
	err := os.Remove(pidfile.GetPidfilePath())
	if err != nil {
		fmt.Println("Error: Error removiendo pid file: ", err.Error())
	}
	os.Exit(0)
}

func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	fmt.Println("Presiono CTRL-C; saliendooo")
	salir()
}

func main() {
	go signalCatcher()
	for {
		Conf, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error: Revise la configuracion")
			os.Exit(1)
		}
		time.Sleep(time.Second * time.Duration(Conf.Global.CheckNewConfigInterval))
		for _, value := range Conf.Tasks {
			fmt.Println(value.ID)
			fmt.Println(value.Command)
		}
	}
}
