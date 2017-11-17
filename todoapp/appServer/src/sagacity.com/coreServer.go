package main

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"
	"runtime"
	"sync"

	"sagacity.com/appRepo"
	"sagacity.com/config"
	"sagacity.com/logger"
	"sagacity.com/dataCache"
	"sagacity.com/webServer"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(2)

	appRepo.PConfigParameters = config.Init()
	logger.Init()

	go logger.LogDispatcher(&wg)

	if isSuccess := dataCache.Init(); !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Couldn't create user-cache.")
		os.Exit(1)
	}

	fmt.Printf("ServerWebServicePort: %d\n", appRepo.PConfigParameters.ServerConfigParams.ServerWebServicePort)

	go webServer.Init(&wg)

	logger.Log(logger.CORESERVER, logger.DEBUG, "Starting server.")

	time.Sleep(time.Second * 3)
	var strInput string
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nUser response:\n")
	for {
		fmt.Printf("Do you want to terminate the application: ")
		strInput, _ = reader.ReadString('\n')

		/* converts CRLF to LF */
		strInput = strings.Replace(strInput, "\n", "", -1)

		switch strInput {
			case "y", "Y", "yes", "Yes", "YES":
				logger.Log(logger.CORESERVER, logger.DEBUG, "Server terminating on user response: \"%s\".", strInput)
				fmt.Printf("Server terminating on user response: \"%s\"\n", strInput)
				time.Sleep(time.Second * 3)  /* waits for 3 seconds so that logger can push pending messages, though not guaranteed :). */
				os.Exit(0)
		}
	}
}
