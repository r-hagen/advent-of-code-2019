package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Grid map[Location]int64

type Location struct {
	x int64
	y int64
}

const (
	EMPTY  int64 = 0
	WALL   int64 = 1
	BLOCK  int64 = 2
	PADDLE int64 = 3
	BALL   int64 = 4
)

type Arcade struct {
	grid   Grid
	score  int64
	paddle int64
	ball   Location
	run  int64
}

func MakeArcade() *Arcade {
	return &Arcade{make(Grid), 0, 0, Location{0, 0}, 0}
}

func (r *Arcade) HasBlocks() bool {
	if r.run == 0 {
		return true
	}

	for _, v := range r.grid {
		if v == BLOCK {
			return true
		}
	}

	return false
}

func (r *Arcade) Joystick() Input {
	px := r.paddle
	bx := r.ball.x

	dx := bx - px

	if dx == 0 {
		return Input{0}
	} else if dx < 0 {
		return Input{-1}
	}
	return Input{1}
}

func (r *Arcade) Update(output Output) {
	r.run += 1

	for i := 0; i < len(output); i = i + 3 {
		x, y := output[i], output[i+1]
		v := output[i+2]

		if x == -1 && y == 0 {
			r.score = v
		} else {
			if v == PADDLE {
				r.paddle = x
			} else if v == BALL {
				r.ball = Location{x, y}
			}
			r.grid[Location{x, y}] = v
		}
	}
}

func (r *Arcade) Plot() {
	// xmin, xmax, ymin, ymax
	bounds := make([]int64, 4)
	for loc, val := range r.grid {
		if val == 0 {
			continue
		}
		if loc.x < bounds[0] {
			bounds[0] = loc.x
		}
		if loc.x > bounds[1] {
			bounds[1] = loc.x
		}
		if loc.y < bounds[2] {
			bounds[2] = loc.y
		}
		if loc.y > bounds[3] {
			bounds[3] = loc.y
		}
	}

	fmt.Println("\nRun: ", r.run)
	for y := bounds[3]; y >= bounds[2]; y-- {
		for x := bounds[0]; x <= bounds[1]; x++ {
			tile := r.grid[Location{x, y}]
			if tile == 0 {
				fmt.Print(" ")
			}
			if tile == 1 {
				fmt.Print("|")
			}
			if tile == 2 {
				fmt.Print("#")
			}
			if tile == 3 {
				fmt.Print("-")
			}
			if tile == 4 {
				fmt.Print("o")
			}
		}
		fmt.Print("\n")
	}
}

type Memory []int64

type Input []int64

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

func (c *Computer) Run(s Input) {
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
	input, _ := reader.ReadString('\n')
	program := parseProgram(input)

	computer1 := MakeComputer(program, Input{})
	computer1.Run(Input{})
	ans1 := 0
	for i := 2; i < len(computer1.out); i = i + 3 {
		id := computer1.out[i]
		if id == BLOCK {
			ans1 += 1
		}
	}
	fmt.Println("ans1", ans1)

	arcade2 := MakeArcade()
	computer2 := MakeComputer(program, Input{})
	computer2.mem[0] = 2

	for arcade2.HasBlocks() {
		computer2.Run(arcade2.Joystick())
		arcade2.Update(computer2.out)

		if arcade2.run == 1 || arcade2.run%1000 == 0 {
			arcade2.Plot()
		}
	}
	arcade2.Plot()
	fmt.Println("ans2", arcade2.score)
}
