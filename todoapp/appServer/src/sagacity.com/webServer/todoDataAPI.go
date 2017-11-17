/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/webServer
Filename    : sagacity.com/webServer/todoDataAPI.go
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
	"encoding/json"

	"sagacity.com/appConstants"
	"sagacity.com/logger"
	"sagacity.com/models/dataCacheModels"
	"sagacity.com/models/webDataModels"
	"sagacity.com/dataCache"
)


// addtodorec
func addTodoRec(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /addtodorec.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_ADD_TODO_REC, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.AddJSONResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_ADD_TODO_REC,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecAddJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	todoCacheRec := dataCacheModels.TodoCacheRec {
		Title: jsonReq.Title,
		Status: jsonReq.Status,
		DueDate: jsonReq.DueDate,
		EDA: jsonReq.EDA,
		ShortDesc: jsonReq.ShortDesc,
	}

	isSuccess, id := dataCache.TodoCacheAddRec(&todoCacheRec)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Error while adding todo-record for title: %s", todoCacheRec.Title)
		result = false
		return
	}

	jsonResponse.ID = id
	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Added todo-record for title: %s", todoCacheRec.Title)
	return
}


// /gettodoreclist
func getTodoRecList(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /gettodoreclist.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_GET_TODO_RECLIST, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.TodoRecGetListJSONResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_GET_TODO_RECLIST,
	}

	isSuccess, todoCacheRecList := dataCache.TodoCacheGetList()
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Failed to fetch todo-rec list.")
		result = false
		return
	}

	todoRecList := make([]webDataModels.TodoRec, len(todoCacheRecList))
	for i, todoCacheRec := range todoCacheRecList {
		todoRecList[i].ID = todoCacheRec.ID
		todoRecList[i].Title = todoCacheRec.Title
		todoRecList[i].DueDate = todoCacheRec.DueDate
		todoRecList[i].EDA = todoCacheRec.EDA
		todoRecList[i].Status = todoCacheRec.Status
		todoRecList[i].CompletedOn = todoCacheRec.CompletedOn
		todoRecList[i].ShortDesc = todoCacheRec.ShortDesc
	}

	jsonResponse.TodoRecList = todoRecList
	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	return
}


// /gettodorec
func getTodoRec(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /gettodorec.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_GET_TODO_REC, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.TodoRecGetJSONResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_GET_TODO_REC,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecGetJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoRec := dataCache.TodoCacheGetRecCopy(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}

	jsonResponse.TodoRec = todoRec
	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Fetched todo-record for ID: %d", jsonReq.ID)
	return
}


// /updatetodorec
func updateTodoRec(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodorec.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.Title = jsonReq.TodoRec.Title
	todoCacheRec.DueDate = jsonReq.TodoRec.DueDate
	todoCacheRec.Status = jsonReq.TodoRec.Status
	todoCacheRec.EDA = jsonReq.TodoRec.EDA
	todoCacheRec.ShortDesc = jsonReq.TodoRec.ShortDesc

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated todo record successfully for ID: %d", jsonReq.ID)
	return
}


// /updatetodorectitle
func updateTodoRecTitle(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodorectitle.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC_TITLE, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC_TITLE,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateStrTypeJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.Title = jsonReq.StrVal

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated title of todo-record successfully for ID: %d", jsonReq.ID)
	return
}


// /updatetodorecshortdesc
func updateTodoRecShortDesc(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodorecshortdesc.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC_SHORTDESC, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC_SHORTDESC,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateStrTypeJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.ShortDesc = jsonReq.StrVal

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated short description of todo-record successfully for ID: %d", jsonReq.ID)
	return
}


// /updatetodorecduedate
func updateTodoRecDueDate(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodorecduedate.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC_DUEDATE, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC_DUEDATE,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateTimeTypeJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.DueDate = jsonReq.TimeVal

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated duedate of todo-record successfully for ID: %d", jsonReq.ID)
	return
}


// /updatetodoreceda
func updateTodoRecEDA(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodoreceda.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC_EDA, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC_EDA,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateTimeTypeJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.EDA = jsonReq.TimeVal

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated EDA of todo-record successfully for ID: %d", jsonReq.ID)
	return
}


// /updatetodorecstatus
func updateTodoRecStatus(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /updatetodorecstatus.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_UPDATE_TODO_REC_STATUS, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_UPDATE_TODO_REC_STATUS,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecUpdateStatusJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	isSuccess, todoCacheRec := dataCache.TodoCacheGetRec(jsonReq.ID)
	if !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "No todo item found for todoID: %d", jsonReq.ID)
		result = false
		return
	}
	defer dataCache.TodoCacheReleaseRec(todoCacheRec)

	todoCacheRec.Status = jsonReq.Status

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Updated status of todo-record successfully for ID: %d", jsonReq.ID)
	return
}


// /deletetodorec
func deleteTodoRec(responseWriter http.ResponseWriter, pReq *http.Request) {
	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Invoked /deletetodorec.")

	result := true

	defer func() {
		if !result {
			http.Error(responseWriter, appConstants.MSG_ERR_DELETE_TODO_REC, http.StatusInternalServerError)
		}
	}()

	jsonResponse := webDataModels.GenericResponse {
		Success: true,
		Message: appConstants.MSG_SUCCESS_DELETE_TODO_REC,
	}

	if pReq == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil request.")
		result = false
		return
	}

	if pReq.Body == nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Nil data in the request.")
		result = false
		return
	}

	jsonReq := webDataModels.TodoRecDeleteJSONRequest{}
	if err := json.NewDecoder(pReq.Body).Decode(&jsonReq); err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Decode error.")
		result = false
		return
	}

	if isSuccess := dataCache.TodoCacheDeleteRec(jsonReq.ID); !isSuccess {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Todo-record with ID: %d couldn't be removed", jsonReq.ID)
		result = false
		return
	}

	jsonResponseByteData, err := json.Marshal(jsonResponse)
	if err != nil {
		logger.Log(logger.WEBSERVICE, logger.ERROR, "Encode error while creating response data.")
		result = false
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonResponseByteData)

	logger.Log(logger.WEBSERVICE, logger.DEBUG, "Removed todo-record (%d) successfully.", jsonReq.ID)
	return
}
