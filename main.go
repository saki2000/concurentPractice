package main

import (
	"fmt"
	"time"
)

// func printSth(x string) {
// 	fmt.Println(x)
// }

// func sendData(c chan int) {
// 	for i := 1; i <= 10; i++ {
// 		c <- i
// 	}
// }

// func reciveData(c chan int) {
// 	for j := 1; j <= 10; j++ {
// 		data := <-c
// 		fmt.Println(data)
// 	}
// }

var requests chan operation = make(chan operation)

var done chan struct{} = make(chan struct{})

type operation struct {
	action    string
	parameter string
	response  chan string
}

func OrderCoffee(coffeeType string) {
	order := operation{action: "order",
		parameter: coffeeType}

	requests <- order
}

func Start() {
	go monitorRequests()
}

func Stop() {
	shutdown := operation{action: "shutdown", parameter: "", response: nil}
	requests <- shutdown
	<-done
}

func monitorRequests() {
	// our protected data
	lastCoffeeMade := "nothing started yet"

	for op := range requests {
		fmt.Println("Actor processing " + op.action + " " + op.parameter)

		switch op.action {
		case "order":
			requestedCoffee := op.parameter
			lastCoffeeMade = requestedCoffee
			makeCoffee(requestedCoffee)

		case "lastmade":
			op.response <- lastCoffeeMade

		case "shutdown":
			// Stop accepting new requests
			fmt.Println("Shutting down")
			close(requests)
		}
	}

	// Signal all requests completed
	fmt.Println("All requests processed")
	close(done)
}

func GetLastCoffeeMade() string {
	answer := make(chan string)

	lastmade := operation{action: "lastmade",
		parameter: "",
		response:  answer}
	requests <- lastmade
	return <-answer
}

func makeCoffee(coffee string) {
	fmt.Println("Brewing " + coffee)
	time.Sleep(1 * time.Second)
	fmt.Println(coffee + " now ready")
}

func main() {

	// go printSth("one")
	// go printSth("two")
	// go printSth("three")
	// go printSth("four")
	// go printSth("five")
	// go printSth("six")
	// go printSth("seven")
	// go printSth("eight")
	// go printSth("nine")
	// go printSth("ten")
	// time.Sleep(2 * time.Second)

	// channel := make(chan int)

	// go sendData(channel)

	// go reciveData(channel)

	// time.Sleep(time.Second)

	Start()
	defer Stop()

	go OrderCoffee("Mocha")
	go OrderCoffee("Choca")
	go fmt.Println("Last made answer: " + GetLastCoffeeMade())
	go OrderCoffee("Tea, milk, no sugar")
	go fmt.Println("Last made answer: " + GetLastCoffeeMade())
	go OrderCoffee("Banana Shake")
	go fmt.Println("Last made answer: " + GetLastCoffeeMade())

}
