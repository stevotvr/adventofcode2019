package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var incache []int

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

	close(c.channel)
}

func createComputer() computer {
	var ints []int

	if incache == nil {
		input, _ := ioutil.ReadFile("input.txt")
		ints = make([]int, 0)
		for _, v := range strings.Split(string(input), ",") {
			i, _ := strconv.Atoi(v)
			ints = append(ints, i)
		}

		incache = make([]int, len(ints))
		copy(incache, ints)
	} else {
		ints = make([]int, len(incache))
		copy(ints, incache)
	}

	c := computer{program: ints, channel: make(chan int)}
	go c.run()

	return c
}

func part1() {
	area := 0
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			c := createComputer()
			c.channel <- x
			c.channel <- y
			area += <-c.channel
		}
	}

	fmt.Println(area)
}

func part2() {
	for x := 0; ; x++ {
		for y := 0; ; y++ {
			c := createComputer()
			c.channel <- x + 99
			c.channel <- y
			if <-c.channel == 1 {
				c = createComputer()
				c.channel <- x
				c.channel <- y + 99
				if <-c.channel == 1 {
					fmt.Println(x*10000 + y)
					return
				}

				break
			}
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
