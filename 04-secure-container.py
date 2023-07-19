input = "130254-678275"
lower, upper = tuple(map(int, input.split("-")))

ans1 = 0
ans2 = 0

for password in range(lower, upper):
    adjacent = False
    adjacent2 = False
    decreasing = False

    digits = str(password)

    for i in range(0, len(digits)-1):
        if digits[i] == digits[i+1]:
            adjacent = True

            nb = [digits[i-1:i], digits[i+2:i+3]]
            if all(map(lambda x: x != digits[i], nb)):
                adjacent2 = True

        if digits[i] > digits[i+1]:
            decreasing = True

    if adjacent and not decreasing:
        ans1 += 1
    if adjacent2 and not decreasing:
        ans2 += 1

print(ans1, ans2)
