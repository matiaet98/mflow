#!/bin/bash

#path absoluto hacia ejecutable de mflow, ej: /usr/local/mflow/mflow
mflow="/home/matiaet98/go/src/mflow/mflow"
#path absoluto hacia el pid de mflow, ej: /usr/local/mflow/mflow.pid
PIDFILE="/home/matiaet98/go/src/mflow/mflow.pid"

if [[ $USER != "matiaet98" ]]; then
    echo "Se debe correr con el usuario hadoop"
    exit 1
fi

start(){
    if [ -f $PIDFILE ] && kill -0 $(cat $PIDFILE); then
        echo 'El servicio ya esta corriendo' >&2
        return 1
    fi
    echo 'Iniciando servicio…' >&2
    $mflow &>/dev/null & echo $! > $PIDFILE
    echo 'Servicio iniciado' >&2
    if [[ ! -f $PIDFILE ]]; then
        echo "El proceso no inicio correctamente, revise logs"
    fi
}

stop(){
    if [ ! -f $PIDFILE ] || ! kill -0 $(cat "$PIDFILE"); then
        echo 'El servicio no esta corriendo' >&2
        return 1
    fi
    echo 'Finalizando servicio…' >&2
    kill -SIGTERM $(cat "$PIDFILE") && rm -f "$PIDFILE"
    echo 'Servicio finalizado' >&2
}

status(){
    if [ -f $PIDFILE ] && kill -0 $(cat $PIDFILE); then
        echo 'El servicio se encuentra corriendo' >&2
        return 1
    else
        echo "El proceso no esta corriendo"
    fi
}

case "$1" in
   start) start ;;
   stop)  stop;;
   status) status;;
   retart)
    stop
    start
    ;;
   *) echo "usage $0 start|stop|restart|status" >&2
      exit 1
    ;;
esac