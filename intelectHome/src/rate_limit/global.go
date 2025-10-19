package rate_limit

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type GlobalRateLimiter struct {
	tokenBucket atomic.Int64
	maxTokens   int64
	duration    time.Duration
	ctx         context.Context
	ctxCancel   context.CancelFunc
	stopped     atomic.Bool
	attacked    atomic.Bool
}

func MakeGlobalRateLimiter(maxRequestInSecond int, StartQuantityTokens int64) *GlobalRateLimiter {
	if maxRequestInSecond < 0 {
		panic("request in seconds <= 0")
	}
	if StartQuantityTokens > 50 {
		log.Printf("ATTENTIONS! A large number start quantity tokens %d\n", StartQuantityTokens)
	}
	rl := &GlobalRateLimiter{}
	rl.tokenBucket.Store(StartQuantityTokens)
	rl.ctx, rl.ctxCancel = context.WithCancel(context.Background())
	rl.duration = time.Duration(1000/maxRequestInSecond) * time.Millisecond
	fmt.Println(rl.duration)
	rl.maxTokens = int64(maxRequestInSecond)
	rl.stopped.Store(false)
	rl.attacked.Store(false)
	go rl.refillBucket()
	return rl
}

func (rl *GlobalRateLimiter) refillBucket() {
	ticker := time.NewTicker(rl.duration)

	for {
		select {
		case <-ticker.C:
			tokens := rl.tokenBucket.Load()
			if tokens < rl.maxTokens {
				rl.tokenBucket.Add(1)
			}
		case <-rl.ctx.Done():
			log.Printf("Tokens refil stop, count tokens: %d\n", rl.tokenBucket.Load())
			rl.stopped.Store(true)
			return
		}
	}
}

func (rl *GlobalRateLimiter) StopRefillToken(isAttacked bool) {
	rl.attacked.Store(isAttacked)
	rl.ctxCancel()
}

func (rl *GlobalRateLimiter) Allow() bool {
	if rl.attacked.Load() {
		return false
	}
	if rl.stopped.Load() {
		return true
	}
	// fmt.Printf("Before: %d\n", rl.tokenBucket.Load())
	if tokens := rl.tokenBucket.Load(); tokens > 0 {
		rl.tokenBucket.Add(-1)
		// fmt.Printf("After: %d\n", rl.tokenBucket.Load())
		return true
	}
	return false
}
