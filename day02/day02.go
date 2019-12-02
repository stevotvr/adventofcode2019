package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getInts() []int {
	input, _ := ioutil.ReadFile("input.txt")
	ints := make([]int, 0)
	for _, v := range strings.Split(string(input), ",") {
		i, _ := strconv.Atoi(v)
		ints = append(ints, i)
	}

	return ints
}

func part1() {
	ints := getInts()
	ints[1], ints[2] = 12, 2

	pos := 0
	for ints[pos] != 99 {
		switch ints[pos] {
		case 1:
			ints[ints[pos+3]] = ints[ints[pos+1]] + ints[ints[pos+2]]
		case 2:
			ints[ints[pos+3]] = ints[ints[pos+1]] * ints[ints[pos+2]]
		}

		pos += 4
	}

	fmt.Println(ints[0])
}

func part2() {
	ints := getInts()
	state := make([]int, len(ints))

	for n := 0; n < 100; n++ {
		for v := 0; v < 100; v++ {
			copy(state, ints)
			state[1], state[2] = n, v
			pos := 0
			for state[pos] != 99 {
				switch state[pos] {
				case 1:
					state[state[pos+3]] = state[state[pos+1]] + state[state[pos+2]]
				case 2:
					state[state[pos+3]] = state[state[pos+1]] * state[state[pos+2]]
				}

				pos += 4
			}

			if state[0] == 19690720 {
				fmt.Println(100*n + v)
				return
			}
		}
	}
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
