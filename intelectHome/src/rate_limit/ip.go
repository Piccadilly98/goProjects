package rate_limit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
ip rate limited - 100 requests in minutes
no gorutine
cleaning ip map once 5 minutes
*/

var (
	fiveMinutesInSecond = int64(5 * time.Minute.Seconds())
)

type IpRateLimiter struct {
	storage              map[string]*ipTokensBucket
	countIpAddres        int64
	isStoped             bool
	isAttacked           bool
	mtx                  sync.RWMutex
	maxRequestIpInSecond int64
	startQuantityTokens  int64
	maxToken             int64
}

type ipTokensBucket struct {
	tokens             atomic.Int64
	maxToken           atomic.Int64
	lastAction         atomic.Int64
	isBlocked          atomic.Bool
	maxRequestInMinute atomic.Int64
}

func MakeIpRateLimiter(maxRequestIpInSecond int64, startTokensQuantity int64) *IpRateLimiter {
	if maxRequestIpInSecond <= 0 {
		panic("max request in minute <= 0!")
	}
	if startTokensQuantity < 0 {
		panic("max start quantuty tokens <0")
	}
	ipRl := &IpRateLimiter{
		storage:              make(map[string]*ipTokensBucket),
		maxRequestIpInSecond: maxRequestIpInSecond,
		startQuantityTokens:  startTokensQuantity,
		maxToken:             maxRequestIpInSecond,
	}
	go ipRl.startCleningStorage()
	return ipRl
}

func (iRl *IpRateLimiter) Allow(ip string) bool {
	iRl.mtx.Lock()
	defer iRl.mtx.Unlock()
	if iRl.isAttacked {
		return false
	}
	if iRl.isStoped {
		return true
	}
	_, ok := iRl.storage[ip]
	if !ok {
		iRl.createNewIp(ip)
	}
	v := iRl.storage[ip]
	times := time.Now().Unix() - v.lastAction.Load()
	tokensToLoad := times * iRl.maxRequestIpInSecond
	if tokensToLoad+v.tokens.Load() > iRl.maxToken {
		v.tokens.Store(iRl.maxToken)
	} else {
		v.tokens.Store(tokensToLoad)
	}
	tokensAfterSub := v.tokens.Load() - 1
	if tokensAfterSub >= 0 {
		v.tokens.Add(-1)
		v.lastAction.Store(time.Now().Unix())
		return true
	}
	return false
}

func (iRl *IpRateLimiter) Stop(isAttacked bool) {
	iRl.mtx.Lock()
	iRl.isAttacked = isAttacked
	iRl.isStoped = true
	iRl.mtx.Unlock()
}

func (iRl *IpRateLimiter) LimitRequestsForEveryone(maxRequestIpInSecond int64) bool {
	if maxRequestIpInSecond <= 0 {
		panic("max request in minute <= 0!")
	}
	iRl.mtx.Lock()
	defer iRl.mtx.Unlock()
	if iRl.countIpAddres <= 0 {
		return false
	}
	iRl.maxRequestIpInSecond = maxRequestIpInSecond / iRl.countIpAddres

	for _, v := range iRl.storage {
		v.maxRequestInMinute.Store(iRl.maxRequestIpInSecond)
	}
	return true
}

func (iRl *IpRateLimiter) createNewIp(ip string) {
	v := &ipTokensBucket{}
	v.maxToken.Store(iRl.maxRequestIpInSecond)
	v.tokens.Store(iRl.maxToken)
	v.maxRequestInMinute.Store(iRl.maxRequestIpInSecond)
	iRl.countIpAddres++
	iRl.storage[ip] = v
}

func (iRl *IpRateLimiter) BlockedIp(ip string) bool {
	iRl.mtx.Lock()
	defer iRl.mtx.Unlock()
	ipBucket, ok := iRl.storage[ip]
	if !ok {
		return false
	}
	ipBucket.isBlocked.Store(true)
	return true
}

func (iRl *IpRateLimiter) startCleningStorage() {
	ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			iRl.mtx.Lock()
			for k, v := range iRl.storage {
				if time.Now().Unix()-v.lastAction.Load() > fiveMinutesInSecond {
					fmt.Println(time.Now().Unix()-v.lastAction.Load() > fiveMinutesInSecond)
					delete(iRl.storage, k)
				}
			}
			iRl.mtx.Unlock()
		default:
			continue
		}
	}
}
