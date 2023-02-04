package main

import (
	"fmt"
)

func main() {
	var set = make(map[string]int)
	set["11"] = 1
	set["44"] = 4
	set["22"] = 2
	set["33"] = 3
	set["55"] = 5
	set["66"] = 6
	fmt.Println(set)
	delete(set, "11")
	fmt.Println(set["Bob"])
	for _, v := range set {
		fmt.Println(v)
	}
}
