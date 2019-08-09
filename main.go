package main

import (
	"sync"
)

func main() {
	// 暴力的防止主进程退出
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
