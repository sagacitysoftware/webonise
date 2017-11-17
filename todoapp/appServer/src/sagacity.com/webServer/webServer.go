/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/webServer
Filename    : sagacity.com/webServer/webServer.go
File-type   : golang source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        :
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */
package webServer

import (
	"fmt"
	"sync"
	"net/http"

	"sagacity.com/appRepo"
	"sagacity.com/logger"
)

func Init(pVromWG *sync.WaitGroup) {
	defer pVromWG.Done()

	if appRepo.PConfigParameters == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "ERROR: Empty server configurations.\n")
		return
	}

	http.HandleFunc("/addtodorec", handleRequest)
	http.HandleFunc("/gettodoreclist", handleRequest)
	http.HandleFunc("/gettodorec", handleRequest)
	http.HandleFunc("/updatetodorec", handleRequest)
	http.HandleFunc("/updatetodorectitle", handleRequest)
	http.HandleFunc("/updatetodorecshortdesc", handleRequest)
	http.HandleFunc("/updatetodorecduedate", handleRequest)
	http.HandleFunc("/updatetodoreceda", handleRequest)
	http.HandleFunc("/updatetodorecstatus", handleRequest)
	http.HandleFunc("/deletetodorec", handleRequest)

	serverWebServicePort := fmt.Sprintf(":%d", appRepo.PConfigParameters.ServerConfigParams.ServerWebServicePort)
	http.ListenAndServe(serverWebServicePort, nil)
}
