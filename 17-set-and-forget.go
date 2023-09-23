package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ---------- INTCODE COMPUTER START ----------

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

func MakeComputer(mem Memory, in Input) Computer {
	memory := mem.Copy()
	for i := 0; i <= 10_000; i++ {
		memory = append(memory, 0)
	}
	return Computer{0, 0, memory, in, Output{}}
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

// ---------- INTCODE COMPUTER END ----------

type Grid = map[Point]int64

func display(grid Grid) {
	xmax := 0
	ymax := 0

	for p := range grid {
		if p.x > xmax {
			xmax = p.x
		}
		if p.y > ymax {
			ymax = p.y
		}
	}

	for y := 0; y <= ymax; y++ {
		for x := 0; x <= xmax; x++ {
			v, ok := grid[Point{x, y}]

			if !ok {
				panic("invalid coordinate")
			}

			fmt.Print(string(v))
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
}

type Point struct {
	x int
	y int
}

type Robot struct {
	x int
	y int
	d int64
}

const (
	SCAFFOLD     int64 = 35
	OPENSPACE    int64 = 46
	NEWLINE      int64 = 10
	INTERSECTION int64 = 111
	UP           int64 = 94
	DOWN         int64 = 118
	LEFT         int64 = 60
	RIGHT        int64 = 62
)

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	input, _ := reader.ReadString('\n')
	program := parseProgram(input)

	computer := MakeComputer(program, Input{})

	x, y := 0, 0
	grid := make(Grid)
	// robot := Robot{0, 0, 0}

	computer.Run(Input{})

	for _, v := range computer.out {
		if v == SCAFFOLD || v == OPENSPACE {
			grid[Point{x, y}] = v
			x = x + 1
		} else if v == NEWLINE {
			x = 0
			y = y + 1
		} else if v == UP || v == DOWN || v == LEFT || v == RIGHT {
			grid[Point{x, y}] = v
			// robot = Robot{x, y, v}
			x = x + 1
		} else {
			panic("unexpected camera output")
		}
	}

	// find intersections
	N := []Point{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
	sum := 0
	for p, v := range grid {
		if v != SCAFFOLD {
			continue
		}

		nbs := []Point{}

		for _, d := range N {
			pn := Point{p.x + d.x, p.y + d.y}
			vn, ok := grid[pn]

			if ok && vn == SCAFFOLD {
				nbs = append(nbs, pn)
			}
		}

		if len(nbs) == 4 {
			grid[p] = INTERSECTION
			sum += p.x * p.y
		}
	}

	display(grid)

	fmt.Println("part1", sum)
}
