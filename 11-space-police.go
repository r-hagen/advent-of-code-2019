package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Panels map[Location]int64

type Direction string

const (
	Up    Direction = "UP"
	Down  Direction = "DOWN"
	Left  Direction = "LEFT"
	Right Direction = "RIGHT"
)

type Location struct {
	x int
	y int
}

type Robot struct {
	pos    Location
	dir    Direction
	panels Panels
}

func MakeRobot() *Robot {
	return &Robot{Location{0, 0}, Up, make(Panels)}
}

func (r *Robot) Camera() Input {
	color := r.panels[r.pos]
	return Input{int64(color)}
}

func (r *Robot) Paint(output Output) {
	color := output[0]
	r.panels[r.pos] = color

	turn := output[1]
	switch turn {
	case 0: // left 90deg
		switch r.dir {
		case Up:
			r.dir = Left
			r.pos = Location{r.pos.x - 1, r.pos.y}
		case Left:
			r.dir = Down
			r.pos = Location{r.pos.x, r.pos.y - 1}
		case Down:
			r.dir = Right
			r.pos = Location{r.pos.x + 1, r.pos.y}
		case Right:
			r.dir = Up
			r.pos = Location{r.pos.x, r.pos.y + 1}
		default:
			panic("invalid direction")
		}
	case 1: // right 90deg
		switch r.dir {
		case Up:
			r.dir = Right
			r.pos = Location{r.pos.x + 1, r.pos.y}
		case Left:
			r.dir = Up
			r.pos = Location{r.pos.x, r.pos.y + 1}
		case Down:
			r.dir = Left
			r.pos = Location{r.pos.x - 1, r.pos.y}
		case Right:
			r.dir = Down
			r.pos = Location{r.pos.x, r.pos.y - 1}
		default:
			panic("invalid direction")
		}
	default:
		panic("invalid turn")
	}
}

func (r *Robot) Plot() {
	b := make([]int, 4)

	for l, v := range r.panels {
		if v == 0 {
			continue
		}

		if l.x < b[0] {
			b[0] = l.x
		}
		if l.x > b[1] {
			b[1] = l.x
		}
		if l.y < b[2] {
			b[2] = l.y
		}
		if l.y > b[3] {
			b[3] = l.y
		}
	}

	for y := b[3]; y >= b[2]; y-- {
		for x := b[0]; x <= b[1]; x++ {
			color := r.panels[Location{x, y}]
			if color == 0 {
				fmt.Print(" ")
			}
			if color == 1 {
				fmt.Print("#")
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

	robot1 := MakeRobot()
	computer1 := MakeComputer(program, Input{})
	for !computer1.Halted() {
		computer1.Run(robot1.Camera())
		robot1.Paint(computer1.out)
	}
	fmt.Println("ans1", len(robot1.panels))

	robot2 := MakeRobot()
	robot2.panels[robot2.pos] = 1
	computer2 := MakeComputer(program, Input{})
	for !computer2.Halted() {
		computer2.Run(robot2.Camera())
		robot2.Paint(computer2.out)
	}
	fmt.Println("ans2")
	robot2.Plot()
}
