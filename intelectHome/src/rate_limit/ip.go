package rate_limit

import (
	"fmt"
	"sync"
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
	maxRequestIpInSecond int64
	startQuantityTokens  int64
	maxToken             int64
	mtx                  sync.Mutex
}

type ipTokensBucket struct {
	tokens             int64
	maxToken           int64
	lastAction         int64
	isBlocked          bool
	maxRequestInMinute int64
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
		fmt.Println("attack")
		return false
	}
	if iRl.isStoped {
		fmt.Println("stop")
		return true
	}
	_, ok := iRl.storage[ip]
	if !ok {
		iRl.createNewIp(ip)
	}
	v := iRl.storage[ip]
	times := time.Now().UnixMicro() - v.lastAction
	tokensToLoad := int64(float64(times) * float64(iRl.maxRequestIpInSecond) / 1000000.0)
	if v.tokens <= iRl.maxToken {
		if tokensToLoad+v.tokens > iRl.maxToken {
			v.tokens = iRl.maxToken
		} else {
			v.tokens += tokensToLoad
		}
	}

	tokensAfterSub := v.tokens - 1
	if tokensAfterSub >= 0 {
		v.tokens -= 1
		v.lastAction = time.Now().UnixMicro()
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
		v.maxRequestInMinute = iRl.maxRequestIpInSecond
	}
	return true
}

func (iRl *IpRateLimiter) createNewIp(ip string) {
	v := &ipTokensBucket{}
	v.maxToken = iRl.maxRequestIpInSecond
	v.tokens = iRl.startQuantityTokens
	v.maxRequestInMinute = iRl.maxRequestIpInSecond
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
	ipBucket.isBlocked = true
	return true
}

func (iRl *IpRateLimiter) startCleningStorage() {
	ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			iRl.mtx.Lock()
			for k, v := range iRl.storage {
				if time.Now().Unix()-v.lastAction > fiveMinutesInSecond {
					fmt.Println(time.Now().Unix()-v.lastAction > fiveMinutesInSecond)
					delete(iRl.storage, k)
				}
			}
			iRl.mtx.Unlock()
		default:
			continue
		}
	}
}
