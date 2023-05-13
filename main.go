package main

import (
	"container/ring"
	"fmt"
	"math/rand"
	"time"
)

const ringCount = 50
const timeout = 100

func main() {
	rand.Seed(time.Now().UnixNano())
	init := func(done <-chan int, r *ring.Ring) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			r.Do(func(p any) {
				time.Sleep(timeout * time.Millisecond)
				select {
				case <-done:
					return
				case output <- p.(int):
				}
			})
		}()
		return output
	}

	plusNumber := func(done <-chan int, input <-chan int) <-chan int {
		plusedStream := make(chan int)
		go func() {
			defer close(plusedStream)
			for {
				select {
				case <-done:
					return
				case i, isChannelOpen := <-input:
					if !isChannelOpen {
						return
					}
					if i < 0 {
						continue
					}
					select {
					case plusedStream <- i:
					case <-done:
						return
					}
				}
			}

		}()
		return plusedStream
	}

	multiples3 := func(done <-chan int, input <-chan int) <-chan int {
		multiples3edStream := make(chan int)
		go func() {
			defer close(multiples3edStream)
			for {
				select {
				case <-done:
					return
				case i, isChannelOpen := <-input:
					if !isChannelOpen {
						return
					}
					if i == 0 || i%3 != 0 {
						continue
					}
					select {
					case multiples3edStream <- i:
					case <-done:
						return
					}
				}
			}
		}()
		return multiples3edStream
	}

	r := ring.New(ringCount)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		random := rand.Intn(20)
		if random > 10 {
			r.Value = i
		} else {
			r.Value = -i
		}
		r = r.Next()
	}

	done := make(chan int)
	defer close(done)
	intStream := init(done, r)

	pipeline := multiples3(done, plusNumber(done, intStream))
	for v := range pipeline {
		fmt.Println(v)
	}
}
