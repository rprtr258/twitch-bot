from subprocess import run
import sys; sys.stdout.reconfigure(encoding="utf-8")
from json import load

from flask import Flask, request


app = Flask(__name__)
with open("model.json", "r") as fd:
    vocabulary = load(fd)["prior"].keys()


@app.route("/p")
def pyth():
    if not request.args.get('p'):
        return "provide program with 'p' request parameter"
    pyth_prog = request.args.get('p').replace(' ', '+').replace('_', ' ')
    print(repr(pyth_prog))
    with open("/app/c", "w") as fd:
        fd.write(pyth_prog)
    ps = run(['echo', '42'], check=True, capture_output=True)
    processNames = run(["python", "pyth/pyth.py", "c"], input=ps.stdout, capture_output=True)
    out = processNames.stdout.decode('utf-8')
    out = out.replace('\r', '').replace('\n', ' ')
    for w in ["ban", "disconnect", "timeout", "unban", "slow", "slowoff", "followers", "followersoff", "subscribers", "subscribersoff", "clear", "uniquechat", "uniquechatoff", "emoteonly", "emoteonlyoff", "commercial", "host", "unhost", "raid", "unraid", "marker"]:
        if f"/{w}" in out or f".{w}" in out:
            return "пососи"
    print(repr(out))
    return out

cache = {}
def levenshteinDistance(s1, s2):
    import edlib
    s1, s2 = s1.lower(), s2.lower()
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
    from random import choices
    import math
    from quotes_generator.ngram import NGram
    model = NGram(3)
    model.load("model.json")
    if request.args.get('m'):
        message = request.args.get('m').replace(' ', '+').replace('_', ' ').split()
        if len(message) > 1:
            message = list(zip(message, message[1:]))
        message = sorted(message, key=len)[-5:]
        begins = sum([get_begins(w, vocabulary) for w in message], [])
        print(begins, flush=True)
        begin = choices([w for w, _ in begins], [math.pow(2, -p) for _, p in begins])
        return model.generate(begin)
    return model.generate()

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)
