package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Memory []int

type Input []int

type Signal []int

type Output []int

type Computer struct {
	p   int
	mem Memory
	in  Input
	out Output
}

func (m Memory) Copy() Memory {
	cpy := make(Memory, 0, len(m))
	return append(cpy, m...)
}

func MakeComputer(mem Memory, in Input) *Computer {
	p := Computer{0, mem.Copy(), in, Output{}}
	return &p
}

func (c *Computer) Opcode() int {
	return c.mem[c.p]
}

func (c *Computer) Op() int {
	return c.Opcode() % 100
}

func (c *Computer) Param(n int) int {
	mode := (c.Opcode() / int(math.Pow(float64(10), float64(n+1)))) % 10
	offset := c.p + n
	if mode == 0 {
		address := c.mem[offset]
		return c.mem[address]
	} else {
		return c.mem[offset]
	}
}

func (c *Computer) Address(n int) int {
	return c.mem[c.p+n]
}

func (c *Computer) Signal() int {
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
				c.in = []int{}
			}
			c.p += 2
		} else if op == 4 { // out
			c.out = append(c.out, c.Param(1))
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
		} else if op == 99 {
			break
		} else {
			fmt.Println(c.Opcode())
			panic("unexpected opcode")
		}
	}
}

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	nums := strings.Split(line, ",")

	program := Memory{}
	for _, num := range nums {
		i, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		program = append(program, i)
	}

	// part 1
	ans1 := 0
	phases := Permute([]int{0, 1, 2, 3, 4})
	for _, phase := range phases {
		amp1 := MakeComputer(program, Input{phase[0]})
		amp2 := MakeComputer(program, Input{phase[1]})
		amp3 := MakeComputer(program, Input{phase[2]})
		amp4 := MakeComputer(program, Input{phase[3]})
		amp5 := MakeComputer(program, Input{phase[4]})

		amp1.Run(Signal{0})
		amp2.Run(Signal(amp1.out))
		amp3.Run(Signal(amp2.out))
		amp4.Run(Signal(amp3.out))
		amp5.Run(Signal(amp4.out))

		if amp5.Signal() > ans1 {
			ans1 = amp5.Signal()
		}
	}
	fmt.Println("ans1", ans1)

	// part 2
	ans2 := 0
	phases = Permute([]int{5, 6, 7, 8, 9})
	for _, phase := range phases {
		amp1 := MakeComputer(program, Input{phase[0]})
		amp2 := MakeComputer(program, Input{phase[1]})
		amp3 := MakeComputer(program, Input{phase[2]})
		amp4 := MakeComputer(program, Input{phase[3]})
		amp5 := MakeComputer(program, Input{phase[4]})

		amp5.out = append(amp5.out, 0)

		for {
			amp1.Run(Signal(amp5.out))
			amp2.Run(Signal(amp1.out))
			amp3.Run(Signal(amp2.out))
			amp4.Run(Signal(amp3.out))
			amp5.Run(Signal(amp4.out))

			if amp5.Halted() {
				if amp5.Signal() > ans2 {
					ans2 = amp5.Signal()
				}
				break
			}
		}
	}
	fmt.Println("ans2", ans2)
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
