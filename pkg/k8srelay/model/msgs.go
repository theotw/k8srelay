/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package model

import (
	"fmt"
	"github.com/google/uuid"
)

const K8SRelayRequestMessageSubjectSuffix = "k8s-relay-req"
const NatssyncMessagePrefix = "natssyncmsg"
const CloudID = "cloud-master"

type CallRequest struct {
	// Path is the path part of the URL for this call
	Path string `json:"path"`

	Headers map[string]string `json:"headers"`

	// Method the HTTP method to perform
	Method string `json:"method"`
	// InBody is the input body for the call, which may be nil
	InBody      []byte `json:"inBody,omitempty"`
	QueryString string `json:"queryString"`
	Stream      bool   `json:"stream,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

func NewCallReq() *CallRequest {
	x := new(CallRequest)
	x.Headers = make(map[string]string, 0)
	return x
}
func (t *CallRequest) AddHeader(k, v string) {
	t.Headers[k] = v
}

type CallResponse struct {
	// Path is the path part of the URL for this call
	Path string `json:"path"`

	// Headers.  HTTP Headers, only set on the first response message on a multi message response
	Headers map[string]string `json:"headers"`

	StatusCode int `json:"statusCode"`
	// InBody is the input body for the call, which may be nil
	OutBody []byte `json:"inBody,omitempty"`

	// LastMessage indicates it the final in a multi message response
	LastMessage bool `json:"lastMessage"`
}

func NewCallResponse() *CallResponse {
	x := new(CallResponse)
	x.Headers = make(map[string]string, 0)
	return x
}
func (t *CallResponse) AddHeader(k, v string) {
	t.Headers[k] = v
}
func MakeReplySubject(replyToLocationID string) string {
	replySubject := fmt.Sprintf("%s.%s.%s", NatssyncMessagePrefix, replyToLocationID, GenerateUUID())
	return replySubject
}
func MakeNBReplySubject() string {
	replySubject := fmt.Sprintf("%s.%s.%s", NatssyncMessagePrefix, CloudID, GenerateUUID())
	return replySubject
}

func MakeMessageSubject(locationID string, params string) string {
	if len(params) == 0 {
		return fmt.Sprintf("%s.%s", NatssyncMessagePrefix, locationID)
	}
	return fmt.Sprintf("%s.%s.%s", NatssyncMessagePrefix, locationID, params)
}
func GenerateUUID() string {
	x, _ := uuid.NewUUID()
	return x.String()
}
