package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func getInts() []int {
	input, _ := ioutil.ReadFile("input.txt")
	output := make([]int, len(input))
	for i, c := range input {
		n, _ := strconv.Atoi(string(c))
		output[i] = n
	}

	return output
}

func part1() {
	list, mask := getInts(), []int{0, 1, 0, -1}
	for i := 0; i < 100; i++ {
		newlist := make([]int, len(list))
		for j := 0; j < len(list); j++ {
			total := 0
			for k := j; k < len(list); k++ {
				total += list[k] * mask[((k+1)/(j+1))%4]
			}

			newlist[j] = int(math.Abs(float64(total % 10)))
		}

		list = newlist
	}

	for _, i := range list[:8] {
		fmt.Print(i)
	}

	fmt.Println()
}

func part2() {
	list := getInts()
	offset := list[0]*1000000 + list[1]*100000 + list[2]*10000 + list[3]*1000 + list[4]*100 + list[5]*10 + list[6]

	newlist := make([]int, len(list)*10000)
	for i := 0; i < len(newlist); i += len(list) {
		copy(newlist[i:], list)
	}

	list = newlist

	newlist = make([]int, len(list))
	for i := 0; i < 100; i++ {
		newlist[len(list)-1] = list[len(list)-1]
		for j := len(list) - 2; j >= 0; j-- {
			newlist[j] = (list[j] + newlist[j+1]) % 10
		}

		list = newlist
	}

	for _, i := range list[offset : offset+8] {
		fmt.Print(i)
	}

	fmt.Println()
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
