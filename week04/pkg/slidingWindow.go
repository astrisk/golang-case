package pkg
/*
计数器实现：计数器用map来实现bucket存储，通过key来定义不同的bucket
窗口实现：用时间戳实现窗口，默认窗口大小为10秒
参考实现：https://github.com/afex/hystrix-go/blob/master/hystrix/rolling/rolling.go
 */

import (
	"sync"
	"time"
)

type Count struct {
	Buckets map[int64]*countBucket
	Mutex   *sync.RWMutex
}

type countBucket struct {
	Value float64
}

func NewCount() *Count {
	return &Count{
		Buckets: make(map[int64]*countBucket),
		Mutex:   &sync.RWMutex{},
	}
}

// 获取当前的bucket
func (c *Count) getCurrentBucket() *countBucket {
	now := time.Now().Unix()
	var bucket *countBucket
	var ok bool

	if bucket, ok = c.Buckets[now]; !ok {
		bucket = &countBucket{}
		c.Buckets[now] = bucket
	}
	return bucket
}

// 删除过期数据
func (c *Count) removeExpiredBuckets() {
	now := time.Now().Unix() - 10
	for timestamp := range c.Buckets {
		if timestamp <= now {
			delete(c.Buckets, timestamp)
		}
	}
}

// 增加计数
func (c *Count) Add(i float64) {
	if i == 0 {
		return
	}
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	b := c.getCurrentBucket()
	b.Value += i
	c.removeExpiredBuckets()
}

// 获取汇总数据
func (c *Count) Sum(now time.Time) float64 {
	sum := float64(0)

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	for timestamp, bucket := range c.Buckets {
		if timestamp >= now.Unix()-10 {
			sum += bucket.Value
		}
	}
	return sum
}
