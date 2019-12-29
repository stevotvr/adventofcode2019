package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func part1() {
	input, _ := ioutil.ReadFile("input.txt")
	w, h := 25, 6
	chunk := w * h

	min, out := math.MaxInt64, 0
	for i := 0; i < len(input); i += chunk {
		str := string(input[i : i+chunk])
		zeros := strings.Count(str, "0")
		if zeros < min {
			min = zeros
			out = strings.Count(str, "1") * strings.Count(str, "2")
		}
	}

	fmt.Println(out)
}

func part2() {
	input, _ := ioutil.ReadFile("input.txt")
	w, h := 25, 6
	chunk := w * h

	output := make([]rune, chunk)

	for i := 0; i < len(input); i += chunk {
		for j, px := range input[i : i+chunk] {
			if output[i%chunk+j] == 0 {
				if px == '0' {
					output[i%chunk+j] = ' '
				} else if px == '1' {
					output[i%chunk+j] = '#'
				}
			}
		}
	}

	for i := 0; i < len(output); i += w {
		fmt.Println(string(output[i : i+w]))
	}
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
