package ablog

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestKIPError(t *testing.T) {
	requestID := uuid.New().String()
	err := errors.New("test error")
	var logData LogDataMap
	logData = make(LogDataMap)
	logData["error"] = err
	Error(requestID, err, logData)
}

func TestKIPWarning(t *testing.T) {
	requestID := uuid.New().String()
	var logData LogDataMap
	logData = make(LogDataMap)
	logData["warning message"] = "warning"
	Warning(requestID, "test warnign", logData)
}

func TestKIPInfo(t *testing.T) {
	requestID := uuid.New().String()
	info := "test info"
	var logData LogDataMap
	logData = make(LogDataMap)
	logData["info"] = info
	Info(requestID, info, logData)
}

func TestKIPTrace(t *testing.T) {
	requestID := uuid.New().String()
	trace := "test trace"
	var logData LogDataMap
	logData = make(LogDataMap)
	logData["trace"] = trace
	Trace(requestID, trace, logData)
}
