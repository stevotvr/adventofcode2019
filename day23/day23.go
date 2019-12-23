package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type vector2d struct {
	x int
	y int
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

	close(c.channel)
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

func createNetwork() (devices []computer, queue [][]vector2d) {
	devices = make([]computer, 50)
	queue = make([][]vector2d, len(devices))
	for i := 0; i < len(devices); i++ {
		devices[i] = createComputer()
		devices[i].channel <- i
		queue[i] = make([]vector2d, 0)
	}

	return
}

func part1() {
	devices, queue := createNetwork()

	for {
		for i, dev := range devices {
			var packet vector2d
			if len(queue[i]) > 0 {
				packet = queue[i][0]
			} else {
				packet = vector2d{-1, -1}
			}

			select {
			case dest := <-dev.channel:
				x := <-dev.channel
				y := <-dev.channel
				if dest == 255 {
					fmt.Println(y)
					return
				}

				queue[dest] = append(queue[dest], vector2d{x, y})
			case dev.channel <- packet.x:
				if packet.x != -1 {
					dev.channel <- packet.y
					queue[i] = queue[i][1:]
				}
			}
		}
	}
}

func part2() {
	devices, queue := createNetwork()

	nat, natsent := vector2d{}, vector2d{-1, -1}

	for {
		idle := true
		for i, dev := range devices {
			var packet vector2d
			if len(queue[i]) > 0 {
				packet = queue[i][0]
			} else {
				packet = vector2d{-1, -1}
			}

			select {
			case dest := <-dev.channel:
				x := <-dev.channel
				y := <-dev.channel
				if dest == 255 {
					nat = vector2d{x, y}
				} else {
					queue[dest] = append(queue[dest], vector2d{x, y})
				}

				idle = false
			case dev.channel <- packet.x:
				if packet.x != -1 {
					dev.channel <- packet.y
					queue[i] = queue[i][1:]
					idle = false
				}
			}
		}

		if !idle {
			continue
		}

		for _, q := range queue {
			if len(q) > 0 {
				idle = false
				break
			}
		}

		if idle {
			if nat == natsent {
				fmt.Println(nat.y)
				return
			}

			queue[0], natsent = append(queue[0], nat), nat
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
