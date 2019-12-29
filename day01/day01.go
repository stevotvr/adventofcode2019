package main

import (
	"fmt"
	"os"
)

func part1() {
	input, _ := os.Open("input.txt")
	output, mass := 0, 0
	for {
		_, err := fmt.Fscanln(input, &mass)
		if err != nil {
			break
		}

		output += mass/3 - 2
	}

	fmt.Println(output)
}

func part2() {
	input, _ := os.Open("input.txt")
	output, mass := 0, 0
	for {
		_, err := fmt.Fscanln(input, &mass)
		if err != nil {
			break
		}

		for {
			mass = mass/3 - 2
			if mass <= 0 {
				break
			}

			output += mass
		}
	}

	fmt.Println(output)
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
