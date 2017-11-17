/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/models/webDataCommModels
Filename    : sagacity.com/models/webDataCommModels/webDataCommModels.go
File-type   : golang source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        :
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */
package webDataModels

import (
	"time"

	"sagacity.com/models/dataCacheModels"
)


// add todo-rec
type TodoRecAddJSONRequest struct {
	Title string `json:"title"`
	Status uint8 `json:"status"`	// 0: created, 1: assigned, 2: wip, 3: done
	DueDate time.Time `json:"duedate"`
	EDA time.Time `json:"eda"`
	ShortDesc string `json:"shortdesc"`
}
// used while sending responsed to add-record. web client should save the entity-id in its own cashe.
type AddJSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      uint32 `json:"todoid"`
}

// to fetch list of todo records.
type TodoRec struct {
	ID uint32 `json:"todoid"`
	Title string `json:"title"`
	Status uint8 `json:"status"`	// 0: created, 1: assigned, 2: wip, 3: done
	DueDate time.Time `json:"duedate"`
	EDA time.Time `json:"eda"`
	CompletedOn time.Time `json:"completedon"`
	ShortDesc string `json:"shortdesc"`
}
type TodoRecGetListJSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	TodoRecList []TodoRec `json"todoreclist"`
}

// to fetch a particular todo record.
type TodoRecGetJSONRequest struct {	// request
	ID uint32 `json:"todoid"`
}
type TodoRecGetJSONResponse struct {	// response
	Success bool   `json:"success"`
	Message string `json:"message"`
	TodoRec dataCacheModels.TodoCacheRec `json:"todorec"`
}

// to update todo record. response-type: GenericResponse
type TodoUpdateRec struct {
	Title string `json:"title"`
	Status uint8 `json:"status"`	// 0: created, 1: assigned, 2: wip, 3: done
	DueDate time.Time `json:"duedate"`
	EDA time.Time `json:"eda"`
	ShortDesc string `json:"shortdesc"`
}
type TodoRecUpdateJSONRequest struct {
	ID uint32 `json:"todoid"`
	TodoRec TodoUpdateRec `json:"todorec"`
}

// to update string type members. used for updating any of the members from: {Title, ShortDesc}
// response-type: GenericResponse 
type TodoRecUpdateStrTypeJSONRequest struct {
	ID uint32 `json:"todoid"`
	StrVal string `json:"strval"`
}

// to update time.Time type members. used for updating any of the members from: {DueDate, EDA}
// response-type: GenericResponse
type TodoRecUpdateTimeTypeJSONRequest struct {
	ID uint32 `json:"todoid"`
	TimeVal time.Time `json:"timeval"`
}

// to update status
// response-type: GenericResponse
type TodoRecUpdateStatusJSONRequest struct {
	ID uint32 `json:"todoid"`
	Status uint8 `json:"status"`	// 0: created, 1: assigned, 2: wip, 3: done
}

// to delete todo record. response-type: GenericResponse
type TodoRecDeleteJSONRequest struct {
	ID uint32 `json:"todoid"`
}

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
