package rate_limit

import (
	"context"
	"fmt"
	"log"
	"sync"
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
	mtx         sync.Mutex
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
			if rl.tokenBucket.Load() < rl.maxTokens {
				rl.tokenBucket.Add(1)
			}
		case <-rl.ctx.Done():
			log.Printf("Tokens refil stop, count tokens: %d\n", rl.tokenBucket.Load())
			return
		}
	}
}

func (rl *GlobalRateLimiter) StopRefillToken(isAttacked bool) {
	rl.attacked.Store(isAttacked)
	rl.stopped.Store(true)
	rl.mtx.Lock()
	rl.ctxCancel()
	rl.mtx.Unlock()
}

func (rl *GlobalRateLimiter) Allow() bool {
	if rl.attacked.Load() {
		fmt.Println(1)
		return false
	}
	if rl.stopped.Load() {
		fmt.Println(2)
		return true
	}
	if tokens := rl.tokenBucket.Load(); tokens > 0 {
		fmt.Println(3)
		rl.tokenBucket.Add(-1)
		return true
	}
	fmt.Println(4)
	return false
}

func (rl *GlobalRateLimiter) Restart() {
	rl.mtx.Lock()
	rl.ctxCancel()
	rl.ctx, rl.ctxCancel = context.WithCancel(context.Background())
	rl.attacked.Store(false)
	rl.stopped.Store(false)
	rl.mtx.Unlock()
	go rl.refillBucket()
}

func (rl *GlobalRateLimiter) GetAttackedStatus() bool {
	return rl.attacked.Load()
}
