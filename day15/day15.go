package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type vector2d struct {
	x int
	y int
}

type move struct {
	direction int
	position  vector2d
}

type computer struct {
	program  []int
	position int
	offset   int
	channel  chan int
}

func (c *computer) parseInst() (int, []int) {
	inst := c.read(c.position)
	opcode := inst % 100
	modes := make([]int, 4)
	inst /= 100
	for i := 0; inst > 0; i++ {
		modes[i] = inst % 10
		inst /= 10
	}

	return opcode, modes
}

func (c *computer) parseParams(modes []int, num int) []int {
	params := make([]int, num)
	for i := 0; i < num; i++ {
		if i >= len(modes) || modes[i] == 0 {
			params[i] = c.read(c.read(c.position + i + 1))
		} else if modes[i] == 1 {
			params[i] = c.read(c.position + i + 1)
		} else {
			params[i] = c.read(c.read(c.position+i+1) + c.offset)
		}
	}

	return params
}

func (c *computer) read(pos int) int {
	if pos >= len(c.program) {
		return 0
	}

	return c.program[pos]
}

func (c *computer) write(pos, val, mode int) {
	if mode == 2 {
		pos += c.offset
	}

	if pos >= len(c.program) {
		prog := make([]int, pos*2)
		copy(prog, c.program)
		c.program = prog
	}

	c.program[pos] = val
}

func (c *computer) run() {
	for c.read(c.position) != 99 {
		inst, modes := c.parseInst()
		switch inst {
		case 1:
			params := c.parseParams(modes, 2)
			c.write(c.read(c.position+3), params[0]+params[1], modes[2])
			c.position += 4
		case 2:
			params := c.parseParams(modes, 2)
			c.write(c.read(c.position+3), params[0]*params[1], modes[2])
			c.position += 4
		case 3:
			c.write(c.read(c.position+1), <-c.channel, modes[0])
			c.position += 2
		case 4:
			params := c.parseParams(modes, 1)
			c.channel <- params[0]
			c.position += 2
		case 5:
			params := c.parseParams(modes, 2)
			if params[0] != 0 {
				c.position = params[1]
			} else {
				c.position += 3
			}
		case 6:
			params := c.parseParams(modes, 2)
			if params[0] == 0 {
				c.position = params[1]
			} else {
				c.position += 3
			}
		case 7:
			params := c.parseParams(modes, 2)
			if params[0] < params[1] {
				c.write(c.read(c.position+3), 1, modes[2])
			} else {
				c.write(c.read(c.position+3), 0, modes[2])
			}
			c.position += 4
		case 8:
			params := c.parseParams(modes, 2)
			if params[0] == params[1] {
				c.write(c.read(c.position+3), 1, modes[2])
			} else {
				c.write(c.read(c.position+3), 0, modes[2])
			}
			c.position += 4
		case 9:
			params := c.parseParams(modes, 1)
			c.offset += params[0]
			c.position += 2
		}
	}
}

func createComputer() computer {
	input, _ := ioutil.ReadFile("input.txt")
	ints := make([]int, 0)
	for _, v := range strings.Split(string(input), ",") {
		i, _ := strconv.Atoi(v)
		ints = append(ints, i)
	}

	c := computer{program: ints, channel: make(chan int)}
	go c.run()

	return c
}

func getAdjacent(v vector2d) []vector2d {
	return []vector2d{{v.x, v.y - 1}, {v.x, v.y + 1}, {v.x - 1, v.y}, {v.x + 1, v.y}}
}

func makeMap() (int, map[vector2d]int, vector2d) {
	c := createComputer()

	min, visited, pos, path, oxy := math.MaxFloat64, make(map[vector2d]int), vector2d{}, []move{move{}}, vector2d{}
	reverse := []int{2, 1, 4, 3}
	for {
		directions := getAdjacent(pos)
		moved := false
		for i, d := range directions {
			_, ok := visited[d]
			if ok {
				continue
			}

			c.channel <- i + 1
			result := <-c.channel
			visited[d] = result
			if result == 0 {
				continue
			}

			if result == 2 {
				oxy = d
				min = math.Min(min, float64(len(path)))
			}

			pos = d
			path = append(path, move{direction: i, position: pos})
			moved = true
			break
		}

		if moved {
			continue
		} else if (pos == vector2d{}) {
			break
		}

		c.channel <- reverse[path[len(path)-1].direction]
		<-c.channel
		pos = path[len(path)-2].position
		path = path[:len(path)-1]
	}

	return int(min), visited, oxy
}

func part1() {
	steps, _, _ := makeMap()
	fmt.Println(steps)
}

func part2() {
	_, grid, oxy := makeMap()

	queue, visited := []vector2d{oxy}, make(map[vector2d]bool)
	minutes := -1
	for len(queue) > 0 {
		minutes++
		newqueue := make([]vector2d, 0)
		for _, d := range queue {
			visited[d] = true

			for _, a := range getAdjacent(d) {
				if !visited[a] && grid[a] == 1 {
					newqueue = append(newqueue, a)
				}
			}
		}

		queue = newqueue
	}

	fmt.Println(minutes)
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
