import urllib
import subprocess
import os
from random import choices
import sys
import json

from flask import Flask, request

import lis


sys.stdout.reconfigure(encoding="utf-8")
app = Flask(__name__)
with open("quotes_generator/model.json", "r") as fd:
    vocabulary = json.load(fd)["prior"].keys()

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

cache = {}
def levenshteinDistance(s1, s2):
    import edlib
    if (s1, s2) in cache:
        return cache[(s1, s2)]
    cache[(s1, s2)] = edlib.align(s1, s2)["editDistance"]
    return cache[(s1, s2)]

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
        begins = sum([get_begins(w, vocabulary) for w in message], [])
        print(begins, flush=True)
        begin = random.choices([w for w, _ in begins], [math.pow(2, -p) for _, p in begins])
        return model.generate(begin)
    return model.generate()

def reset_chess_save():
    with open("history", "w") as fd:
        fd.write("")

@app.route("/h/reset")
def chess_reset():
    reset_chess_save()
    return "Chess were reset"

@app.route("/h/history")
def chess_history():
    with open("history", "r") as fd:
        return "History: " + ", ".join(map(lambda s: s.strip(), fd))

@app.route("/h/position")
def chess_position():
    from sunfish.sunfish import Position, initial, print_pos, MATE_LOWER, parse, MATE_UPPER, render, print_pos
    import re
    hist = [Position(initial, 0, (True,True), (True,True), 0, 0)]
    raw_hist = []
    with open("history", "r") as fd:
        i = 0
        for line in fd:
            if line.strip() == "":
                continue
            match = re.match('([a-h][1-8])'*2, line.strip())
            old_move = parse(match.group(1)), parse(match.group(2))
            if i % 2 == 1:
                old_move = (119-old_move[0], 119-old_move[1])
            hist.append(hist[-1].move(old_move))
            i += 1
    return print_pos(hist[-1])

def render_move(move):
    from sunfish.sunfish import render
    return render(move[0]) + render(move[1])

@app.route("/h")
def chess():
    from sunfish.sunfish import Position, Searcher, initial, print_pos, MATE_LOWER, parse, MATE_UPPER
    import re
    import time
    _input = request.args.get('m')
    match = re.match('([a-h][1-8])'*2, _input)
    if match:
        move = parse(match.group(1)), parse(match.group(2))
    else:
        return "Please enter a move like g8f6"
    hist = [Position(initial, 0, (True,True), (True,True), 0, 0)]
    searcher = Searcher()

    with open("history", "r") as fd:
        i = 0
        for line in fd:
            if line.strip() == "":
                continue
            match = re.match('([a-h][1-8])'*2, line.strip())
            old_move = parse(match.group(1)), parse(match.group(2))
            if i % 2 == 1:
                old_move = (119-old_move[0], 119-old_move[1])
            hist.append(hist[-1].move(old_move))
            i += 1
    if move not in hist[-1].gen_moves():
        moves = ", ".join(map(render_move, hist[-1].gen_moves()))
        return "Correct moves are: " + moves
    with open("history", "a") as fd:
        fd.write(_input + "\n")
    hist.append(hist[-1].move(move))
    if hist[-1].score <= -MATE_LOWER:
        reset_chess_save()
        return "You won"

    start = time.time()
    for _depth, move, score in searcher.search(hist[-1], hist):
        if time.time() - start > 1:
            break
    the_move = render_move((119-move[0],119-move[1]))
    hist.append(hist[-1].move(move))
    if score == MATE_UPPER:
        reset_chess_save()
        return "Checkmate! " + the_move
    else:
        with open("history", "a") as fd:
            fd.write(the_move + "\n")
    if hist[-1].score <= -MATE_LOWER:
        reset_chess_save()
        return "You lost"
    return "My move:" + the_move

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)
