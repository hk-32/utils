import json

txt = open("basic.json").read()

print(json.loads(txt))