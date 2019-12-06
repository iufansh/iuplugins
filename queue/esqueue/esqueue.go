package esqueue

import (
	"encoding/json"
	"errors"
	"iufan/cplugins/queue"
	"runtime"
	"sync/atomic"
)

type esCache struct {
	value interface{}
	mark  bool
}

// lock free queue
type EsQueue struct {
	capaciity uint32
	capMod    uint32
	putPos    uint32
	getPos    uint32
	cache     []esCache
}

func NewEsQueue() queue.Queue {
	q := EsQueue{}
	return &q
}

func (q *EsQueue) Capaciity() uint32 {
	return q.capaciity
}

func (q *EsQueue) Quantity() uint32 {
	var putPos, getPos uint32
	var quantity uint32
	getPos = q.getPos
	putPos = q.putPos

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + putPos - getPos
	}

	return quantity
}
func (q *EsQueue) PutByName(name string, val interface{}) error {
	return nil
}

func (q *EsQueue) GetByName(name string) interface{} {
	return nil
}

// put queue functions
func (q *EsQueue) Put(val interface{}) error {
	var putPos, putPosNew, getPos, posCnt uint32
	var cache *esCache
	capMod := q.capMod
	for {
		getPos = q.getPos
		putPos = q.putPos

		if putPos >= getPos {
			posCnt = putPos - getPos
		} else {
			posCnt = capMod + putPos - getPos
		}

		if posCnt >= capMod {
			runtime.Gosched()
			return errors.New("put queue error posCnt >= capMod")
		}

		putPosNew = putPos + 1
		if atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	cache = &q.cache[putPosNew&capMod]

	for {
		if !cache.mark {
			cache.value = val
			cache.mark = true
			return nil
		} else {
			runtime.Gosched()
		}
	}
}

// get queue functions
func (q *EsQueue) Get() interface{} {
	var putPos, getPos, getPosNew, posCnt uint32
	var cache *esCache
	capMod := q.capMod
	for {
		putPos = q.putPos
		getPos = q.getPos

		if putPos >= getPos {
			posCnt = putPos - getPos
		} else {
			posCnt = capMod + putPos - getPos
		}

		if posCnt < 1 {
			runtime.Gosched()
			return nil
		}

		getPosNew = getPos + 1
		if atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	cache = &q.cache[getPosNew&capMod]

	for {
		if cache.mark {
			val := cache.value
			cache.mark = false
			return val
		} else {
			runtime.Gosched()
		}
	}
}

func (q *EsQueue) ClearAll() error {
	return nil
}

func (q *EsQueue) StartAndGC(config string) error {
	var cf map[string]uint32
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["capaciity"]; !ok {
		cf = make(map[string]uint32)
		cf["capaciity"] = 100
	}
	q.capaciity = minQuantity(cf["capaciity"])
	q.capMod = q.capaciity - 1
	q.cache = make([]esCache, q.capaciity)
	return nil
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func init() {
	queue.Register("list", NewEsQueue)
}
