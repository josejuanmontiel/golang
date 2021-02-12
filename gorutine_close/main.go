package main

import (
	"fmt"
	"strconv"
)

var end = make(chan bool, 1)

func Close(channel chan bool) {
	select {
	case <-channel:
		fmt.Println("Channel is closed! SURE?")
	default:
		fmt.Println("Channel is not closed!")
		close(channel)
		fmt.Println("NOW we closed it!")
	}
}

func Funcion(iter int, channel chan bool) {
	go func() {
		for i := 0; i < iter; i++ {
			<-channel
			fmt.Println("Iteration: " + strconv.Itoa(i))
		}
		end <- true
	}()

	defer Close(channel)
	// Close(channel)
}

func main() {
	iter := 4
	closed := make(chan bool, iter)
	for i := 0; i < iter; i++ {
		closed <- true
	}

	Funcion(iter, closed)

	<-end

	fmt.Println("Break!")
}
