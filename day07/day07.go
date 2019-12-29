package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

func runProg(ints []int, ch chan int) {
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
			ints[ints[pos+1]] = <-ch
			pos += 2
		case 4:
			params := parseParams(ints, pos, modes, 1)
			ch <- params[0]
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

// https://stackoverflow.com/a/30226442/361721
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func part1() {
	ints := getInts()
	state := make([]int, len(ints))
	ch := make(chan int)

	max := 0.0
	for _, phases := range permutations([]int{0, 1, 2, 3, 4}) {
		output := 0
		for _, p := range phases {
			copy(state, ints)
			go runProg(state, ch)
			ch <- p
			ch <- output
			output = <-ch
		}

		max = math.Max(max, float64(output))
	}

	fmt.Println(int(max))
}

func part2() {
	ints := getInts()
	chans := make([]chan int, 5)

	max := 0.0
	for _, phases := range permutations([]int{5, 6, 7, 8, 9}) {
		for i := range phases {
			state := make([]int, len(ints))
			copy(state, ints)
			chans[i] = make(chan int)
			go runProg(state, chans[i])
			chans[i] <- phases[i]
		}

		output, halt := 0, false
		for !halt {
			for i := range phases {
				select {
				case chans[i] <- output:
				default:
					halt = true
				}

				if halt {
					break
				}

				output = <-chans[i]
			}
		}

		max = math.Max(max, float64(output))
	}

	fmt.Println(int(max))
}

func main() {
	part := 0
	if len(os.Args) == 2 {
		fmt.Sscan(os.Args[1], &part)
	} else {
		fmt.Print("Enter 1 or 2 to select part: ")
		fmt.Scanf("%d\n", &part)
	}

	switch part {
	case 1:
		part1()
	case 2:
		part2()
	default:
		fmt.Println("Error: Invalid part.")
	}
}
