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

func getMap(chars []int) (map[vector2d]bool, vector2d, int) {
	grid, pos, dir := make(map[vector2d]bool), vector2d{}, 0
	x, y := 0, 0
	for _, i := range chars {
		if i == 10 {
			x = 0
			y++
			continue
		}

		if i == '#' || i == '^' || i == '>' || i == 'v' || i == '<' {
			grid[vector2d{x, y}] = true
			if i == '^' || i == '>' || i == 'v' || i == '<' {
				pos = vector2d{x, y}
				switch i {
				case '^':
					dir = 0
				case '>':
					dir = 1
				case 'v':
					dir = 2
				case '<':
					dir = 3
				}
			}
		}
		x++
	}

	return grid, pos, dir
}

func breakInstructions(in string) []string {
	for l1 := 2; ; l1++ {
		if in[l1-1] != ',' {
			continue
		}

		off1 := l1
		p1 := in[:l1]

		for p1 == in[off1:off1+l1] {
			off1 += l1
		}

		for l2 := 2; l2 < 20; l2++ {
			if in[off1+l2-1] != ',' {
				continue
			}

			off2 := off1 + l2
			p2 := in[off1 : off1+l2]

			for {
				if p1 == in[off2:off2+l1] {
					off2 += l1
					continue
				}

				if p2 == in[off2:off2+l2] {
					off2 += l2
					continue
				}

				break
			}

			for l3 := 2; off2+l3 < len(in); l3++ {
				p3 := in[off2 : off2+l3]
				if p3[len(p3)-1] == ',' {
					p3 = p3[:len(p3)-1]
				}

				tmp := strings.ReplaceAll(in, p1[:len(p1)-1], "")
				tmp = strings.ReplaceAll(tmp, p2[:len(p2)-1], "")
				tmp = strings.ReplaceAll(tmp, p3, "")
				tmp = strings.ReplaceAll(tmp, ",", "")

				if len(tmp) == 0 {
					keys := strings.ReplaceAll(in, p1[:len(p1)-1], "A")
					keys = strings.ReplaceAll(keys, p2[:len(p2)-1], "B")
					keys = strings.ReplaceAll(keys, p3, "C")
					return []string{keys, p1[:len(p1)-1], p2[:len(p2)-1], p3}
				}
			}
		}
	}
}

func part1() {
	c := createComputer()
	chars := make([]int, 0)
	for i := range c.channel {
		chars = append(chars, i)
	}

	grid, _, _ := getMap(chars)
	directions := []vector2d{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	total := 0
	for pos, val := range grid {
		if !val {
			continue
		}

		for d := range directions {
			a1, a2, a3 := directions[d], directions[(d+1)%4], directions[(d+2)%4]
			if grid[vector2d{pos.x + a1.x, pos.y + a1.y}] && grid[vector2d{pos.x + a2.x, pos.y + a2.y}] && grid[vector2d{pos.x + a3.x, pos.y + a3.y}] {
				total += pos.x * pos.y
				break
			}
		}
	}

	fmt.Println(total)
}

func part2() {
	c := createComputer()
	c.program[0] = 2
	chars := make([]int, 0)
	prev := 0
	for {
		i := <-c.channel
		if prev == 10 && i == 10 {
			break
		}

		chars = append(chars, i)
		prev = i
	}

	grid, pos, dir := getMap(chars)
	directions := []vector2d{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	path, steps := make([]string, 0), 0
	for {
		forward := directions[dir]
		forward = vector2d{pos.x + forward.x, pos.y + forward.y}
		if grid[forward] {
			steps++
			pos = forward
			continue
		}

		right := directions[(dir+1)%4]
		right = vector2d{pos.x + right.x, pos.y + right.y}
		if grid[right] {
			path = append(path, strconv.Itoa(steps), "R")
			steps = 1
			dir = (dir + 1) % 4
			pos = right
			continue
		}

		left := directions[(dir+3)%4]
		left = vector2d{pos.x + left.x, pos.y + left.y}
		if grid[left] {
			path = append(path, strconv.Itoa(steps), "L")
			steps = 1
			dir = (dir + 3) % 4
			pos = left
			continue
		}

		path = append(path, strconv.Itoa(steps))
		break
	}

	lines := breakInstructions(strings.Join(path[1:], ","))

	for _, line := range lines {
		for <-c.channel != 10 {
		}

		for _, char := range line {
			c.channel <- int(char)
		}

		c.channel <- 10
	}

	for <-c.channel != 10 {
	}

	c.channel <- int('n')
	c.channel <- 10

	output := 0
	for {
		i, ok := <-c.channel
		if !ok {
			break
		}

		output = i
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
