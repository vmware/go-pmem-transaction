package transaction

import (
	"log"
	"reflect"
	"unsafe"
)

const (
	LOGSIZE     int = 4 * 1024 * 1024
	LBUFFERSIZE     = 512 * 1024
	BUFFERSIZE      = 4 * 1024
)

// transaction interface
type (
	TX interface {
		Begin() error
		Log(...interface{}) error
		ReadLog(interface{}) (interface{}, error)
		Exec(...interface{}) (error, []reflect.Value)
		FakeLog(interface{})
		End() error
		abort() error
	}
)

func Init(logHeadPtr unsafe.Pointer, logType string) unsafe.Pointer {
	switch logType {
	case "undo":
		return initUndoTx(logHeadPtr)
	case "redo":
		return initRedoTx(logHeadPtr)
	default:
		log.Panic("initializing unsupported transaction! Try undo/redo")
	}
	return nil
}

func Release(t TX) {
	switch v := t.(type) {
	case *undoTx:
		releaseUndoTx(v)
	case *redoTx:
		releaseRedoTx(v)
	default:
		log.Panic("Releasing unsupported transaction!")
	}
}
