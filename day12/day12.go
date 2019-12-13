package main

import (
	"fmt"
	"math"
	"os"
)

type Vector3d struct {
	x int
	y int
	z int
}

type Moon struct {
	pos Vector3d
	vel Vector3d
}

func (m *Moon) Move() {
	m.pos.x += m.vel.x
	m.pos.y += m.vel.y
	m.pos.z += m.vel.z
}

func (m *Moon) Energy() int {
	pos := math.Abs(float64(m.pos.x)) + math.Abs(float64(m.pos.y)) + math.Abs(float64(m.pos.z))
	vel := math.Abs(float64(m.vel.x)) + math.Abs(float64(m.vel.y)) + math.Abs(float64(m.vel.z))
	return int(pos * vel)
}

func getMoons() []Moon {
	input, _ := os.Open("input.txt")
	moons := make([]Moon, 0)
	x, y, z := 0, 0, 0
	for {
		_, err := fmt.Fscanf(input, "<x=%d, y=%d, z=%d>", &x, &y, &z)
		if err != nil {
			break
		}

		moons = append(moons, Moon{pos: Vector3d{x, y, z}, vel: Vector3d{0, 0, 0}})
		fmt.Fscanln(input)
	}

	return moons
}

func applyGravity(moons []Moon) {
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

func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func Lcm(a, b int) int {
	return a * b / Gcd(a, b)
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

	fmt.Println(moons)
	fmt.Println(total)
}

func part2() {
	moons := getMoons()

	counts, found := make([]int, 3), 0
	for i := 1; found < 3; i++ {
		applyGravity(moons)

		for m := range moons {
			moons[m].Move()
		}

		if counts[0] == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.x != 0 {
					inv = true
					break
				}
			}

			if !inv {
				counts[0] = i * 2
				found++
			}
		}

		if counts[1] == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.y != 0 {
					inv = true
					break
				}
			}

			if !inv {
				counts[1] = i * 2
				found++
			}
		}

		if counts[2] == 0 {
			inv := false
			for m := range moons {
				if moons[m].vel.z != 0 {
					inv = true
					break
				}
			}

			if !inv {
				counts[2] = i * 2
				found++
			}
		}
	}

	lcm := Lcm(Lcm(counts[0], counts[1]), counts[2])

	fmt.Println(lcm)
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
