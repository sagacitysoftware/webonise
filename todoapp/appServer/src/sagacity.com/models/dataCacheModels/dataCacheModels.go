/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/models/dataCacheModels
Filename    : sagacity.com/models/dataCacheModels/dataCacheModels.go
File-type   : golang source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        :
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */
package dataCacheModels

import (
	"time"
	"sync"
)

type TodoCacheRec struct {
	ID uint32 `json:"todoid"`
	Title string `json:"title"`
	Status uint8 `json:"status"`	// 0: created, 1: assigned, 2: wip, 3: done
	DueDate time.Time `json:"duedate"`
	EDA time.Time `json:"eda"`
	CompletedOn time.Time `json:"completedon"`
	ShortDesc string `json:"shortdesc"`
	RecLock *sync.Mutex    // record lock: a successful search through the cache returns locked record. any transaction on the record is mutually exclusive.
}
type TodoCacheMap map[uint32]*TodoCacheRec
