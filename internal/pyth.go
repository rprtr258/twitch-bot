package internal

// @app.route("/p")
// def pyth():
//     if not request.args.get('p'):
//         return "provide program with 'p' request parameter"
//     pyth_prog = request.args.get('p').replace(' ', '+').replace('_', ' ')
//     print(repr(pyth_prog))
//     with open("/app/c", "w") as fd:
//         fd.write(pyth_prog)
//     ps = run(['echo', '42'], check=True, capture_output=True)
//     processNames = run(["python", "pyth/pyth.py", "c"], input=ps.stdout, capture_output=True)
//     out = processNames.stdout.decode('utf-8')
//     out = out.replace('\r', '').replace('\n', ' ')
//     for w in ["ban", "disconnect", "timeout", "unban", "slow", "slowoff", "followers", "followersoff", "subscribers", "subscribersoff", "clear", "uniquechat", "uniquechatoff", "emoteonly", "emoteonlyoff", "commercial", "host", "unhost", "raid", "unraid", "marker"]:
//         if f"/{w}" in out or f".{w}" in out:
//             return "пососи"
//     print(repr(out))
//     return out
