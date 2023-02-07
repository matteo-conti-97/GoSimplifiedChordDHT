package main

import (
	"fmt"
)

func main() {
	var x string
	var y = "0"
	fmt.Println("Start")
	fmt.Scanln(&x, &y)
	fmt.Println(x)
	fmt.Println(y)
}
