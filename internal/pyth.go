package internal

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

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

const (
	pythCmd = "?pyth"
)

var running atomic.Int64

func (s *Services) pyth(message twitch.PrivateMessage) (string, error) {
	if !strings.ContainsRune(message.Message, ' ') {
		return "Pyth docs: https://pyth.readthedocs.io/en/latest/getting-started.html", nil
	}

	program := strings.TrimPrefix(message.Message, pythCmd+" ")

	if running.Load() > 2 {
		return "Сервер загружен, попробуйте позже", nil
	}

	running.Add(1)
	defer func() {
		running.Add(-1)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--memory=500m", "--cpus=1", "--rm", "pyth", program)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if err != nil {
		if err.Error() == "signal: killed" {
			return "Was running too long, killed", nil
		}
		log.Println("pyth stopped", program, err.Error())

		stderrStr, err := io.ReadAll(&stderr)
		if err != nil {
			return "", err
		}
		return string(stderrStr), nil
	}

	return "> " + string(stdout), nil
}
