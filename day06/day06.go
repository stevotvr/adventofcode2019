package main

import (
	"fmt"
	"os"
	"strings"
)

func getMap() map[string]string {
	input, _ := os.Open("input.txt")
	line, orbits := "", make(map[string]string)
	for {
		_, err := fmt.Fscanln(input, &line)
		if err != nil {
			break
		}

		o := strings.SplitN(line, ")", 2)
		orbits[o[1]] = o[0]
	}

	return orbits
}

func part1() {
	orbits := getMap()
	total := 0
	for _, v := range orbits {
		for {
			total++
			p, ok := orbits[v]
			if !ok {
				break
			}
			v = p
		}
	}

	fmt.Println(total)
}

func part2() {
	orbits := getMap()
	path1, path2 := make([]string, 0), make([]string, 0)

	o, ok := "YOU", false
	for {
		o, ok = orbits[o]
		if !ok {
			break
		}
		path1 = append(path1, o)
	}

	o = "SAN"
	for {
		o, ok = orbits[o]
		if !ok {
			break
		}
		path2 = append(path2, o)
	}

	common := 0
	for i, j := len(path1)-1, len(path2)-1; i >= 0 && j >= 0; {
		if path1[i] == path2[j] {
			common++
		} else {
			break
		}
		i--
		j--
	}

	fmt.Println(len(path1) + len(path2) - common*2)
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
