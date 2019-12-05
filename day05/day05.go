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

func parseInst(i int) (int, []bool) {
	str := strconv.Itoa(i)
	inst, _ := strconv.Atoi(str[len(str)-1:])
	modes := make([]bool, 0)
	for i := len(str) - 3; i >= 0; i-- {
		modes = append(modes, str[i] == '1')
	}

	return inst, modes
}

func parseParams(ints []int, pos int, modes []bool, num int) []int {
	params := make([]int, num)
	for i := 0; i < num; i++ {
		if i < len(modes) && modes[i] {
			params[i] = ints[pos+i+1]
		} else {
			params[i] = ints[ints[pos+i+1]]
		}
	}

	return params
}

func part1() {
	ints := getInts()
	input := 1

	pos := 0
	for ints[pos] != 99 {
		inst, modes := parseInst(ints[pos])
		switch inst {
		case 1:
			params := parseParams(ints, pos, modes, 2)
			ints[ints[pos+3]] = params[0] + params[1]
			pos += 4
		case 2:
			params := parseParams(ints, pos, modes, 2)
			ints[ints[pos+3]] = params[0] * params[1]
			pos += 4
		case 3:
			ints[ints[pos+1]] = input
			pos += 2
		case 4:
			params := parseParams(ints, pos, modes, 1)
			fmt.Println(params[0])
			pos += 2
		}
	}
}

func part2() {
	ints := getInts()
	input := 5

	pos := 0
	for ints[pos] != 99 {
		inst, modes := parseInst(ints[pos])
		switch inst {
		case 1:
			params := parseParams(ints, pos, modes, 2)
			ints[ints[pos+3]] = params[0] + params[1]
			pos += 4
		case 2:
			params := parseParams(ints, pos, modes, 2)
			ints[ints[pos+3]] = params[0] * params[1]
			pos += 4
		case 3:
			ints[ints[pos+1]] = input
			pos += 2
		case 4:
			params := parseParams(ints, pos, modes, 1)
			fmt.Println(params[0])
			pos += 2
		case 5:
			params := parseParams(ints, pos, modes, 2)
			if params[0] != 0 {
				pos = params[1]
			} else {
				pos += 3
			}
		case 6:
			params := parseParams(ints, pos, modes, 2)
			if params[0] == 0 {
				pos = params[1]
			} else {
				pos += 3
			}
		case 7:
			params := parseParams(ints, pos, modes, 2)
			if params[0] < params[1] {
				ints[ints[pos+3]] = 1
			} else {
				ints[ints[pos+3]] = 0
			}
			pos += 4
		case 8:
			params := parseParams(ints, pos, modes, 2)
			if params[0] == params[1] {
				ints[ints[pos+3]] = 1
			} else {
				ints[ints[pos+3]] = 0
			}
			pos += 4
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
