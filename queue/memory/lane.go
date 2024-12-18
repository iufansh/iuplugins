package memory

import (
	"errors"
	"iuplugins/queue"

	"fmt"

	"github.com/oleiade/lane"
)

type LaneQueue struct {
	q *lane.Queue
}

func NewLaneQueue() queue.Queue {
	return &LaneQueue{q: lane.NewQueue()}
}

func (lq *LaneQueue) PutByName(name string, val interface{}) error {
	return nil
}

func (lq *LaneQueue) GetByName(name string) interface{} {
	return nil
}

// put queue functions
func (lq *LaneQueue) Put(val interface{}) error {
	if lq.q == nil {
		return errors.New("queue is not init")
	}
	lq.q.Enqueue(val)
	return nil
}

// get queue functions
func (lq *LaneQueue) Get() interface{} {
	if lq.q == nil {
		return nil
	}
	return lq.q.Dequeue()
}

func (lq *LaneQueue) ClearAll() error {
	return nil
}

func (lq *LaneQueue) StartAndGC(config string) error {
	return nil
}

func init() {
	fmt.Println("init memory queue")
	queue.Register("memory", NewLaneQueue)
}
