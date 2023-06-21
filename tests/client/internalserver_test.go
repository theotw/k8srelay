// Copyright (c) 2023.  NetApp, Inc. All Rights Reserved.

package client

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHealthcheckURL(t *testing.T) {
	resp, err := http.Get("http://localhost:1701/healthcheck")
	if err != nil {
		t.Errorf("Unable to get healthcheck %s", err.Error())
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Healthcheck failed %d", resp.StatusCode)
		return
	}
}
func TestMertrics(t *testing.T) {
	urlToUse := "http://localhost:1701/metrics"
	var body []byte
	t.Run("Metrics Query", func(t2 *testing.T) {

		resp, err := http.Get(urlToUse)
		if err != nil {
			t.Errorf("Unable to get healthcheck %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			t.Errorf("Healthcheck failed %d", resp.StatusCode)
			return
		}

		body, _ = io.ReadAll(resp.Body)
	})

	t.Run("Metrics Values", func(t2 *testing.T) {
		//scan the data see if it has some metrics in it
		sbody := string(body)
		n := strings.Index(sbody, "process_virtual_memory_max_bytes")
		if n < 0 {
			t.Errorf("Unable to find process_virtual_memory_max_bytes")
			return
		}
		n = strings.Index(sbody, "promhttp_metric_handler_requests_total")
		if n < 0 {
			t.Errorf("Unable to find promhttp_metric_handler_requests_total")
			return
		}
	})
}
func TestAbout(t *testing.T) {
	resp, err := http.Get("http://localhost:1701/about")
	if err != nil {
		t.Errorf("Unable to get about %s", err.Error())
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("About failed %d", resp.StatusCode)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Errorf("About body is empty")
	}
}
