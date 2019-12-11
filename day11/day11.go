package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	x int
	y int
}

type Robot struct {
	pos  Position
	dir  int
	grid map[Position]int
}

func (r *Robot) TurnLeft() {
	r.dir = (r.dir + 3) % 4
}

func (r *Robot) TurnRight() {
	r.dir = (r.dir + 1) % 4
}

func (r *Robot) Move() {
	switch r.dir {
	case 0:
		r.pos.y--
	case 1:
		r.pos.x++
	case 2:
		r.pos.y++
	case 3:
		r.pos.x--
	}
}

func (r *Robot) Paint(color int) {
	r.grid[r.pos] = color
}

func (r *Robot) GetColor() int {
	return r.grid[r.pos]
}

func (r *Robot) Run(c Computer) {
	for {
		_, ok := <-c.waiting
		if !ok {
			break
		}

		c.channel <- r.GetColor()
		color := <-c.channel
		turn := <-c.channel

		r.Paint(color)
		if turn == 0 {
			r.TurnLeft()
		} else {
			r.TurnRight()
		}
		r.Move()
	}
}

type Computer struct {
	program  []int
	position int
	offset   int
	channel  chan int
	waiting  chan bool
}

func (c *Computer) parseInst() (int, []int) {
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

func (c *Computer) parseParams(modes []int, num int) []int {
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

func (c *Computer) read(pos int) int {
	if pos >= len(c.program) {
		return 0
	}

	return c.program[pos]
}

func (c *Computer) write(pos, val, mode int) {
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

func (c *Computer) run() {
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
			c.waiting <- true
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

	close(c.channel)
	close(c.waiting)
}

func createRobot() Robot {
	return Robot{grid: make(map[Position]int)}
}

func createComputer() Computer {
	input, _ := ioutil.ReadFile("input.txt")
	ints := make([]int, 0)
	for _, v := range strings.Split(string(input), ",") {
		i, _ := strconv.Atoi(v)
		ints = append(ints, i)
	}

	c := Computer{program: ints, channel: make(chan int), waiting: make(chan bool)}
	go c.run()

	return c
}

func part1() {
	r := createRobot()
	c := createComputer()

	r.Run(c)

	fmt.Println(len(r.grid))
}

func part2() {
	r := createRobot()
	r.grid[r.pos] = 1
	c := createComputer()

	r.Run(c)

	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for k := range r.grid {
		minX = math.Min(minX, float64(k.x))
		minY = math.Min(minY, float64(k.y))
		maxX = math.Max(maxX, float64(k.x))
		maxY = math.Max(maxY, float64(k.y))
	}

	for y := int(minY); y <= int(maxY); y++ {
		for x := int(minX); x <= int(maxX); x++ {
			c, ok := r.grid[Position{x, y}]
			if ok && c == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
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
