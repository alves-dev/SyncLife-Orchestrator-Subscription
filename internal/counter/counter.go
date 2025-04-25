package counter

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type DailyCounter struct {
	count    uint64
	lastDate time.Time
	mutex    sync.Mutex
}

func NewDailyCounter() *DailyCounter {
	return &DailyCounter{
		lastDate: time.Now().Truncate(24 * time.Hour),
	}
}

func (c *DailyCounter) Increment(eventType string) uint64 {
	now := time.Now().Truncate(24 * time.Hour)

	c.mutex.Lock()
	if now.After(c.lastDate) {
		atomic.StoreUint64(&c.count, 0)
		c.lastDate = now
	}
	c.mutex.Unlock()

	newCount := atomic.AddUint64(&c.count, 1)
	log.Printf("[DailyCounter] Type: %s | Daily Total: %d", eventType, newCount)
	return newCount
}
