# use to parse reading plan from this source:
# https://lectionary.library.vanderbilt.edu/daily.php
# you should manually copy all relevant plain text from
# 3 pages (for years A, B and C) and put that in rcl.txt file
# after running it you will have a valid JSON array of reading days
# in the same file (rcl.txt)
# NOTE: this script will deal with MOST of the reading days
# in the plan, but will not deal with any of the special days
# (such as Easter, Christmas Day etc) - those are to be filled
# by manually looking up those from this link:
# https://lectionary.library.vanderbilt.edu/texts.php?id=1
# navigate through all the liturgical seasons displayed on the left panel,
# then copypaste and format

with open('rcl.txt', 'r', encoding='utf-8') as file:
    data = file.readlines()

i = 0
for line in data:
    lr = line.split(": ")
    if len(lr) < 2:
        continue
    r = lr[1]
    coll = r.split("; ")
    for j in range(len(coll)):
        w = coll[j]
        w = w.strip("\n")
        w = "\"{}\"".format(w)
        coll[j] = w
    spis = ", ".join(coll)
    spis = "[{}]".format(spis)
    if "Complementary" in lr[0]:
        i -= 1
        spis = "[{}{}],\n".format(data[i].strip("\n"), spis)
    else:
        spis = "[{}],\n".format(spis)
    data[i] = spis
    i += 1

with open('rcl.txt', 'w', encoding='utf-8') as file:
    file.writelines(data)
