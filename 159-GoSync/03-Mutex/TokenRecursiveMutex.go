package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// TokenRecursiveMutex Token方式的递归锁
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// Lock 请求锁，需要传入token
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { // 如果传入的token和持有锁的token一致，说明是递归调用
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// Unlock 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}

	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有退回到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0)
	m.Mutex.Unlock()
}

func main() {

}