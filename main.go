package main

import (
	"fmt"
	"strings"
)

func main() {

	test := []int32{10, 20, 20, 10, 10, 30, 50, 10, 20}

	p := sockMerchant(9, test)
	fmt.Println(p)

	s := strings.Split("abc", "")
	fmt.Println(s)

	t := countingValleys(8, "UDDDUDUU")
	fmt.Println(t)

}

func countingValleys(n int32, s string) int32 {

	var lvl int32
	var v int32
	sr := strings.Split(s, "")
	for _, num := range sr {

		if num == "U" {
			lvl = lvl + 1
		}
		if num == "D" {
			lvl = lvl - 1
		}
		if lvl == 0 && num == "U" {
			v = v + 1
		}

	}

	return v
}

func sockMerchant(n int32, ar []int32) int32 {

	var pairs int32
	m := make(map[int32]int32)

	for _, num := range ar {

		if _, ok := m[num]; !ok {
			m[num] = num
		} else {
			pairs = pairs + 1
			delete(m, num)
		}

	}

	return pairs
}
