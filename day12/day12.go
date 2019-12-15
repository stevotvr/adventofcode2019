package main

import (
	"fmt"
	"math"
	"os"
)

type vector3d struct {
	x int
	y int
	z int
}

type moon struct {
	pos vector3d
	vel vector3d
}

func (m *moon) Move() {
	m.pos.x += m.vel.x
	m.pos.y += m.vel.y
	m.pos.z += m.vel.z
}

func (m *moon) Energy() int {
	pos := math.Abs(float64(m.pos.x)) + math.Abs(float64(m.pos.y)) + math.Abs(float64(m.pos.z))
	vel := math.Abs(float64(m.vel.x)) + math.Abs(float64(m.vel.y)) + math.Abs(float64(m.vel.z))
	return int(pos * vel)
}

func getMoons() []moon {
	input, _ := os.Open("input.txt")
	moons := make([]moon, 0)
	x, y, z := 0, 0, 0
	for {
		_, err := fmt.Fscanf(input, "<x=%d, y=%d, z=%d>", &x, &y, &z)
		if err != nil {
			break
		}

		moons = append(moons, moon{pos: vector3d{x, y, z}, vel: vector3d{0, 0, 0}})
		fmt.Fscanln(input)
	}

	return moons
}

func applyGravity(moons []moon) {
	for m1 := range moons {
		for m2 := range moons {
			if m1 == m2 {
				continue
			}

			if moons[m1].pos.x > moons[m2].pos.x {
				moons[m1].vel.x--
			} else if moons[m1].pos.x < moons[m2].pos.x {
				moons[m1].vel.x++
			}

			if moons[m1].pos.y > moons[m2].pos.y {
				moons[m1].vel.y--
			} else if moons[m1].pos.y < moons[m2].pos.y {
				moons[m1].vel.y++
			}

			if moons[m1].pos.z > moons[m2].pos.z {
				moons[m1].vel.z--
			} else if moons[m1].pos.z < moons[m2].pos.z {
				moons[m1].vel.z++
			}
		}
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func part1() {
	moons := getMoons()
	for i := 0; i < 1000; i++ {
		applyGravity(moons)

		for m := range moons {
			moons[m].Move()
		}
	}

	total := 0
	for m := range moons {
		total += moons[m].Energy()
	}

	fmt.Println(total)
}

func part2() {
	moons := getMoons()

	period, found := vector3d{}, 0
	for i := 1; found < 3; i++ {
		applyGravity(moons)

		for m := range moons {
			moons[m].Move()
		}

		if period.x == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.x != 0 {
					inv = true
					break
				}
			}

			if !inv {
				period.x = i * 2
				found++
			}
		}

		if period.y == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.y != 0 {
					inv = true
					break
				}
			}

			if !inv {
				period.y = i * 2
				found++
			}
		}

		if period.z == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.z != 0 {
					inv = true
					break
				}
			}

			if !inv {
				period.z = i * 2
				found++
			}
		}
	}

	fmt.Println(lcm(lcm(period.x, period.y), period.z))
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
