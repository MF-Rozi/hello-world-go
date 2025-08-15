package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)

	go print(5, "Hello from Goroutine")
	print(5, "Hello from Main")

	time.Sleep(1 * time.Second)
	fmt.Println("All Goroutines finished")
}
func print(till int, message string) {
	for i := 1; i <= till; i++ {
		fmt.Println("current Active Goroutines:", runtime.NumGoroutine())
		fmt.Println(message, i)
		// time.Sleep(250 * time.Millisecond)
	}
}
