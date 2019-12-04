package main

import (
	"fmt"
	"os"
	"strconv"
)

func getRange() (int, int) {
	input, _ := os.Open("input.txt")
	min, max := 0, 0
	fmt.Fscanf(input, "%d-%d", &min, &max)

	return min, max
}

func getDigits(n int) []int {
	str := strconv.Itoa(n)
	digits := make([]int, len(str))
	for i, v := range str {
		d, _ := strconv.Atoi(string(v))
		digits[i] = d
	}

	return digits
}

func valid(n int, part2 bool) bool {
	digits := getDigits(n)
	hasRepeat := false
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] > digits[i+1] {
			return false
		}

		if digits[i] == digits[i+1] && (!part2 || ((i == 0 || digits[i] != digits[i-1]) && (i == len(digits)-2 || digits[i] != digits[i+2]))) {
			hasRepeat = true
		}
	}

	return hasRepeat
}

func part1() {
	min, max := getRange()
	count := 0
	for i := min; i <= max; i++ {
		if valid(i, false) {
			count++
		}
	}

	fmt.Println(count)
}

func part2() {
	min, max := getRange()
	count := 0
	for i := min; i <= max; i++ {
		if valid(i, true) {
			count++
		}
	}

	fmt.Println(count)
}

func main() {
	part := 0
	fmt.Sscan(os.Args[1], &part)
	switch part {
	case 1:
		part1()
	case 2:
		part2()
	}
}
