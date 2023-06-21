/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/theotw/k8srelay/pkg/k8srelay/server"
	"net/http"
	"os"
	"testing"
)

func TestServerMain(t *testing.T) {
	fmt.Println("TestServerMain")
	working := os.Getenv("WORKING_DIR")
	if len(working) > 0 {
		os.Chdir(working)
	}
	t.Log(os.Getwd())

	t.Log("Test")
	os.Setenv("LOG_LEVEL", "debug")
	log.Info("Starting Relay Server")
	http.HandleFunc("/kill", KillIt)
	http.HandleFunc("/ready", Ready)
	go http.ListenAndServe(":8080", nil)
	srv, err := server.NewServer()
	if err != nil {
		log.Errorf("Unable to create server %s", err.Error())
		os.Exit(1)
	}
	fmt.Println("Running TestServerMain")
	srv.RunRelayServer(true)
	t.Log("Test 2")

	//	runtime.Goexit()
}
