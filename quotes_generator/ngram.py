import random
import re
import json

def normalize_line(line):
    res = line.strip()
    return res

class NGram():
    def __init__(self, n):
        assert(n >= 1)
        self.n = n
        self.prior = {}
        self.prob = {}
    
    def save(self, filename):
        with open(filename, "w", encoding="utf-8") as f:
            f.write(json.dumps({
                "prior": self.prior,
                "prob": self.prob
            }))

    def load(self, filename):
        with open(filename, "r", encoding="utf-8") as f:
            data = json.loads(f.read())
            self.prior = data["prior"]
            self.prob = data["prob"]
    
    def train(self, dataset_filename):
        transitions = {}
        count = {}
        with open(dataset_filename, "r", encoding="utf-8") as f:
            for line in f:
                line = normalize_line(line)
                words = line.split()
                B = " ".join(words[:self.n - 1])
                if not B in self.prior:
                    self.prior[B] = 0
                self.prior[B] += 1
                for i in range(0, len(words) - self.n + 1):
                    B = " ".join(words[i : i + self.n - 1])
                    A = words[i + self.n - 1]
                    if not B in transitions:
                        transitions[B] = []
                    if not f"{B} {A}" in count:
                        count[f"{B} {A}"] = 0
                    if not B in count:
                        count[B] = 0
                    count[f"{B} {A}"] += 1
                    count[B] += 1
                    if not A in transitions[B]:
                        transitions[B].append(A)
        for B in transitions.keys():
            self.prob[B] = {}
            for A in transitions[B]:
                self.prob[B][A] = count[f"{B} {A}"] / count[B]

    def generate(self, begin=None, length=100):
        if begin is None:
            begin = random.choices(list(self.prior.keys()), weights=self.prior.values())
            begin = begin[0].split()
        else:
            begin = begin.split()
        res = begin
        for _ in range(length):
            B = " ".join(res[-self.n + 1:])
            if B not in self.prob:
                break
            distribution = self.prob[B].values()
            A = random.choices(list(self.prob[B].keys()), weights=distribution)
            res += A
        return " ".join(res)

