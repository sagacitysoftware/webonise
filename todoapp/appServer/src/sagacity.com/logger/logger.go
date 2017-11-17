/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Softwares Pvt. Ltd.. All rights reserved.
Package     : sagacity.com/logger
Filename    : sagacity.com/logger/logger.go
File-type   : golang source code file.

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        :
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
- logger library.
- Log API logger.Log() has been exposed.
- This API captures current timestamp, creates a string type variable, and dumps that variable in the
channel.
- Log message has following structure:
[module_name] [timestamp] [log_level] [path-of-source-file: lineno] [package_name.APIname]: log_message

Words enclosed in [] aren't optional. It's how they're placed in each log message. For instance:
[CORESERVER] [17-11-2017:07-37-01-89360970-IST] [DEBUG] [/sagacity.com/dataCache/todoCache.go: 211] [sagacity.com/dataCache.TodoCacheAddRec]:
dataCacheModels.TodoCacheRec todo-ID.name: 1.1st todo item added successfully.

- A dispatcher go routine fetches the channel, extracts the created message, and dumps the message in the log-file.
- Location of logfile is as mentioned in the appLinux.conf file.
- Max allowed size of a logfile is 20 MB (20971520 Bytes) and logfiles are rolled over after 10 log files.
**************************************************************************** */
package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"sync"

	"sagacity.com/appRepo"
)

const DEFAULTGOPATH string = "/appServer/src"

/* log levels */
const DEBUG string = "DEBUG"
const ERROR string = "ERROR"
const INFO string = "INFO"
const WARNING string = "WARNING"

var loglevel map[string]uint8

/* components */
const CORESERVER string = "CORESERVER"
const WEBSERVICE string = "WEBSERVICE"

type logmessage struct {
	componentFlag  int8
	component      string
	logmsg         string
}

// buffered channel with size 10.
var chanbuffLog chan logmessage

// log-file file handler.
var pServerLogFile *os.File

var currentLogfileCnt uint8 = 1
var logfileNameList []string
var dummyLogfile string


/* ****************************************************************************
Receiver    : na

Arguments   :
1> sourceFilePath string: Absolute path of source file where logger.Log() has been called from.
2> defaultPath string: Default path component.

Return value:
1> bool: true is successful, false otherwise.
2> string: Absolute-path less default path.

Description :
- Extracts sourceFilePath - defaultPath from sourceFilePath.

Additional note: na
**************************************************************************** */
func getFilePath(sourceFilePath string, defaultPath string) (bool, string) {
	filePath := ""
	if len(defaultPath) > len(sourceFilePath) {
		return false, filePath
	}

	length := len(sourceFilePath) - len(defaultPath)
	var i int
	for i = 0; i < length ; i++ {
		if sourceFilePath[i] == defaultPath[0] {
			if sourceFilePath[i : i + len(defaultPath)] == defaultPath {
				break
			}
		}
	}

	return true, sourceFilePath[i + len(defaultPath) : len(sourceFilePath)]
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> strcomponent string: Either of the following:
CORESERVER: for core server log message
WEBSERVICE: for webservice log message

2> loglevelStr string:
- There exist 4 loglevels: ERROR, WARNING, INFO, and DEBUG.
The loglevels are incremental where DEBUG being the highest one and
includes all log levels.

Return value: na

Description :
- Constructs a type logmessage variable.
- Dumps the same in the logmsg_buffered_channel

Additional note: na
**************************************************************************** */
func Log(strcomponent string, loglevelStr string, msg string, args ...interface{}) {
	if appRepo.PConfigParameters == nil {
		fmt.Printf("ERROR: Empty server configurations.\n")
		os.Exit(1)
	}

	configLoglevelVal := uint8(loglevel[appRepo.PConfigParameters.LogConfigParams.LogLevel]) /* 0: DEBUG, 1: INFO, 2: WARNING, 3: ERROR */
	msgLoglevelVal := loglevel[loglevelStr]
	if msgLoglevelVal < configLoglevelVal {
		return
	}

	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d-%02d-%02d-%06d-%s", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	pc, fn, line, _ := runtime.Caller(1)
	_, filePath := getFilePath(fn, DEFAULTGOPATH)

	logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n", strcomponent, msgTimeStamp, loglevelStr, filePath, line, runtime.FuncForPC(pc).Name())
	logMsg = fmt.Sprintf(logMsg+msg, args...)
	logMsg = logMsg + "\n"
	logMessage := logmessage{
		componentFlag: 0,
		component:      strcomponent,
		logmsg:         logMsg,
	}

	chanbuffLog <- logMessage
}

/* ****************************************************************************
Prototype   :
func LogDispatcher()

TODO: will get back to this signature once waitgroup are implemented in the coreserver.
func LogDispatcher(wg *sync.WaitGroup)

Arguments   : na for now.
1> wg *sync.WaitGroup: waitgroup handler for conveying done status to the caller.

Description :
- A go routine, invoked through Logger()
- Waits onto buffered channel name chanbuffLog infinitely.
- Extracts data from the channel, it's of type logmessage.
- Dumps log into the file pointed by pServerLogFile.

Assumptions :

TODO        :
db dispatch.

Return Value: na
**************************************************************************** */
/* func LogDispatcher(wg *sync.WaitGroup)	{ */
func LogDispatcher(pVromWG *sync.WaitGroup) {
	defer pVromWG.Done()

	for {
		select {
			case logMsg := <-chanbuffLog:	// pushes dummy logmessage onto the channel
				dumpServerLog(logMsg.logmsg)
		}
	}

}

/* ****************************************************************************
Prototype   :
func dumpServerLog(logMsg string)

Arguments   :
1> logMsg string: log message to be dumped in the logfile defined by appRepo.PConfigParameters.LogConf.LogFile

Description :
- Dumps logMsg into target logfile pointed to by plogfile file handler.
- Dumps logMsg into the database table.

Assumptions :

TODO        :
- dumping of logMsg into database table.

Return Value:
**************************************************************************** */
func dumpServerLog(logMsg string) {
	if pServerLogFile == nil {
		fmt.Println("Error: nil file handler.")
		os.Exit(1)
	}

	pServerLogFile.WriteString(logMsg)
	fmt.Printf(logMsg) /* TODO-REM: remove this fmp.Printf() call later */

	fi, err := pServerLogFile.Stat()
	if err != nil {
		fmt.Printf("Couldn't obtain stat, handle error: %s\n", err.Error())
		return
	}

	fileSize := fi.Size()
	if fileSize >= appRepo.PConfigParameters.LogConfigParams.LogFileSize {
		pServerLogFile.Close()
		pServerLogFile = nil
		err = os.Rename(logfileNameList[0], dummyLogfile)
		if err != nil {
			fmt.Println("error while mv %s to %s: ", logfileNameList[0], dummyLogfile, err.Error())
			pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR | os.O_APPEND | os.O_CREATE, 0644)
			return
		}

		pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR | os.O_APPEND | os.O_CREATE, 0644)
		if err != nil	{
			fmt.Printf("Error while recreating logfile: %s,  error: %s\n", logfileNameList[0], err.Error())
			return
		}

		if currentLogfileCnt < 10 {
			currentLogfileCnt = currentLogfileCnt + 1
		}

		go handleLogRotate()
	}
}

func handleLogRotate() {
	for i := currentLogfileCnt; i > 2; i-- {
		err := os.Rename(logfileNameList[i - 2], logfileNameList[i - 1])
		if err != nil {
			fmt.Println("error while mv %s to %s. error: %s\n", logfileNameList[i - 2], logfileNameList[i - 1], err.Error())
			return
		}
	}

	err := os.Rename(dummyLogfile, logfileNameList[1])
	if err != nil {
		fmt.Println("error while mv %s to %s. error: %s\n", dummyLogfile, logfileNameList[1], err.Error())
		return
	}
}

func Init() {
	if appRepo.PConfigParameters == nil {
		fmt.Printf("ERROR: Empty server configurations.\n")
		os.Exit(1)
	}

	logfileNameList = make([]string, appRepo.PConfigParameters.LogConfigParams.LogMaxFiles)
	dummyLogfile = appRepo.PConfigParameters.LogConfigParams.LogFile + ".1.dummy"

	chanbuffLog = make(chan logmessage, 10)
	var err error

	logFile := appRepo.PConfigParameters.LogConfigParams.LogFile + ".1"
	pServerLogFile, err = os.OpenFile(logFile, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error while creating logfile: %s, Error: %s\n", logFile, err.Error())
		os.Exit(1)
	}

	loglevel = make(map[string]uint8)
	loglevel["DEBUG"] = uint8(0)
	loglevel["INFO"] = uint8(1)
	loglevel["WARNING"] = uint8(2)
	loglevel["ERROR"] = uint8(3)

	for i := int8(0); i < appRepo.PConfigParameters.LogConfigParams.LogMaxFiles; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", appRepo.PConfigParameters.LogConfigParams.LogFile, i + 1)
	}

	/*
	for i := range logfileNameList {
		fmt.Printf("[dbgrm]:  %d: logfile-name: %s\n", i, logfileNameList[i])
	}
	*/
}
