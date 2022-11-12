package cmds

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"
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

type PythCmd struct {
	running atomic.Int64
}

func (*PythCmd) Command() string {
	return "?pyth"
}

func (*PythCmd) Description() string {
	return "Eval in pyth"
}

func (cmd *PythCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, msg message.TwitchMessage) (string, error) {
	if !strings.ContainsRune(msg.Text, ' ') {
		return "Pyth docs: https://pyth.readthedocs.io/en/latest/getting-started.html", nil
	}

	program := strings.TrimPrefix(msg.Text, cmd.Command()+" ")

	if cmd.running.Load() > 2 {
		return "Сервер загружен, попробуйте позже", nil
	}

	cmd.running.Add(1)
	defer func() {
		cmd.running.Add(-1)
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// TODO: secure
	// proc := exec.CommandContext(ctx, "docker", "run", "--memory=500m", "--cpus=1", "--rm", "pyth", program)
	proc := exec.CommandContext(ctx, "python3", "pyth/pyth.py", "-c", program)
	var stderr bytes.Buffer
	proc.Stderr = &stderr
	stdout, err := proc.Output()
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

	stdoutStr := string(stdout)
	if strings.HasPrefix(stdoutStr, "/") || strings.HasPrefix(stdoutStr, ".") {
		return "пососи", nil
	}

	return stdoutStr, nil
}
