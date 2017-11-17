/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software. All rights reserved.
Package     : sagacity.com/dataCache
Filename    : sagacity.com/dataCache/todoCache.go 
File-type   : golang source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Date        : 13-Nov-2017-210418-678-IST
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
- Defines todo-cache APIs
**************************************************************************** */
package dataCache

import (
	"sync"
	"time"

	"sagacity.com/logger"
	"sagacity.com/models/dataCacheModels"
)


/* ****************************************************************************
Receiver    : na

Arguments   : na

Return value: na

Description : API to acquire RD store lock over todo-cache.

Additional note: na
**************************************************************************** */
func TodoCacheRDLock() {
	todoCacheLock.RLock()
}


/* ****************************************************************************
Receiver    : na

Arguments   : na

Return value: na

Description : API to release RD store lock over todo-cache.

Additional note: na
**************************************************************************** */
func TodoCacheRDUnlock() {
	todoCacheLock.RUnlock()
}


/* ****************************************************************************
Receiver    : na

Arguments   : na

Return value: na

Description : API to acquire WR store lock over todo-cache.

Additional note: na
**************************************************************************** */
func TodoCacheWRLock() {
	todoCacheLock.Lock()
}


/* ****************************************************************************
Receiver    : na

Arguments   : na

Return value: na

Description : API to release WR store lock over todo-cache.

Additional note: na
**************************************************************************** */
func TodoCacheWRUnlock() {
	todoCacheLock.Unlock()
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> pRec *dataCacheModels.TodoCacheRec: reference to todo-cache record

Return value: na

Description :
- Releases locked record. Nulls the cache record so that it can be protected from accidental access.

Additional note: na
**************************************************************************** */
func TodoCacheReleaseRec(pRec *dataCacheModels.TodoCacheRec) {
	if (pRec != nil) && (pRec.RecLock != nil) {
		pRec.RecLock.Unlock() // releases record lock
		pRec = nil
	}
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> pTmpRec *dataCacheModels.TodoCacheRec: Source todo-cache record.

Return value:
1> *dataCacheModels.TodoCacheRec: Newly created todo-cache record.

Description :
- Creates a todo-cache record dynamically.
- Creates mutex for record locking also.

Additional note: na
**************************************************************************** */
func TodoCacheCreateRecFromStruct(pTmpRec *dataCacheModels.TodoCacheRec) *dataCacheModels.TodoCacheRec {
	pRec := &dataCacheModels.TodoCacheRec {
		Title: pTmpRec.Title,
		DueDate: pTmpRec.DueDate,
		EDA: pTmpRec.EDA,
		ShortDesc: pTmpRec.ShortDesc,
	}
	pRec.RecLock = &sync.Mutex{}

	return pRec
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> tmpRecList []dataCacheModels.TodoCacheRec): Slice of todo-cache records.

Return value:
1> bool: true for success. false for failure.

Description :
- Typically used at server start-up where todo-cache is populated from the result of a db query.
- Receives slice of todo-cache records, each record in the slice thereby goes in the todo-cache.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheAddRecFromList(tmpRecList []dataCacheModels.TodoCacheRec) bool {
	TodoCacheWRLock()
	defer TodoCacheWRUnlock()

	for _, todoCacheRec := range tmpRecList {
		idAutoIncr += 1
		pRec := &dataCacheModels.TodoCacheRec {
			ID: idAutoIncr,
			Title: todoCacheRec.Title,
			DueDate: todoCacheRec.DueDate,
			ShortDesc: todoCacheRec.ShortDesc,
		}
		pRec.RecLock = &sync.Mutex{}

		todoCache[idAutoIncr] = pRec
	}

	return true
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> pTmpRec *dataCacheModels.TodoCacheRec: Source record of type TodoCacheRec, to be added to todo-cache.

Return value:
1> bool: true for success. false for failure.

Description :
- Adds todo-cache record to todo-cache based on todo-ID key.
- Creates todo-cache record dynamically from pTmpRec.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheAddRec(pTmpRec *dataCacheModels.TodoCacheRec) (bool, uint32) {
	TodoCacheWRLock()
	defer TodoCacheWRUnlock()

	pRec, ok := todoCache[pTmpRec.ID]
	if ok == true {
		logger.Log(logger.CORESERVER, logger.WARNING, "dataCacheModels.TodoCacheRec todo-ID.name: %d.%s already exists.", pRec.ID, pRec.Title)
		return true, 0
	}

	pRec = TodoCacheCreateRecFromStruct(pTmpRec)
	idAutoIncr += 1
	pRec.ID = idAutoIncr
	todoCache[idAutoIncr] = pRec
	logger.Log(logger.CORESERVER, logger.DEBUG, "dataCacheModels.TodoCacheRec todo-ID.name: %d.%s added successfully.", pRec.ID, pRec.Title)

	return true, pRec.ID
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> ID uint32: ID of the todo-record whose details need to be fetched.

Return value:
1> bool: true if record is found. false if it'sn't.
2> *dataCacheModels.TodoCacheRec: Fetched todo-cache record.

Description :
- Function fetches records from todo-cache based on todo-ID as a key.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function. Locking sequence, however, should be maintained.
- Fetched record is in the locked state. It's the caller's responsibility to release the locked record.
**************************************************************************** */
func TodoCacheGetRec(todoID uint32) (bool, *dataCacheModels.TodoCacheRec) {
	TodoCacheRDLock()
	defer TodoCacheRDUnlock()

	pRec, isOK := todoCache[todoID]
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todoCache-record with ID: %d doesn't exist.", todoID)
		return false, nil
	}

	pRec.RecLock.Lock()  // record is locked

	return true, pRec
}

func TodoCacheGetRecCopy(todoID uint32) (bool, dataCacheModels.TodoCacheRec) {
	isOK, pRec := TodoCacheGetRec(todoID)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todo-cache record for ID %d doesn't exist.", todoID)
		return false, dataCacheModels.TodoCacheRec{}
	}

	rec := *pRec
	TodoCacheReleaseRec(pRec)

	return true, rec
}


/* ****************************************************************************
Receiver    : na

Arguments   : na

Return value: na

Description :
- Function lists all the records in todo-cache.
- This's an iterator. So, it runs in the WR store lock over the cache.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheGetList() (bool, []dataCacheModels.TodoCacheRec) {
	TodoCacheWRLock()
	defer TodoCacheWRUnlock()

	logger.Log(logger.CORESERVER, logger.DEBUG, "todo-cache-len: %d", len(todoCache))

	if len(todoCache) < 1 {
		logger.Log(logger.CORESERVER, logger.WARNING, "empty todo-cache")
		return false, []dataCacheModels.TodoCacheRec{}
	}

	todoCacheRecList := make([]dataCacheModels.TodoCacheRec, len(todoCache))
	if todoCacheRecList == nil {
		logger.Log(logger.CORESERVER, logger.ERROR, "Failed to create todo-cache reclist.")
		return false, []dataCacheModels.TodoCacheRec{}
	}

	i := 0
	for key, pRec := range todoCache {
		logger.Log(logger.CORESERVER, logger.DEBUG, "[key: %d]: todo-cache-rec: %#v", key, *pRec)
		todoCacheRecList[i] = *pRec
		i++
	}

	return true, todoCacheRecList
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> pTmpRec *dataCacheModels.TodoCacheRec: Source todo-cache record.

Return value:
1> bool: true for success. false for failure.

Description :
- Function fetches todo-cache record based on ID key and updates the same
from the source todo-cache record.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheUpdateRec(pTmpRec *dataCacheModels.TodoCacheRec) bool {
	isOK, pRec := TodoCacheGetRec(pTmpRec.ID)
	TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todoCache record for ID: %d doesn't exist.", pTmpRec.ID)
		return false
	}

	pRec.Title = pTmpRec.Title
	pRec.DueDate = pTmpRec.DueDate
	pRec.Status = pTmpRec.Status
	pRec.CompletedOn = pTmpRec.CompletedOn
	pRec.ShortDesc = pTmpRec.ShortDesc

	logger.Log(logger.CORESERVER, logger.DEBUG, "dataCacheModels.TodoCacheRec todoID: %d updated successfully.", pTmpRec.ID)

	return true
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> newTitle string: New title for the Todo-cache record.
2> id uint32: Todo-cache record ID whose EDA needs to be updated.

Return value:
1> bool: true if Title is updated successfully. false otherwise.

Description :

Additional note:
**************************************************************************** */
func TodoCacheUpdateTitle(newTitle string, id uint32) bool {
	isOK, pRec := TodoCacheGetRec(id)
	TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todoCache record for ID: %d doesn't exist.", id)
		return false
	}

	pRec.Title = newTitle
	logger.Log(logger.CORESERVER, logger.ERROR, "Title of todo-cache record updated. ID.Title :: %d.%s.", id, newTitle)

	return true
}


/* ****************************************************************************
Receiver    :

Arguments   :
1> newShortDesc string: New ShortDesc for the todo-cache record.
2> id uint32: Todo-cache record ID whose EDA needs to be updated.

Return value:

Description :

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheUpdateShortDesc(newShortDesc string, id uint32) bool {
	isOK, pRec := TodoCacheGetRec(id)
	TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todo-cache record for ID: %d doesn't exist.", id)
		return false
	}

	pRec.ShortDesc = newShortDesc
	logger.Log(logger.CORESERVER, logger.ERROR, "ShortDesc of todo-cache record updated. ID.ShortDesc :: %d.%s.", id, newShortDesc)

	return true
}



/* ****************************************************************************
Receiver    :

Arguments   :
1> newDueDate time.Time: new DueDate.
2> id uint32: Todo-cache record ID whose EDA needs to be updated.

Return value:
1> bool: true if DueDate is updated successfully. false otherwise.

Description :
- Function updates DueDate of the todo-cache record. Todo-cache record is fetched from the given id.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheUpdateDueDate(newDueDate time.Time, id uint32) bool {
	isOK, pRec := TodoCacheGetRec(id)
	TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todo-cache record for ID: %d doesn't exist.", id)
		return false
	}

	pRec.DueDate = newDueDate
	logger.Log(logger.CORESERVER, logger.ERROR, "ShortDesc of todo-cache record updated. ID.DueDate :: %d.%v.", id, newDueDate.Format("02-Jan-2006::15:04:05"))

	return true
}


/* ****************************************************************************
Receiver    :

Arguments   :
1> newEDA time.Time: new EDA.
2> id uint32: Todo-cache record ID whose EDA needs to be updated.

Return value:
1> bool: true if EDA is updated successfully. false otherwise.

Description :
- Function updates EDA of the todo-cache record. Todo-cache record is fetched from the given id.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheUpdateEDA(newEDA time.Time, id uint32) bool {
	isOK, pRec := TodoCacheGetRec(id)
	TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todo-cache record for ID: %d doesn't exist.", id)
		return false
	}

	pRec.EDA = newEDA
	logger.Log(logger.CORESERVER, logger.ERROR, "ShortDesc of todo-cache record updated. ID.DueDate :: %d.%v.", id, newEDA.Format("02-Jan-2006::15:04:05"))

	return true
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> status uint8: New todo-cache record status.
2> todoID uint32: Todo-cache record ID whose status needs to be updated.

Return value:
1> bool: true if status for todo-cache rec updated successfully. false otherwise.

Description :
- Function updates status of todo-cache record.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheUpdateStatus(status uint8, todoID uint32) bool {
	isOK, pRec := TodoCacheGetRec(todoID)
	defer TodoCacheReleaseRec(pRec)
	if !isOK {
		logger.Log(logger.CORESERVER, logger.ERROR, "todoCache record for ID: %d doesn't exist.", todoID)
		return false
	}

	pRec.Status = status
	logger.Log(logger.CORESERVER, logger.ERROR, "todoCache record status updated, ID.status: %d.%d", todoID, status)

	return true
}


/* ****************************************************************************
Receiver    : na

Arguments   :
1> todoID uint32: Key of todo-cache record to be deleted.

Return value:
1> bool: returns true in any case, i.e., if record for the given key is found and deleted or even if a record for the key isn't found.

Description :
- Function removes todo-cache record with key todoID.

Additional note:
- Function takes and releases all the necessary locks.
- Caller need not take any lock prior to calling this function.
- Locking sequence, however, should be maintained.
**************************************************************************** */
func TodoCacheDeleteRec(todoID uint32) bool {
	TodoCacheWRLock()
	defer TodoCacheWRUnlock()

	if _, isOK := todoCache[todoID]; !isOK {
		logger.Log(logger.CORESERVER, logger.WARNING, "todoCache record for ID: %d doesn't exist.", todoID)
		return true
	}

	delete(todoCache, todoID)
	logger.Log(logger.CORESERVER, logger.DEBUG, "todoCache record for ID: %d deleted.", todoID)

	return true
}
