input = [x.strip() for x in open("in").readlines()]

deck = [x for x in range(0, 10007)]


def deal_into_stack(deck):
    stack = deck.copy()
    stack.reverse()
    return stack


def cut_n_cards(deck, n):
    stack = deck[n:] + deck[0:n]
    return stack


def deal_with_increment_n(deck, n):
    stack = [0] * len(deck)
    top = 0
    idx = 0
    while top < len(deck):
        stack[idx] = deck[top]
        top += 1
        idx = (idx+n) % len(deck)
    return stack


for instruction in input:
    words = instruction.split(" ")
    if words[0] == "cut":
        deck = cut_n_cards(deck, int(words[1]))
    elif words[0] == "deal" and words[1] == "with":
        deck = deal_with_increment_n(deck, int(words[3]))
    elif words[0] == "deal" and words[1] == "into":
        deck = deal_into_stack(deck)


print("part1", deck.index(2019))
