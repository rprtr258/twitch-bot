import urllib
import subprocess
import os
from random import choices

from flask import Flask, request

import lis
from vptree import VPTree, editIgnoreCaseDistance


app = Flask(__name__)
with open("quotes_generator/dict2", "r", encoding="utf-8") as fin:
    words = sorted(set(fin.read().split()))
print("building tree...")
tree = VPTree(words, editIgnoreCaseDistance)
print("tree built!")

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
        if "/timeout" in out:
            return "пососи"
        if "/ban" in out:
            return "охуел?"
        print(repr(out))
        return out
    elif request.args.get('c'):
        lisp_prog = request.args.get('c').replace(' ', '+').replace('_', ' ')
        print(repr(lisp_prog))
        res = lis.eval(lis.parse(lisp_prog))
        print(repr(res))
        return lis.lispstr(res)
 
@app.route("/g")
def genn():
    from quotes_generator.ngram import NGram
    model = NGram(3)
    model.load("quotes_generator/model.json")
    word = request.args.get('w')
    if word != "" and word is not None:
        sims = [(res, tree.dist(word, res)) for res in tree.findSimilar(word)]
        begin=choices([w for w, _ in sims], weights=[1 / (2 ** p) for _, p in sims])
    else:
        begin = None
    return model.generate(begin)

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)
