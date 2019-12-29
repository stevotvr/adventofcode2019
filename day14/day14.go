package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type chemical struct {
	name     string
	quantity int
}

type reaction struct {
	product  chemical
	reagents []chemical
}

func getReactions() map[string]reaction {
	reactions := make(map[string]reaction)
	input, _ := ioutil.ReadFile("input.txt")

	for _, line := range strings.Split(string(input), "\n") {
		parts := strings.Split(line, " => ")
		name, quantity := "", 0
		fmt.Sscanf(parts[1], "%d %s", &quantity, &name)
		product := chemical{name: name, quantity: quantity}

		reagents := make([]chemical, 0)
		for _, s := range strings.Split(parts[0], ", ") {
			fmt.Sscanf(s, "%d %s", &quantity, &name)
			reagents = append(reagents, chemical{name: name, quantity: quantity})
		}

		reactions[product.name] = reaction{product: product, reagents: reagents}
	}

	return reactions
}

func orePerFuel(reactions map[string]reaction, quantity int) int {
	chems, excess, raw := make(map[string]chemical), make(map[string]int), make(map[string]chemical)
	chems["FUEL"] = chemical{name: "FUEL", quantity: quantity}
	for len(chems) > 0 {
		newchems := make(map[string]chemical)

		for _, c := range chems {
			react := reactions[c.name]

			if react.reagents[0].name == "ORE" {
				c2, ok := raw[c.name]
				if ok {
					raw[c.name] = chemical{name: c.name, quantity: c.quantity + c2.quantity}
				} else {
					raw[c.name] = chemical{name: c.name, quantity: c.quantity}
				}

				continue
			}

			if excess[c.name] > 0 {
				c.quantity -= excess[c.name]
				if c.quantity <= 0 {
					excess[c.name] = -c.quantity
					continue
				}

				excess[c.name] = 0
			}

			units := c.quantity / react.product.quantity
			if c.quantity%react.product.quantity != 0 {
				units++
			}

			for _, r := range react.reagents {
				needs := units * r.quantity

				c2, ok := newchems[r.name]
				if ok {
					newchems[r.name] = chemical{name: r.name, quantity: needs + c2.quantity}
				} else {
					newchems[r.name] = chemical{name: r.name, quantity: needs}
				}
			}

			excess[c.name] += react.product.quantity*units - c.quantity
		}

		chems = newchems
	}

	ore := 0
	for _, r := range raw {
		react := reactions[r.name]

		q := react.product.quantity
		if r.quantity%q == 0 {
			q = r.quantity / q
		} else {
			q = r.quantity/q + 1
		}

		ore += react.reagents[0].quantity * q
	}

	return ore
}

func part1() {
	reactions := getReactions()
	fmt.Println(orePerFuel(reactions, 1))
}

func part2() {
	reactions := getReactions()

	fuel, incr := 1, 1000000
	for {
		needs := orePerFuel(reactions, fuel)
		if needs > 1000000000000 {
			if incr == 1 {
				fuel--
				break
			}

			fuel -= incr
			incr /= 10
		}

		fuel += incr
	}

	fmt.Println(fuel)
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
