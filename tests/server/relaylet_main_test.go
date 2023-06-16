/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/theotw/k8srelay/pkg/k8srelay/relaylet"
	"github.com/theotw/k8srelay/pkg/natsmodel"
	testing2 "github.com/theotw/k8srelay/pkg/testing"
	"net/http"
	"os"
	"os/signal"
	"testing"
)

func TestRelayletMain(t *testing.T) {
	os.Setenv("LOG_LEVEL", "trace")
	log.Info("Starting Relaylet")
	http.HandleFunc("/kill", KillIt)
	http.HandleFunc("/ready", Ready)
	go http.ListenAndServe(":8080", nil)
	_, err := relaylet.NewRelaylet()
	if err != nil {
		log.Errorf("Unable to create server %s", err.Error())
		os.Exit(1)
	}

	con := natsmodel.GetNatsConnection()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	testing2.NotifyOnAppExitMessageGenericNats(con, quit)
	log.Info("Server Started blocking on channel")

	//wait for the signal to exit, REST or nats message
	<-quit

	log.Info("Shutdown Relaylet ...")

	log.Info("Relaylet exiting")
}
