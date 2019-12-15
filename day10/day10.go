package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

type asteroid struct {
	x int
	y int
}

func getGrid() []asteroid {
	input, _ := os.Open("input.txt")
	grid := make([]asteroid, 0)
	line, y := "", 0
	for {
		_, err := fmt.Fscanln(input, &line)
		if err != nil {
			break
		}

		for x := range line {
			if line[x] == '#' {
				grid = append(grid, asteroid{x, y})
			}
		}

		y++
	}

	return grid
}

func los(grid []asteroid, a1, a2 asteroid) bool {
	dist := int((math.Sqrt(math.Pow(float64(a2.x-a1.x), 2) + math.Pow(float64(a2.y-a1.y), 2))) * 1000000000)
	for _, a3 := range grid {
		if a3 == a1 || a3 == a2 {
			continue
		}

		dist1 := math.Sqrt(math.Pow(float64(a3.x-a1.x), 2) + math.Pow(float64(a3.y-a1.y), 2))
		dist2 := math.Sqrt(math.Pow(float64(a3.x-a2.x), 2) + math.Pow(float64(a3.y-a2.y), 2))

		if dist == int((dist1+dist2)*1000000000) {
			return false
		}
	}

	return true
}

func getAngle(a1, a2 asteroid) float64 {
	rads := math.Atan2(float64(a2.y-a1.y), float64(a2.x-a1.x))
	rads = math.Mod(rads+2*math.Pi, 2*math.Pi)
	return math.Mod(rads+math.Pi/2, 2*math.Pi)
}

func part1() {
	grid := getGrid()
	max := 0.0
	for _, a1 := range grid {
		count := 0
		for _, a2 := range grid {
			if a2 != a1 && los(grid, a1, a2) {
				count++
			}
		}

		max = math.Max(max, float64(count))
	}

	fmt.Println(int(max))
}

func part2() {
	grid := getGrid()
	max, station := 0, grid[0]
	for _, a1 := range grid {
		count := 0
		for _, a2 := range grid {
			if a2 != a1 && los(grid, a1, a2) {
				count++
			}
		}

		if count > max {
			station = a1
			max = count
		}
	}

	num, dead := 0, asteroid{}
	for {
		visible, angles := make(map[float64]int), make([]float64, 0)
		for i, a := range grid {
			if a != dead && a != station && los(grid, station, a) {
				angle := getAngle(station, a)
				angles = append(angles, angle)
				visible[angle] = i
			}
		}

		sort.Float64s(angles)
		for _, angle := range angles {
			i := visible[angle]
			num++
			if num == 200 {
				fmt.Println(grid[i].x*100 + grid[i].y)
				return
			}

			grid[i] = dead
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
