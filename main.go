package main

import (
	"fmt"
	"sync"
	"time"
)

type IP string

type Bucket struct {
	Tokens int
	sync.RWMutex
}

type RateLimiter struct {
	Buckets     map[IP]*Bucket
	MaxCapacity int
	RefillTime  int  // seconds
	RefillRate  int  // tokens added
	InfoLog     bool // for debug and dev environments
}

func NewRateLimiter(bucketCapacity, refillTime, refillRate int, infoLog bool) *RateLimiter {
	if refillRate > bucketCapacity {
		fmt.Printf("RATELIMITER WARNING: Refill rate is bigger than buckets max capacity\n")
	}
	rl := &RateLimiter{
		Buckets:     make(map[IP]*Bucket),
		MaxCapacity: bucketCapacity,
		RefillTime:  refillTime,
		RefillRate:  refillRate,
		InfoLog:     infoLog,
	}
	go rl.Refill()
	return rl
}

func (rl *RateLimiter) Refill() {
	for {
		time.Sleep(time.Duration(rl.RefillTime) * time.Second)
		for key, val := range rl.Buckets {
			val.Lock()
			if rl.InfoLog {
				fmt.Printf("RATELIMITER INFO %v: Bucket %s capacity %d added %d\n", time.Now().Format("15:04:05"), key, val.Tokens, rl.MaxCapacity-val.Tokens)
			}
			val.Tokens += rl.RefillRate
			if val.Tokens > rl.MaxCapacity {
				val.Tokens = rl.MaxCapacity
			}
			val.Unlock()
		}
	}
}

func (rl *RateLimiter) CheckRequest(ip IP) bool {
	bucket, ok := rl.Buckets[ip]
	if !ok {
		// IP is new need to be added to the ratelimiter
		bucket := &Bucket{Tokens: rl.MaxCapacity - 1} // Initialize the new bucket and subtract 1 because of new request
		rl.Buckets[ip] = bucket
		return true
	}
	if bucket.IsEmpty() {
		// no more tokens available, request needs to be dropped
		return false
	}
	bucket.Lock()

	bucket.Tokens-- // remove one token
	bucket.Unlock()
	return true
}

func (b *Bucket) IsEmpty() bool {
	b.RLock()
	isEmpty := b.Tokens == 0
	b.RUnlock()
	return isEmpty
}