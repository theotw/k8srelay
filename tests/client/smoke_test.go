// Copyright (c) 2023.  NetApp, Inc. All Rights Reserved.

package client

import (
	"context"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/theotw/k8srelay/pkg/natsmodel"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"os"
	"testing"
	"time"
)

const caCertReq = "k8srelay.cacert.req"
const caCertResp = "k8srelay.cacert.resp"

var caPubCert []byte

func TestBasics(t *testing.T) {
	natsURL := os.Getenv("NATS_SERVER_URL")
	err := natsmodel.InitNats(natsURL, "smoketest", time.Minute*2)
	if !assert.Nil(t, err) {
		t.Fatalf("Unable to init NATS %s", err.Error())
	}

	nc := natsmodel.GetNatsConnection()
	caSub, _ := nc.SubscribeSync(caCertResp)
	nc.Publish(caCertReq, nil)
	msg, err := caSub.NextMsg(2 * time.Minute)
	if err != nil {
		t.Fatalf("Unable to get CA Pub Cert %s", err.Error())
	}
	caSub.Unsubscribe()

	b64cert := string(msg.Data)
	certBits, err := base64.StdEncoding.DecodeString(b64cert)
	if err != nil {
		t.Fatalf("Unable to decode CA cert %s", err.Error())
	}
	caPubCert = certBits

	//bits, err := os.ReadFile("/Users/masonb/.kube/config")
	//if err != nil {
	//	t.Fatalf("oops")
	//}
	//configX, err := clientcmd.RESTConfigFromKubeConfig(bits)
	//if err != nil {
	//	t.Fatalf("oops")
	//}

	config := new(restclient.Config)
	config.Host = "https://localhost:8443"
	config.TLSClientConfig.CAData = caPubCert
	// Create a new clientset for interacting with the cluster
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create clientset: %v", err)
	}
	listNamespaces(t, clientset)
	listPods(t, clientset)
}

func listNamespaces(t *testing.T, client *kubernetes.Clientset) {

	// List namespaces
	ns, err := client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Greater(t, len(ns.Items), 0, "ns should not be 0")
}
func listPods(t *testing.T, client *kubernetes.Clientset) {
	pods, err := client.CoreV1().Pods("test").List(context.Background(), v1.ListOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Greater(t, len(pods.Items), 0, "pods should not be 0")
}
