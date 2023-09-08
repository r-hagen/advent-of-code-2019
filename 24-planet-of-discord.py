input = [x.strip() for x in open("in").readlines()]

G = {}

for y, row in enumerate(input):
    for x, col in enumerate(row):
        G[x, y] = True if col == "#" else False


def size(A):
    return 1 + max(map(lambda t: t[0], A.keys()))


def display(A, m, d=None):
    print("minute", m, "depth", d)
    s = size(A)
    for y in range(0, s):
        for x in range(0, s):
            print("#" if A[x, y] else ".", end="")
        print("")
    print("")


def bugs(A):
    return sum(A.values())


def part1(G):

    def key(A):
        return tuple([(x, y) for x, y in A if A[x, y]])

    def rating(A):
        s = size(A)
        p = 0
        for y in range(0, s):
            for x in range(0, s):
                if A[x, y]:
                    p += pow(2, y*s+x)
        return p

    N = [(-1, 0), (0, 1), (1, 0), (0, -1)]
    SEEN = set([key(G)])

    m = 0

    while True:
        GN = {}

        for x, y in G.keys():
            bugs = sum([G[x+xn, y+yn] for xn, yn in N if (x+xn, y+yn) in G])
            if G[x, y] and bugs != 1:
                GN[x, y] = False
            elif not G[x, y] and bugs in [1, 2]:
                GN[x, y] = True
            else:
                GN[x, y] = G[x, y]

        G = GN
        m += 1

        k = key(G)
        if k in SEEN:
            break
        else:
            SEEN.add(k)

    return rating(G)


print("part1", part1(G))


def empty():
    G = {}
    for y in range(0, 5):
        for x in range(0, 5):
            G[x, y] = False
    return G


def mutate(v, b):
    if v and b != 1:
        return False
    elif not v and b in [1, 2]:
        return True
    else:
        return v


def part2(G):
    N = [(-1, 0), (0, 1), (1, 0), (0, -1)]

    L = {}
    L[0] = G

    m = 0

    while m < 200:
        LN = {}

        L[min(L.keys()) - 1] = empty()
        L[max(L.keys()) + 1] = empty()

        for d, G in L.items():
            GN = {}

            for x, y in G.keys():
                # center
                if (x, y) == (2, 2):
                    GN[x, y] = False
                    continue
                # outer corners
                elif (x, y) == (0, 0):
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[1, 0], G[0, 1], GP[2, 1], GP[1, 2]])
                    GN[x, y] = mutate(G[x, y], b)
                elif (x, y) == (4, 0):
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[3, 0], G[4, 1], GP[2, 1], GP[3, 2]])
                    GN[x, y] = mutate(G[x, y], b)
                elif (x, y) == (4, 4):
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[4, 3], G[3, 4], GP[3, 2], GP[2, 3]])
                    GN[x, y] = mutate(G[x, y], b)
                elif (x, y) == (0, 4):
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[0, 3], G[1, 4], GP[2, 3], GP[1, 2]])
                    GN[x, y] = mutate(G[x, y], b)
                # outer edges
                elif x == 0 and 1 <= y <= 3:
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[x+1, y], G[x, y+1], G[x, y-1], GP[1, 2]])
                    GN[x, y] = mutate(G[x, y], b)
                elif x == 4 and 1 <= y <= 3:
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[x-1, y], G[x, y+1], G[x, y-1], GP[3, 2]])
                    GN[x, y] = mutate(G[x, y], b)
                elif y == 0 and 1 <= x <= 3:
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[x-1, y], G[x+1, y], G[x, y+1], GP[2, 1]])
                    GN[x, y] = mutate(G[x, y], b)
                elif y == 4 and 1 <= x <= 3:
                    GP = L[d-1] if d-1 in L else empty()
                    b = sum([G[x-1, y], G[x+1, y], G[x, y-1], GP[2, 3]])
                    GN[x, y] = mutate(G[x, y], b)
                # inner edges
                elif x == 2 and y == 1:
                    GP = L[d+1] if d+1 in L else empty()
                    b = sum([G[x-1, y], G[x+1, y], G[x, y-1]])
                    b += sum([GP[xp, 0] for xp in range(0, 5)])
                    GN[x, y] = mutate(G[x, y], b)
                elif x == 3 and y == 2:
                    GP = L[d+1] if d+1 in L else empty()
                    b = sum([G[x, y-1], G[x, y+1], G[x+1, y]])
                    b += sum([GP[4, yp] for yp in range(0, 5)])
                    GN[x, y] = mutate(G[x, y], b)
                elif x == 2 and y == 3:
                    GP = L[d+1] if d+1 in L else empty()
                    b = sum([G[x-1, y], G[x+1, y], G[x, y+1]])
                    b += sum([GP[xp, 4] for xp in range(0, 5)])
                    GN[x, y] = mutate(G[x, y], b)
                elif x == 1 and y == 2:
                    GP = L[d+1] if d+1 in L else empty()
                    b = sum([G[x, y-1], G[x, y+1], G[x-1, y]])
                    b += sum([GP[0, yp] for yp in range(0, 5)])
                    GN[x, y] = mutate(G[x, y], b)
                else:
                    b = sum([G[x+xn, y+yn]
                            for xn, yn in N if (x+xn, y+yn) in G])
                    GN[x, y] = mutate(G[x, y], b)

            LN[d] = GN

        L = LN
        m += 1

    return sum(map(lambda A: bugs(A), L.values()))


print("part2", part2(G))
