package dataCache

import (
	"sagacity.com/logger"
	"sagacity.com/models/dataCacheModels"
)

//dataCacheModels.go: type UserCacheMap map[string]*UserCacheRec
func Init() bool {
	todoCache = make(dataCacheModels.TodoCacheMap)
	if todoCache == nil {
		logger.Log(logger.CORESERVER, logger.ERROR, "Couldn't create todo-cache.")
		return false
	}

	return true
}
