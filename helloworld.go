package main

import (
	"fmt"
)

func RemoveIndex(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func main() {
	set := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	var ret []string
	var temp []string
	for _, v := range set {
		if "4" >= v {
			ret = append(ret, v)
		} else {
			temp = append(temp, v)
		}
	}

	fmt.Println(set)
	set = temp
	fmt.Println(set)
	fmt.Println(ret)
}
