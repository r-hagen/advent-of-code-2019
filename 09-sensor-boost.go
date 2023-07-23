package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Memory []int64
type Input []int64
type Signal []int64
type Output []int64

type Computer struct {
	p   int64
	rel int64
	mem Memory
	in  Input
	out Output
}

func (m Memory) Copy() Memory {
	cpy := make(Memory, 0, len(m))
	return append(cpy, m...)
}

func MakeComputer(mem Memory, in Input) *Computer {
	memory := mem.Copy()
	for i := 0; i <= 10_000; i++ {
		memory = append(memory, 0)
	}
	return &Computer{0, 0, memory, in, Output{}}
}

func (c *Computer) Opcode() int64 {
	return c.mem[c.p]
}

func (c *Computer) Op() int64 {
	return c.Opcode() % 100
}

func (c *Computer) Mode(n int64) int64 {
	return (c.Opcode() / int64(math.Pow(float64(10), float64(n+1)))) % 10
}

func (c *Computer) Param(n int64) int64 {
	mode := c.Mode(n)
	offset := c.p + n
	if mode == 0 {
		address := c.mem[offset] // position
		return c.mem[address]
	} else if mode == 1 { // immediate
		return c.mem[offset]
	} else if mode == 2 { // relative
		address := c.mem[offset] + c.rel
		return c.mem[address]
	} else {
		panic("unknown parameter mode")
	}
}

func (c *Computer) Address(n int64) int64 {
	mode := c.Mode(n)
	offset := c.p + n
	if mode == 0 || mode == 1 {
		return c.mem[offset]
	} else if mode == 2 {
		return c.mem[offset] + c.rel
	} else {
		panic("unknown parameter mode")
	}
}

func (c *Computer) Signal() int64 {
	return c.out[0]
}

func (c *Computer) Halted() bool {
	return c.mem[c.p] == 99
}

func (c *Computer) Run(s Signal) {
	c.out = Output{}
	c.in = append(c.in, s...)

	for {
		op := c.Op()

		if op == 1 { // add
			p1 := c.Param(1)
			p2 := c.Param(2)
			p3 := c.Address(3)
			add := p1 + p2
			c.mem[p3] = add
			c.p += 4
		} else if op == 2 { // mul
			p1 := c.Param(1)
			p2 := c.Param(2)
			p3 := c.Address(3)
			mul := p1 * p2
			c.mem[p3] = mul
			c.p += 4
		} else if op == 3 { // in
			p1 := c.Address(1)
			if len(c.in) == 0 {
				break
			}
			c.mem[p1] = c.in[0]
			if len(c.in) > 1 {
				c.in = c.in[1:]
			} else {
				c.in = Input{}
			}
			c.p += 2
		} else if op == 4 { // out
			p1 := c.Param(1)
			c.out = append(c.out, p1)
			c.p += 2
		} else if op == 5 { // jt
			p1 := c.Param(1)
			p2 := c.Param(2)
			if p1 != 0 {
				c.p = p2
			} else {
				c.p += 3
			}
		} else if op == 6 { // jf
			p1 := c.Param(1)
			p2 := c.Param(2)
			if p1 == 0 {
				c.p = p2
			} else {
				c.p += 3
			}
		} else if op == 7 { // lt
			p1 := c.Param(1)
			p2 := c.Param(2)
			p3 := c.Address(3)
			if p1 < p2 {
				c.mem[p3] = 1
			} else {
				c.mem[p3] = 0
			}
			c.p += 4
		} else if op == 8 { // eq
			p1 := c.Param(1)
			p2 := c.Param(2)
			p3 := c.Address(3)
			if p1 == p2 {
				c.mem[p3] = 1
			} else {
				c.mem[p3] = 0
			}
			c.p += 4
		} else if op == 9 { // rel
			p1 := c.Param(1)
			c.rel += p1
			c.p += 2
		} else if op == 99 {
			break
		} else {
			fmt.Println(c.Opcode())
			panic("unexpected opcode")
		}
	}
}

func parseProgram(line string) Memory {
	line = strings.Trim(line, "\n")
	nums := strings.Split(line, ",")

	program := Memory{}

	for _, num := range nums {
		i, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		program = append(program, int64(i))
	}

	return program
}

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')

	program := parseProgram(line)

	com1 := MakeComputer(program, Input{1})
	com1.Run(Signal{})
	fmt.Println("ans1", com1.out)

	com2 := MakeComputer(program, Input{2})
	com2.Run(Signal{})
	fmt.Println("ans2", com2.out)
}

func Permute(arr []int) [][]int {
	res := [][]int{}

	// Heap's algorithm
	// https://en.wikipedia.org/wiki/Heap%27s_algorithm
	var generate func(int, []int)

	generate = func(k int, A []int) {
		if k == 1 {
			tmp := make([]int, len(A))
			copy(tmp, A)
			res = append(res, tmp)
		} else {
			generate(k-1, A)
			for i := 0; i < k-1; i++ {
				if k%2 == 0 {
					tmp := A[i]
					A[i] = A[k-1]
					A[k-1] = tmp
				} else {
					tmp := A[0]
					A[0] = A[k-1]
					A[k-1] = tmp
				}
				generate(k-1, A)
			}
		}
	}

	generate(len(arr), arr)

	return res
}
