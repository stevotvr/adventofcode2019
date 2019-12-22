package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func newDeck(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}

	return s
}

func deal(s []int) []int {
	s2 := make([]int, len(s))
	for i := range s {
		s2[len(s)-1-i] = s[i]
	}

	return s2
}

func cut(s []int, n int) []int {
	if n >= 0 {
		return append(s[n:], s[:n]...)
	}

	n = -n
	return append(s[len(s)-n:], s[:len(s)-n]...)
}

func deali(s []int, n int) []int {
	s2 := make([]int, len(s))
	for i := range s {
		s2[(i*n)%len(s)] = s[i]
	}

	return s2
}

func part1() {
	input, _ := ioutil.ReadFile("input.txt")
	s := newDeck(10007)
	for _, line := range strings.Split(string(input), "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case line == "deal into new stack":
			s = deal(s)
		case line[:3] == "cut":
			i, _ := strconv.Atoi(line[4:])
			s = cut(s, i)
		case line[:19] == "deal with increment":
			i, _ := strconv.Atoi(line[20:])
			s = deali(s, i)
		}
	}

	for i, v := range s {
		if v == 2019 {
			fmt.Println(i)
			return
		}
	}
}

func part2() {
	input, _ := ioutil.ReadFile("input.txt")

	n, iter := big.NewInt(119315717514047), big.NewInt(101741582076661)
	offset, increment := big.NewInt(0), big.NewInt(1)
	for _, op := range strings.Split(string(input), "\n") {
		op = strings.TrimRight(op, "\r")
		switch {
		case op == "deal into new stack":
			increment.Mul(increment, big.NewInt(-1))
			offset.Add(offset, increment)
		case op[:3] == "cut":
			i, _ := strconv.Atoi(op[4:])
			offset.Add(offset, big.NewInt(0).Mul(big.NewInt(int64(i)), increment))
		case op[:19] == "deal with increment":
			i, _ := strconv.Atoi(op[20:])
			increment.Mul(increment, big.NewInt(0).Exp(big.NewInt(int64(i)), big.NewInt(0).Sub(n, big.NewInt(2)), n))
		}
	}

	finalIncr := big.NewInt(0).Exp(increment, iter, n)

	finalOffs := big.NewInt(0).Exp(increment, iter, n)
	finalOffs.Sub(big.NewInt(1), finalOffs)
	invmod := big.NewInt(0).Exp(big.NewInt(0).Sub(big.NewInt(1), increment), big.NewInt(0).Sub(n, big.NewInt(2)), n)
	finalOffs.Mul(finalOffs, invmod)
	finalOffs.Mul(finalOffs, offset)

	answer := big.NewInt(0).Mul(big.NewInt(2020), finalIncr)
	answer.Add(answer, finalOffs)
	answer.Mod(answer, n)

	fmt.Println(answer)
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
