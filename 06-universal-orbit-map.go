package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("in")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	orbit_map := map[string][]string{}

	for scanner.Scan() {
		line := scanner.Text()

		objects := strings.Split(line, ")")
		object1 := objects[0]
		object2 := objects[1]

		orbit_map[object1] = append(orbit_map[object1], object2)
	}

	ans1 := count_orbits(orbit_map, 0, "COM")
	fmt.Println("ans1", ans1)

	p1 := find_path(orbit_map, "COM", "YOU", make([]string, 0))
	p2 := find_path(orbit_map, "COM", "SAN", make([]string, 0))
	ans2 := orbital_transfers(p1, p2)
	fmt.Println("ans2", ans2)
}

func count_orbits(orbit_map map[string][]string, depth int, object string) int {
	orbit_count := depth

	for _, o := range orbit_map[object] {
		orbit_count += count_orbits(orbit_map, depth+1, o)
	}

	return orbit_count
}

func find_path(orbit_map map[string][]string, origin string, target string, path []string) []string {
	for _, object := range orbit_map[origin] {
		cpath := make([]string, len(path))
		copy(cpath, path)
		cpath = append(cpath, object)

		if object == target {
			return cpath
		}

		cpath = find_path(orbit_map, object, target, cpath)
		if len(cpath) != 0 {
			return cpath
		}
	}

	return make([]string, 0)
}

func orbital_transfers(path1 []string, path2 []string) int {
	for i := len(path1) - 1; i >= 0; i-- {
		for j := len(path2) - 1; j >= 0; j-- {
			if path1[i] == path2[j] {
				return len(path1) - i - 2 + len(path2) - j - 2
			}
		}
	}
	panic("orbital transfer is impossible")
}
