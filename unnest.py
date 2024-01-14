# this is a script to fix the formatting error that
# occurs by running rcl.py (entries with Complementary
# verses have to much nesting)

import json

with open('plans.json', 'r', encoding='utf-8') as file:
    data = json.load(file)


for day in data["RCL"]["days"]:
    lw = type(day[0][0]) is list
    if lw:
        day[0] = day[0][0]


out_file = open("unnested.json", "w")
json.dump(data, out_file, indent=4)
out_file.close()
