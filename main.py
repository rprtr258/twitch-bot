import urllib
import subprocess
import os
from random import choices
import sys

from flask import Flask, request

import lis


sys.stdout.reconfigure(encoding="utf-8")
app = Flask(__name__)

@app.route("/")
def latex():
    if request.args.get('f'):
        formula = request.args.get('f').replace(' ', '+')
        print(repr(formula))
        res = "https://latex.codecogs.com/gif.latex?%5Cdpi%7B300%7D%20%5Chuge%20" + urllib.request.pathname2url(formula.replace(' ', "&space;")).replace('/', '%5C')
        print(repr(res))
        return res
    elif request.args.get('p'):
        pyth_prog = request.args.get('p').replace(' ', '+').replace('_', ' ')
        print(repr(pyth_prog))
        with open("/app/c", "w") as fd:
            fd.write(pyth_prog)
        ps = subprocess.run(['echo', '42'], check=True, capture_output=True)
        processNames = subprocess.run(["python", "pyth/pyth.py", "c"], input=ps.stdout, capture_output=True)
        out = processNames.stdout.decode('utf-8')
        out = out.replace('\r', '').replace('\n', ' ')
        for w in ["ban", "disconnect", "timeout", "unban", "slow", "slowoff", "followers", "followersoff", "subscribers", "subscribersoff", "clear", "uniquechat", "uniquechatoff", "emoteonly", "emoteonlyoff", "commercial", "host", "unhost", "raid", "unraid", "marker"]:
            if f"/{w}" in out or f".{w}" in out:
                return "пососи"
        print(repr(out))
        return out
    elif request.args.get('c'):
        lisp_prog = request.args.get('c').replace(' ', '+').replace('_', ' ')
        print(repr(lisp_prog))
        res = lis.eval(lis.parse(lisp_prog))
        print(repr(res))
        return lis.lispstr(res)

def levenshteinDistance(s1, s2):
    if abs(len(s1) - len(s2)) > 20:
        return 9999999999
    if len(s1) > len(s2):
        s1, s2 = s2, s1
    distances = range(len(s1) + 1)
    for i2, c2 in enumerate(s2):
        distances_ = [i2+1]
        for i1, c1 in enumerate(s1):
            if c1 == c2:
                distances_.append(distances[i1])
            else:
                distances_.append(1 + min((distances[i1], distances[i1 + 1], distances_[-1])))
        distances = distances_
    return distances[-1]

def get_begins(message_word, vocabulary):
    words = {}
    for prior_ngram in vocabulary:
        dist = levenshteinDistance(prior_ngram, message_word)
        if dist < len(message_word) * 1.2:
            words[prior_ngram] = dist
    return list(sorted(words.items(), key=lambda kv:kv[1]))[:10]

@app.route("/g")
def genn():
    import json
    import random
    import math
    from quotes_generator.ngram import NGram
    model = NGram(3)
    model.load("quotes_generator/model.json")
    if request.args.get('m'):
        message = request.args.get('m').replace(' ', '+').replace('_', ' ').split()
        if len(message) > 1:
            message = list(zip(message, message[1:]))
        message = sorted(message, key=len)[-5:]
        with open("quotes_generator/model.json", "r") as fd:
            vocabulary = json.load(fd)["prior"].keys()
        begins = sum([get_begins(w, vocabulary) for w in message], [])
        print(begins, flush=True)
        begin = random.choices([w for w, _ in begins], [math.pow(2, -p) for _, p in begins])
        return model.generate(begin)
    return model.generate()

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)
