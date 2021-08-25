from collections import Counter

with open("dict", "r", encoding="utf-8") as fd:
    cnt = sorted(set(filter(lambda x: all(c in "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" for c in x), fd.read().split())))
for w in cnt:
    print(w)