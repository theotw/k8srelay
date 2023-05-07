package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

const (
	NATS_MSG_PREFIX      = "natssyncmsg"
	RequestForLocationID = "natssync.location.request"
	CLOUD_ID             = "cloud-master"

	// ResponseForLocationID this is the response subject, the data is the location ID, this message can be sent without a request, if the location ID changes
	ResponseForLocationID = "natssync.location.response"
	AppExitTopic          = "natssync.testing.exitapp"
)

func GetEnvWithDefaults(envKey string, defaultVal string) string {
	val := os.Getenv(envKey)
	if len(val) == 0 {
		val = defaultVal
	}
	return val
}
func GenerateUUID() string {
	uuid, _ := uuid.NewUUID()
	ret := uuid.String()
	return ret
}
func MakeReplySubject(replyToLocationID string) string {
	replySubject := fmt.Sprintf("%s.%s.%s", NATS_MSG_PREFIX, replyToLocationID, GenerateUUID())
	return replySubject
}
func MakeNBReplySubject() string {
	replySubject := fmt.Sprintf("%s.%s.%s", NATS_MSG_PREFIX, CLOUD_ID, GenerateUUID())
	return replySubject
}

func MakeMessageSubject(locationID string, params string) string {
	if len(params) == 0 {
		return fmt.Sprintf("%s.%s", NATS_MSG_PREFIX, locationID)
	}
	return fmt.Sprintf("%s.%s.%s", NATS_MSG_PREFIX, locationID, params)
}
