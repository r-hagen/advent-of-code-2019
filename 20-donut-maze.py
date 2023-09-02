import heapq

lines = [x.replace('\n', '') for x in open('in').readlines()]

oy = [0, 121]
iy = [35, 86]

ox = [0, 127]
ix = [35, 92]

grid = {}
pairs = {}
portals = {}
outer = {}
inner = {}

for y, line in enumerate(lines):
    for x, column in enumerate(line):
        if column == ' ':
            continue
        elif column.isupper():
            if x == ox[0]:
                pairs[y, x+2] = lines[y][x] + lines[y][x+1]
                outer[y, x+2] = True
            elif x == ox[1]:
                pairs[y, x-1] = lines[y][x] + lines[y][x+1]
                outer[y, x-1] = True
            elif x == ix[0]:
                pairs[y, x-1] = lines[y][x] + lines[y][x+1]
                inner[y, x-1] = True
            elif x == ix[1]:
                pairs[y, x+2] = lines[y][x] + lines[y][x+1]
                inner[y, x+2] = True
            elif y == oy[0]:
                pairs[y+2, x] = lines[y][x] + lines[y+1][x]
                outer[y+2, x] = True
            elif y == oy[1]:
                pairs[y-1, x] = lines[y][x] + lines[y+1][x]
                outer[y-1, x] = True
            elif y == iy[0]:
                pairs[y-1, x] = lines[y][x] + lines[y+1][x]
                inner[y-1, x] = True
            elif y == iy[1]:
                pairs[y+2, x] = lines[y][x] + lines[y+1][x]
                inner[y+2, x] = True
        elif column == '.':
            grid[y, x] = True
        elif column == '#':
            grid[y, x] = False

for k1, v1 in pairs.items():
    if v1 == 'AA':
        start = k1
    elif v1 == 'ZZ':
        goal = k1
    else:
        pair = [k2 for k2, v2 in pairs.items() if v2 == v1]
        portals[pair[0]] = pair[1]
        portals[pair[1]] = pair[0]


def part1():
    v = set()
    q = []
    heapq.heappush(q, (0, start))

    while len(q):
        steps, pos = heapq.heappop(q)

        if pos == goal:
            return steps

        if pos in v:
            continue

        v.add(pos)

        neighbors = [(-1, 0), (0, -1), (1, 0), (0, 1)]

        for dx, dy in neighbors:
            n = (pos[0]+dy, pos[1]+dx)
            if n in grid and grid[n] is True:
                heapq.heappush(q, (steps+1, n))

        if pos in portals:
            heapq.heappush(q, (steps+1, portals[pos]))


print("part1", part1())


def part2():
    v = set()
    q = []
    heapq.heappush(q, ((0, 0), start))

    while len(q):
        (steps, level), pos = heapq.heappop(q)

        if level == 0 and pos == goal:
            return steps

        if (pos, level) in v:
            continue

        v.add((pos, level))

        neighbors = [(-1, 0), (0, -1), (1, 0), (0, 1)]

        for dx, dy in neighbors:
            n = (pos[0]+dy, pos[1]+dx)
            if n in grid and grid[n] is True:
                heapq.heappush(q, ((steps+1, level), n))

        if pos in portals:
            if pos in inner:
                heapq.heappush(q, ((steps+1, level+1), portals[pos]))
            elif pos in outer and level != 0:
                heapq.heappush(q, ((steps+1, level-1), portals[pos]))


print("part2", part2())
