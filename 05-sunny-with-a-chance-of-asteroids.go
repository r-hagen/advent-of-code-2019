package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	nums := strings.Split(line, ",")

	program := []int{}
	for _, num := range nums {
		i, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		program = append(program, i)
	}

	memory := make([]int, len(program))

	copy(memory, program)
	ans1 := run_program(memory, 1)
	fmt.Println("ans1", ans1)

	copy(memory, program)
	ans2 := run_program(memory, 5)
	fmt.Println("ans2", ans2)
}

func run_program(memory []int, input int) int {
	p := 0
	output := 0
	for {
		opcode := memory[p]
		op := opcode % 100

		param := func(nth int) int {
			mode := (opcode / int(math.Pow(float64(10), float64(nth+1)))) % 10
			offset := p + nth
			if mode == 0 {
				address := memory[offset]
				return memory[address]
			} else {
				return memory[offset]
			}
		}

		address := func(nth int) int {
			return memory[p+nth]
		}

		if op == 1 { // addition
			p1 := param(1)
			p2 := param(2)
			p3 := address(3)
			sum := p1 + p2
			memory[p3] = sum
			p += 4
		} else if op == 2 { // multiplication
			p1 := param(1)
			p2 := param(2)
			p3 := address(3)
			mul := p1 * p2
			memory[p3] = mul
			p += 4
		} else if op == 3 { // input
			p1 := address(3)
			memory[p1] = input
			p += 2
		} else if op == 4 { // output
			output = param(1)
			fmt.Println("output", output)
			p += 2
		} else if op == 5 { // jump if true
			p1 := param(1)
			p2 := param(2)
			if p1 != 0 {
				p = p2
			} else {
				p += 3
			}
		} else if op == 6 { // jump if false
			p1 := param(1)
			p2 := param(2)
			if p1 == 0 {
				p = p2
			} else {
				p += 3
			}
		} else if op == 7 { // less than
			p1 := param(1)
			p2 := param(2)
			p3 := address(3)
			if p1 < p2 {
				memory[p3] = 1
			} else {
				memory[p3] = 0
			}
			p += 4
		} else if op == 8 { // equals
			p1 := param(1)
			p2 := param(2)
			p3 := memory[p+3]
			if p1 == p2 {
				memory[p3] = 1
			} else {
				memory[p3] = 0
			}
			p += 4
		} else if op == 99 {
			break
		} else {
			fmt.Println(opcode)
			panic("unexpected opcode")
		}
	}
	return output
}
