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
	err = pidfile.Write()
	if err != nil {
		fmt.Println("Error al crear PID: " + err.Error())
		os.Exit(1) //Salgo con error porque no pude ni leer la config
	}
}

func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	fmt.Println("Presiono CTRL-C; saliendooo")
	os.Exit(0)
}

func main() {
	go signalCatcher()
	for {
		Conf, err := config.ReadConfig()
		time.Sleep(time.Second * time.Duration(Conf.Global.CheckNewConfigInterval))
		for _, value := range Conf.Tasks {
			fmt.Println(value.ID)
			fmt.Println(value.Command)
		}
	}
}
