package algorithms

import (
	"github.com/n0l3r/limitron/store"
)

type SlidingWindow struct {
	WindowSize int64 // in seconds
	Limit      int64
	Store      store.Store
}

func (sw *SlidingWindow) Allow(key string) (bool, error) {
	windowKey := key + ":sw"
	count, err := sw.Store.Get(windowKey)
	if err != nil {
		return false, err
	}
	if count < sw.Limit {
		err = sw.Store.Set(windowKey, count+1, sw.WindowSize)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
