package blog

import "sync"

var loggerMgr sync.Map

func GetLogger(key any) *Logger {
	if l, ok := loggerMgr.Load(key); ok {
		return l.(*Logger)
	} else {
		// Prevents return nil, consecutive call failure
		Warnf("GetLogger Warn : cannot found by key : %v !", key)
		return Default()
	}
}
