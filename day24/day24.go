package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type vector2d struct {
	x int
	y int
}

func loadInput() int {
	input, _ := ioutil.ReadFile("input.txt")

	output, i := 0, 0
	for _, c := range input {
		if c == '#' {
			output |= 1 << i
			i++
		} else if c == '.' {
			i++
		}
	}

	return output
}

func isBug(grid, x, y int) bool {
	if x < 0 || x >= 5 || y < 0 || y >= 5 {
		return false
	}

	bit := 1 << (y*5 + x)
	return grid&bit == bit
}

func countNeighbors(grid, x, y int) int {
	directions := []vector2d{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	count := 0
	for _, d := range directions {
		if isBug(grid, x+d.x, y+d.y) {
			count++
		}
	}

	return count
}

func countNeighborsRecursive(levels map[int]int, level, x, y int) int {
	directions := []vector2d{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	count := 0
	for _, d := range directions {
		switch {
		case x+d.x == 2 && y+d.y == 2:
			if y == 2 {
				for y2 := 0; y2 < 5; y2++ {
					if (x == 1 && isBug(levels[level+1], 0, y2)) || (x == 3 && isBug(levels[level+1], 4, y2)) {
						count++
					}
				}
			} else {
				for x2 := 0; x2 < 5; x2++ {
					if (y == 1 && isBug(levels[level+1], x2, 0)) || (y == 3 && isBug(levels[level+1], x2, 4)) {
						count++
					}
				}
			}
		case x+d.x < 0 && isBug(levels[level-1], 1, 2), y+d.y < 0 && isBug(levels[level-1], 2, 1), x+d.x >= 5 && isBug(levels[level-1], 3, 2), y+d.y >= 5 && isBug(levels[level-1], 2, 3):
			count++
		case isBug(levels[level], x+d.x, y+d.y):
			count++
		}
	}

	return count
}

func part1() {
	grid := loadInput()

	visited, gridCopy := make(map[int]bool), 0
	for !visited[grid] {
		visited[grid] = true

		gridCopy = 0
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				count := countNeighbors(grid, x, y)
				bit := 1 << (y*5 + x)
				if grid&bit == bit {
					if count == 1 {
						gridCopy |= bit
					}
				} else if count == 1 || count == 2 {
					gridCopy |= bit
				}

			}
		}

		grid = gridCopy
	}

	fmt.Println(grid)
}

func part2() {
	levels := map[int]int{-1: 0, 0: loadInput(), 1: 0}
	var levelsCopy map[int]int
	var min, max int
	for i := 0; i < 200; i++ {
		levelsCopy = make(map[int]int)
		for level, grid := range levels {
			levelsCopy[level] = 0
			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					if x == 2 && y == 2 {
						continue
					}

					count := countNeighborsRecursive(levels, level, x, y)
					bit := 1 << (y*5 + x)
					if grid&bit == bit {
						if count == 1 {
							levelsCopy[level] |= bit
						}
					} else if count == 1 || count == 2 {
						levelsCopy[level] |= bit
					}
				}
			}
		}

		if levelsCopy[min-1] > 0 {
			min--
			levelsCopy[min-1] = 0
		}

		if levelsCopy[max+1] > 0 {
			max++
			levelsCopy[max+1] = 0
		}

		levels = levelsCopy
	}

	total := 0
	for _, grid := range levels {
		for i := 0; i < 25; i++ {
			if grid&(1<<i) == 1<<i {
				total++
			}
		}
	}

	fmt.Println(total)
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
