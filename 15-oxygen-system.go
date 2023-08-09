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

func (l Location) Move(d Location) Location {
	return Location{l.x + d.x, l.y + d.y}
}

type Neighbor struct {
	location  Location
	move      int64
	backtrack int64
}

func (l Location) Neighbors() []Neighbor {
	north := Neighbor{Location{l.x, l.y + 1}, NORTH, SOUTH}
	east := Neighbor{Location{l.x + 1, l.y}, EAST, WEST}
	south := Neighbor{Location{l.x, l.y - 1}, SOUTH, NORTH}
	west := Neighbor{Location{l.x - 1, l.y}, WEST, EAST}

	return []Neighbor{north, east, south, west}
}

const (
	WALL   int64 = 0
	OPEN   int64 = 1
	OXYGEN int64 = 2

	NORTH int64 = 1
	SOUTH int64 = 2
	WEST  int64 = 3
	EAST  int64 = 4
)

type Droid struct {
	grid     Grid
	computer Computer
	position Location
	ans1     int
}

func MakeDroid(com Computer) *Droid {
	start := Location{0, 0}
	grid := make(Grid)
	grid[start] = OPEN
	return &Droid{grid, com, start, math.MaxInt}
}

func (d *Droid) Part1(depth int, location Location) int {
	for _, neighbor := range location.Neighbors() {
		_, visited := d.grid[neighbor.location]

		if visited {
			continue
		}

		d.computer.Run(Input{neighbor.move})
		status := d.computer.out[0]

		if status == WALL {
			d.grid[neighbor.location] = WALL
		} else if status == OPEN {
			d.grid[neighbor.location] = OPEN
			d.position = neighbor.location

			d.Part1(depth+1, neighbor.location)

			d.computer.Run(Input{neighbor.backtrack})
			d.position = location
		} else if status == OXYGEN {
			d.grid[neighbor.location] = OXYGEN
			d.position = neighbor.location

			if depth+1 < d.ans1 {
				d.ans1 = depth + 1
			}

			d.Plot(true)

			d.computer.Run(Input{neighbor.backtrack})
			d.position = location
		}

	}

	return d.ans1
}

func (d *Droid) Part2(minute int) int {
	oxygen := []Location{}
	for location, status := range d.grid {
		if status == OXYGEN {
			oxygen = append(oxygen, location)
		}
	}

	didSpread := false

	for _, location := range oxygen {
		for _, neighbor := range location.Neighbors() {
			status, visited := d.grid[neighbor.location]

			if !visited {
				panic("unexplored location")
			}

			if status == OPEN {
				d.grid[neighbor.location] = OXYGEN
				didSpread = true
			}
		}
	}

	if didSpread {
		return d.Part2(minute + 1)
	}

	return minute
}

func (d *Droid) Plot(drawDroid bool) {
	// [xmin, xmax, ymin, ymax]
	bounds := make([]int64, 4)

	for location := range d.grid {
		if location.x < bounds[0] {
			bounds[0] = location.x
		}
		if location.x > bounds[1] {
			bounds[1] = location.x
		}
		if location.y < bounds[2] {
			bounds[2] = location.y
		}
		if location.y > bounds[3] {
			bounds[3] = location.y
		}
	}

	fmt.Println()

	for y := bounds[3]; y >= bounds[2]; y-- {
		for x := bounds[0]; x <= bounds[1]; x++ {
			location := Location{x, y}
			status, ok := d.grid[location]

			if location == d.position && drawDroid {
				fmt.Print("D")
			} else if !ok {
				fmt.Print(" ")
			} else if status == WALL {
				fmt.Print("#")
			} else if status == OPEN {
				fmt.Print(".")
			} else if status == OXYGEN {
				fmt.Print("o")
			} else {
				panic("unexpected status")
			}
		}
		fmt.Print("\n")
	}

	fmt.Println()
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

func main() {
	file, _ := os.Open("in")
	defer file.Close()

	reader := bufio.NewReader(file)
	input, _ := reader.ReadString('\n')
	program := parseProgram(input)

	computer := MakeComputer(program, Input{})
	droid := MakeDroid(computer)

	ans1 := droid.Part1(0, droid.position)
	droid.Plot(false)
	fmt.Println("ans1", ans1)

	ans2 := droid.Part2(0)
	droid.Plot(false)
	fmt.Println("ans2", ans2)
}
