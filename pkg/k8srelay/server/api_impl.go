/*
 * Copyright (c) The One True Way 2021. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/theotw/k8srelay/pkg"
	"github.com/theotw/k8srelay/pkg/k8srelay/model"
	"github.com/theotw/k8srelay/pkg/natsmodel"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JsonError struct {
	Message string `json:"message"`
}

func HandleError(_ *gin.Context, err error) (int, interface{}) {
	ret := JsonError{
		Message: err.Error(),
	}

	return http.StatusInternalServerError, &ret
}

func RouteHandlerForInternalAPI(c *gin.Context) {
	//normalize out the string
	tmp := c.Request.RequestURI
	if !strings.HasSuffix(tmp, "/") {
		tmp = fmt.Sprintf("%s/", tmp)
	}

	if strings.HasSuffix(tmp, "/about/") {
		aboutGetUnversioned(c)
		return
	}

	if strings.HasSuffix(tmp, "/healthcheck/") {
		healthCheckGetUnversioned(c)
		return
	}

	if strings.HasSuffix(tmp, "/metrics/") {
		metricsHandler(c)
		return
	}

	c.JSON(http.StatusServiceUnavailable, "")
	return

}

type AboutResponse struct {
	AppVersion string `json:"appVersion"`
}

func aboutGetUnversioned(c *gin.Context) {
	var resp AboutResponse
	resp.AppVersion = pkg.VERSION

	c.JSON(http.StatusOK, &resp)
}

func healthCheckGetUnversioned(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func metricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

const bearer = "Bearer "

func genericHandlerHandler(c *gin.Context) {
	parse, err := url.Parse(c.Request.RequestURI)
	userTokenWhichBecomesRouteID := GetRouteIDFromAuthHeader(c)
	log.Infof(userTokenWhichBecomesRouteID)
	log.Infof("URI %s", c.Request.URL.String())
	if err != nil {
		panic(err)
	}

	req := model.NewCallReq()
	for k, v := range c.Request.Header {
		if k != "Authorization" {
			req.AddHeader(k, v[0])
		}
	}
	req.QueryString = c.Request.URL.RawQuery
	req.Path = parse.Path
	req.Method = c.Request.Method
	bodyBits, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Errorf("Unable to read body on call %s - %s error: %s", c.Request.Method, parse.Path, err.Error())
	}

	requestUUID := uuid2.New().String()
	urlString := c.Request.URL.String()
	if strings.HasPrefix(urlString, "/api/v1/namespaces") &&
		strings.Contains(urlString, "pods") &&
		strings.Contains(urlString, "log") &&
		strings.HasSuffix(urlString, "follow=true") {
		// looking for log follow /api/v1/namespaces/<ns>/pods/<pod>/log?container=<container>&follow=true
		log.Infof("got log streaming request, setting stream and UUID %s for the client request", requestUUID)
		req.Stream = true
		req.UUID = requestUUID
	}
	req.InBody = bodyBits
	nc := natsmodel.GetNatsConnection()
	replySub := model.MakeNBReplySubject()
	sbMsgSub := model.MakeMessageSubject(userTokenWhichBecomesRouteID, model.K8SRelayRequestMessageSubjectSuffix)
	bits, err := json.Marshal(&req)
	if err != nil {
		c.Status(502)
		log.WithError(err).Errorf("Returning a 502, got an error Marshal %s ", err.Error())
		c.Header("Content-Type", "text/plain")
		c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
		return
	}
	nm := nats.NewMsg(sbMsgSub)
	nm.Data = bits
	nm.Reply = replySub
	replyChannel, err := nc.SubscribeSync(replySub)
	if err != nil {
		c.Status(502)
		log.WithError(err).Errorf("Returning a 502, got an error Subscribe %s ", err.Error())
		c.Header("Content-Type", "text/plain")
		c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
		return
	}
	err = nc.PublishMsg(nm)
	if err != nil {
		c.Status(502)
		c.Header("Content-Type", "text/plain")
		log.WithError(err).Errorf("Returning a 502, got an error failed to publish message %s ", err.Error())
		c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
		return
	}

	isFirst := true
	for {
		select {
		case <-c.Request.Context().Done():
			log.Info("context done, client might have disconnected, returning")
			if req.Stream {
				endLogStreaming(c, nc, requestUUID)
			}
			return
		default:
			msg, err := replyChannel.NextMsg(time.Minute * 2)
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}
				c.Status(502)
				c.Header("Content-Type", "text/plain")
				log.WithError(err).Errorf("Returning a 502, got an error next message %s ", err.Error())
				c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
				if req.Stream {
					endLogStreaming(c, nc, requestUUID)
				}
				return
			}

			var respMsg model.CallResponse
			err = json.Unmarshal(msg.Data, &respMsg)
			if err != nil {
				if req.Stream {
					endLogStreaming(c, nc, requestUUID)
				}
				c.Status(502)
				c.Header("Content-Type", "text/plain")
				log.WithError(err).Errorf("Returning a 502, got an error on unmarshall %s ", err.Error())
				c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
				return
			}
			if isFirst {
				log.Infof("Got resp status %d ", respMsg.StatusCode)
				for k, v := range respMsg.Headers {
					log.Infof("%s = %s ", k, v)
					c.Header(k, v)
				}
				c.Status(respMsg.StatusCode)
				isFirst = false
			}

			if respMsg.OutBody != nil {
				_, err = c.Writer.Write(respMsg.OutBody)
				if err != nil {
					if err.Error() == "client disconnected" && req.Stream {
						log.Warn("Write err, client disconnected, ending streaming")
						endLogStreaming(c, nc, requestUUID)
						return
					}
				}
				c.Writer.Flush()
			}
			if respMsg.LastMessage {
				return
			}
		}
	}
}

func GetRouteIDFromAuthHeader(c *gin.Context) string {
	var userTokenWhichBecomesRouteID string
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, bearer) {
		userTokenWhichBecomesRouteID = token[len(bearer):]
	} else {
		userTokenWhichBecomesRouteID = "dev"
	}
	return userTokenWhichBecomesRouteID
}

func endLogStreaming(c *gin.Context, nc *nats.Conn, requestUUID string) {
	log.Infof("ending log streaming")
	userTokenWhichBecomesRouteID := GetRouteIDFromAuthHeader(c)
	sbMsgSub := model.MakeMessageSubject(userTokenWhichBecomesRouteID, model.K8SRelayRequestMessageSubjectSuffix+"."+requestUUID+".stopStreaming")
	log.Infof("endLogStreaming: subject for log streaming end: %s", sbMsgSub)
	nm := nats.NewMsg(sbMsgSub)
	err := nc.PublishMsg(nm)
	if err != nil {
		c.Status(502)
		c.Header("Content-Type", "text/plain")
		log.WithError(err).Errorf("Returning a 502, got an error failed to publish message %s ", err.Error())
		c.Writer.Write([]byte(fmt.Sprintf(" gate way error %s", err.Error())))
		return
	}
}
