package algorithms

import (
	"strconv"
	"time"

	"github.com/n0l3r/limitron/store"
)

type FixedWindow struct {
	WindowSize int64 // in seconds
	Limit      int64
	Store      store.Store
}

func (fw *FixedWindow) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	window := now - (now % fw.WindowSize)
	windowKey := key + ":fw:" + strconv.FormatInt(window, 10)
	count, err := fw.Store.Get(windowKey)
	if err != nil {
		return false, err
	}
	if count < fw.Limit {
		err = fw.Store.Set(windowKey, count+1, fw.WindowSize)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
