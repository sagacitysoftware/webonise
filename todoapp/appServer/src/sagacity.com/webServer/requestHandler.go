/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/webServer
Filename    : sagacity.com/webServer/requestHandler.go
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
	"net/http"

	"sagacity.com/logger"
)


func handleRequest(responseWriter http.ResponseWriter, pReq *http.Request) {
	switch pReq.Method {
		case "GET":
			break

		case "POST":
			switch pReq.RequestURI {
				case "/addtodorec":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					addTodoRec(responseWriter, pReq)
					break

				case "/gettodoreclist":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					getTodoRecList(responseWriter, pReq)
					break

				case "/gettodorec":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					getTodoRec(responseWriter, pReq)
					break

				case "/updatetodorec":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRec(responseWriter, pReq)
					break

				case "/updatetodorectitle":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRecTitle(responseWriter, pReq)
					break

				case "/updatetodorecshortdesc":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRecShortDesc(responseWriter, pReq)
					break

				case "/updatetodorecduedate":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRecDueDate(responseWriter, pReq)
					break

				case "/updatetodoreceda":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRecEDA(responseWriter, pReq)
					break

				case "/updatetodorecstatus":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					updateTodoRecStatus(responseWriter, pReq)
					break

				case "/deletetodorec":
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit addtodorec.")
					deleteTodoRec(responseWriter, pReq)
					break

				default:
					logger.Log(logger.WEBSERVICE, logger.DEBUG, "hit %s, incorrect URI", pReq.RequestURI)
					break
			}
			break

		default:
			logger.Log(logger.WEBSERVICE, logger.ERROR, "Incorrect method.")
			break
	}

	return
}
