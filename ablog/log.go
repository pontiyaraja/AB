package ablog

import (
	"encoding/json"
	"fmt"
	"io"
	l "log"
	"os"
)

const (
	//INFO level 1
	INFO = iota
	//HTTP level 2
	HTTP
	//ERROR level 3
	ERROR
	//TRACE level 4
	TRACE
	//WARNING level 5
	WARNING
)

var (
	setLevel = WARNING
	trace    *l.Logger
	info     *l.Logger
	warning  *l.Logger
	httplog  *l.Logger
	errorlog *l.Logger
)

//LogDataMap map of key value pair to log
type LogDataMap map[string]interface{}

func init() {
	logInit(os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stderr)
}

func logInit(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	httpHandle io.Writer,
	errorHandle io.Writer) {

	trace = l.New(traceHandle,
		"TRACE|",
		l.LUTC|l.LstdFlags|l.Lshortfile)

	info = l.New(infoHandle,
		"INFO|",
		l.LUTC|l.LstdFlags|l.Lshortfile)

	warning = l.New(warningHandle,
		"WARNING|",
		l.LUTC|l.LstdFlags|l.Lshortfile)

	httplog = l.New(httpHandle,
		"HTTP|",
		l.LUTC|l.LstdFlags|l.Lshortfile)

	errorlog = l.New(errorHandle,
		"ERROR|",
		l.LUTC|l.LstdFlags|l.Lshortfile)
}

func doLog(cLog *l.Logger, level, callDepth int, v ...interface{}) {
	if level <= setLevel {
		if level == ERROR {
			cLog.SetOutput(os.Stderr)
		}
		cLog.Output(callDepth, fmt.Sprintln(v...))
	}
}

// HTTPLog prints the log in the following format:
//
// If any of the value is irrelevant then two consecutive PIPEs are printed:
// HTTP|TIMESTAMP|ServerIP:PORT|RequestMethod|RequestURL|ResponseStatusCode|ResponseWeight|UserAgent|Duration
func HTTPLog(logMessage string) {
	doLog(httplog, HTTP, 6, logMessage)
}

//Trace system gives facility to helps you isolate your system problems by monitoring selected events Ex. entry and exit
func traceLog(v ...interface{}) {
	doLog(trace, TRACE, 6, v...)
}

//Info dedicated for logging valuable information
func infoLog(v ...interface{}) {
	doLog(info, INFO, 6, v...)
}

//Warning for critical error
func warningLog(v ...interface{}) {
	doLog(warning, WARNING, 3, v...)
}

//Error logging error
func errorLog(v ...interface{}) {
	doLog(errorlog, ERROR, 6, v...)
}

func generateTrackingIDs(requestID string) string {
	var retString string
	if requestID != "" {
		retString += "requestId=" + requestID
	}
	return retString
}

//Error generates error log
func Error(requestID string, e error, data LogDataMap) {
	trackingIDs := generateTrackingIDs(requestID)
	msg := fmt.Sprintf("|%s|%s", trackingIDs, e.Error())
	if data != nil && len(data) > 0 {
		dataBytes, _ := json.Marshal(data)
		dataString := string(dataBytes)
		errorLog(msg, "|", dataString)
	} else {
		errorLog(msg)
	}
}

//Info generates info log
func Info(requestID, infoMessage string, data LogDataMap) {
	trackingIDs := generateTrackingIDs(requestID)
	dataBytes, _ := json.Marshal(data)
	dataString := string(dataBytes)
	msg := fmt.Sprintf("|%s|", trackingIDs)
	if data != nil && len(data) > 0 {
		infoLog(msg, infoMessage, "|", dataString)
	} else {
		infoLog(msg, infoMessage)
	}

}

//Warning generates warning log
func Warning(requestID, warnMessage string, data LogDataMap) {
	trackingIDs := generateTrackingIDs(requestID)
	msg := fmt.Sprintf("|%s|", trackingIDs)
	if data != nil && len(data) > 0 {
		dataBytes, _ := json.Marshal(data)
		dataString := string(dataBytes)
		warningLog(msg, warnMessage, "|", dataString)
	} else {
		warningLog(msg, warnMessage)
	}
}

//Trace generates trace log
func Trace(requestID, traceMessage string, data LogDataMap) {
	trackingIDs := generateTrackingIDs(requestID)
	msg := fmt.Sprintf("|%s|", trackingIDs)
	if data != nil && len(data) > 0 {
		dataBytes, _ := json.Marshal(data)
		dataString := string(dataBytes)
		traceLog(msg, traceMessage, "|", dataString)
	} else {
		traceLog(msg, traceMessage)
	}
}
