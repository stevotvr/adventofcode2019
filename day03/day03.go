package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type coords struct {
	x     int
	y     int
	steps int
}

func (c coords) Equals(c2 coords) bool {
	return c.x == c2.x && c.y == c2.y
}

func (c coords) Dist() float64 {
	return math.Abs(float64(c.x)) + math.Abs(float64(c.y))
}

func getWires() [][]coords {
	input, _ := ioutil.ReadFile("input.txt")
	output := make([][]coords, 0)

	for _, line := range strings.Split(string(input), "\n") {
		wire := make([]coords, 0)
		x, y, dx, dy, steps := 0, 0, 0, 0, 0
		for _, move := range strings.Split(line, ",") {
			dir, dist := byte(0), 0
			fmt.Sscanf(move, "%c%d", &dir, &dist)
			switch dir {
			case 'U':
				dx, dy = 0, -1
			case 'R':
				dx, dy = 1, 0
			case 'D':
				dx, dy = 0, 1
			case 'L':
				dx, dy = -1, 0
			}

			for i := dist; i > 0; i-- {
				x, y = x+dx, y+dy
				steps++
				wire = append(wire, coords{x, y, steps})
			}
		}

		output = append(output, wire)
	}

	return output
}

func part1() {
	wires := getWires()
	min := math.MaxFloat64
	for _, w1 := range wires[0] {
		for _, w2 := range wires[1] {
			if w1.Equals(w2) {
				min = math.Min(min, w1.Dist())
			}
		}
	}

	fmt.Println(int(min))
}

func part2() {
	wires := getWires()
	min := math.MaxFloat64
	for _, w1 := range wires[0] {
		for _, w2 := range wires[1] {
			if w1.Equals(w2) {
				min = math.Min(min, float64(w1.steps+w2.steps))
			}
		}
	}

	fmt.Println(int(min))
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
