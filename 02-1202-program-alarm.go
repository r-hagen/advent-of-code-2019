package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	nums := strings.Split(line, ",")

	program := []int{}
	for _, num := range nums {
		i, _ := strconv.Atoi(num)
		program = append(program, i)
	}

	memory := make([]int, len(program))
	copy(memory, program)

	memory[1] = 12
	memory[2] = 2
	ans1 := run_program(memory)
	fmt.Println("ans1", ans1)

	for noun := 12; noun <= 99; noun++ {
		for verb := 2; verb <= 99; verb++ {
			copy(memory, program)
			memory[1] = noun
			memory[2] = verb

			res := run_program(memory)

			if res == 19690720 {
				ans2 := noun*100 + verb
				fmt.Println("ans2", ans2)
				break
			}
		}
	}
}

func run_program(memory []int) int {
	p := 0
	for {
		op := memory[p]
		if op == 99 {
			break
		}
		if op == 1 {
			x1 := memory[p+1]
			x2 := memory[p+2]
			y := memory[p+3]
			sum := memory[x1] + memory[x2]
			memory[y] = sum
		} else if op == 2 {
			x1 := memory[p+1]
			x2 := memory[p+2]
			y := memory[p+3]
			prod := memory[x1] * memory[x2]
			memory[y] = prod
		} else {
			log.Fatal("unexpected opcode", op)
		}
		p += 4
	}
	return memory[0]
}
