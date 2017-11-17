/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/dataCache
Filename    : sagacity.com/dataCache/dataCacheRepo.go
File-type   : golang source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        :
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version.
- Defines todo-cache.
- Defines todo-cache store locks.
**************************************************************************** */
package dataCache

import (
	"sync"

	"sagacity.com/models/dataCacheModels"
)


// defines todo-cache.
var todoCache dataCacheModels.TodoCacheMap

/* toco-cache store-locks.
wr store lock is mutually exclusive.
rd store lock is mutually inclusive only for other rd store locks. */
var todoCacheLock sync.RWMutex

var idAutoIncr uint32
