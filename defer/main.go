package main

import "fmt"

func main() {
	defer fmt.Println("FINAL 1")
	defer fmt.Println("FINAL 2")
	defer fmt.Println("FINAL 3")

	fmt.Println("DOING SOMETHING")

}
