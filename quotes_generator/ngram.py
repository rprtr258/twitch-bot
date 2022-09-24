from random import choices
import re
import json
from collections import defaultdict


class NGram():
    def __init__(self, n):
        assert n >= 1, "size of n-gram must be at leat 1"
        self.n = n
        self.prior = defaultdict(int)
        self.prob = {}
    
    def save(self, model_file):
        with open(model_file, "w", encoding="utf-8") as fd:
            fd.write(json.dumps({
                "prior": self.prior,
                "prob": self.prob
            }))

    def load(self, model_file):
        with open(model_file, "r", encoding="utf-8") as fd:
            data = json.loads(fd.read())
            self.prior = data["prior"]
            self.prob = data["prob"]
    
    def train(self, dataset_file):
        from tqdm import tqdm
        transitions = defaultdict(set)
        count = defaultdict(int)
        with open(dataset_file, "r", encoding="utf-8") as fd:
            for line in tqdm(fd):
                words = line.strip().split()
                B = " ".join(words[:self.n - 1])
                self.prior[B] += 1
                for i in tqdm(range(0, len(words) - self.n + 1)):
                    B = " ".join(words[i : i + self.n - 1])
                    self.prior[B] += 1
                    A = words[i + self.n - 1]
                    count[f"{B} {A}"] += 1
                    count[B] += 1
                    transitions[B].add(A)
        self.prob = {
            B: {
                A: count[f"{B} {A}"] / count[B]
                for A in transitions[B]
            }
            for B in transitions.keys()
        }

    def generate(self, begin=None, length=100):
        def randchoice(distribution):
            return choices(
                list(distribution.keys()),
                weights=distribution.values()
            )
        if begin is None:
            begin = randchoice(self.prior)
        begin = begin[0].split()
        res = begin
        for _ in range(length):
            B = " ".join(res[-self.n + 1:])
            if B not in self.prob:
                break
            res += randchoice(self.prob[B])
        return " ".join(res)

