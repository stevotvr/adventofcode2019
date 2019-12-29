package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type state struct {
	room  string
	items string
}

type holder struct {
	state    state
	computer *computer
	log      []string
}

type computer struct {
	program  []int
	position int
	offset   int
	channel  chan int
	waiting  chan bool
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

func createComputer() *computer {
	input, _ := ioutil.ReadFile("input.txt")
	ints := make([]int, 0)
	for _, v := range strings.Split(string(input), ",") {
		i, _ := strconv.Atoi(v)
		ints = append(ints, i)
	}

	c := computer{program: ints, channel: make(chan int), waiting: make(chan bool)}
	go c.run()

	return &c
}

func (c *computer) clone() *computer {
	ints := make([]int, len(c.program))
	copy(ints, c.program)

	c2 := computer{program: ints, position: c.position, offset: c.offset, channel: make(chan int), waiting: make(chan bool)}
	go c2.run()

	return &c2
}

func play() {
	c := createComputer()

	inreader := bufio.NewReader(os.Stdin)
	for {
		select {
		case char := <-c.channel:
			fmt.Printf("%c", char)
		case _, ok := <-c.waiting:
			if !ok {
				return
			}

			input, _ := inreader.ReadString('\n')
			for _, inchar := range strings.Trim(input, "\r\n ") {
				c.channel <- int(inchar)
				<-c.waiting
			}

			c.channel <- '\n'
		}
	}
}

func solve() {
	reRoom, _ := regexp.Compile("== ([a-zA-Z ]+) ==")
	reDoors, _ := regexp.Compile("Doors here lead:\n((- [a-z]+\n)+)")
	reItems, _ := regexp.Compile("Items here:\n- ([a-z ]+)\n")
	reCode, _ := regexp.Compile(" ([0-9]+) ")

	queue, visited := make([]holder, 1), make(map[state]bool)
	queue[0] = holder{
		state:    state{},
		computer: createComputer(),
		log:      make([]string, 0),
	}

	var h holder
	for {
		h, queue = queue[0], queue[1:]
		outputBuilder := strings.Builder{}

	readloop:
		for {
			select {
			case c := <-h.computer.channel:
				outputBuilder.WriteRune(rune(c))
			case <-h.computer.waiting:
				break readloop
			}
		}

		output := []byte(outputBuilder.String())

		if h.state.room == "Security Checkpoint" && strings.Contains(string(output), "You may proceed.") {
			fmt.Println("Commands:")
			for _, c := range h.log {
				fmt.Println(c)
			}

			fmt.Println()
			fmt.Print("Password: ")

			m := reCode.FindSubmatch(output)
			fmt.Println(string(m[1]))

			return
		}

		visited[h.state] = true

		matchRoom := reRoom.FindSubmatch(output)
		matchDoors := reDoors.FindSubmatch(output)
		matchItems := reItems.FindSubmatch(output)

		for _, door := range strings.Split(string(matchDoors[1][2:]), "- ") {
			door = door[:len(door)-1]
			nextState := state{room: string(matchRoom[1]), items: h.state.items}
			if !visited[nextState] {
				nextComputer := *h.computer.clone()
				<-nextComputer.waiting

				for _, c := range door {
					nextComputer.channel <- int(c)
					<-nextComputer.waiting
				}

				nextComputer.channel <- '\n'

				log := h.log
				log = append(log, door)
				queue = append(queue, holder{state: nextState, computer: &nextComputer, log: log})
			}

			if len(matchItems) > 1 {
				item := string(matchItems[1])
				if item == "infinite loop" || item == "giant electromagnet" {
					continue
				}

				nextItems := strings.Split(h.state.items, ",")
				if nextItems[0] == "" {
					nextItems = nextItems[1:]
				}

				nextItems = append(nextItems, string(matchItems[1]))
				sort.Strings(nextItems)

				nextState = state{room: string(matchRoom[1]), items: strings.Join(nextItems, ",")}
				if !visited[nextState] {
					nextComputer := *h.computer.clone()
					<-nextComputer.waiting

					take := "take " + string(matchItems[1])
					for _, c := range take {
						nextComputer.channel <- int(c)
						<-nextComputer.waiting
					}

					nextComputer.channel <- '\n'

					var ok bool
				readloop2:
					for {
						select {
						case <-nextComputer.channel:
						case _, ok = <-nextComputer.waiting:
							break readloop2
						}
					}

					if !ok {
						continue
					}

					for _, c := range door {
						nextComputer.channel <- int(c)
						<-nextComputer.waiting
					}

					nextComputer.channel <- '\n'

					log := h.log
					log = append(log, take, door)
					queue = append(queue, holder{state: nextState, computer: &nextComputer, log: log})
				}
			}
		}
	}
}

func main() {
	var mode string
	if len(os.Args) == 2 {
		fmt.Sscan(os.Args[1], &mode)
	} else {
		fmt.Print("Enter mode (play or solve): ")
		fmt.Scanf("%s\n", &mode)
	}

	switch mode {
	case "play":
		play()
	case "solve":
		solve()
	default:
		fmt.Println("Error: Invalid mode.")
	}
}
