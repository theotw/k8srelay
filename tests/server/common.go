/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"syscall"
)

func KillIt(w http.ResponseWriter, r *http.Request) {
	log.Infof("In Kill it, stopping test server")
	w.WriteHeader(200)
	killIt()
}
func Ready(w http.ResponseWriter, r *http.Request) {
	log.Infof("In ready check")
	w.Write([]byte("ready\n"))
	w.WriteHeader(200)
}
func killIt() {
	log.Infof("Sending sig int")
	raise(syscall.SIGINT)
	log.Infof("Done sending sig int")
}
func raise(sig os.Signal) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Errorf("Unable to get local process, bailing on the process %s", err.Error())
		os.Exit(1)
	}
	p.Signal(sig)
}
