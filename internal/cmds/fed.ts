import {type Command} from "../cmds.ts";
import * as message from "../message.ts";
import * as permissions from "../permissions.ts";
import * as services from "../services.ts";

export default {
  RequiredPermissions: ["execute_commands"],
  Command: "?fed",
  Description: "Show how many times word has been used",
  Run(s: services.Services, perms: permissions.PermissionsList, msg: message.TwitchMessage): string {
    const words = msg.Text.split(" ");
    if (words.length !== 2 && words.length !== 3) {
      return `Usage: "${this.Command} <word>" or "${this.Command} <user> <word>" or "${this.Command} * <word>"`;
    }

    const db = s.Backend.DB();

    if (words.length === 3 && words[1] === "*") {
      const theWord = words[2];
      const count: number = db.
        Select("COUNT(*) AS count").
        From("messages").
        Where(dbx.And(
          dbx.NewExp("channel={:channel}", {"channel": msg.Channel}),
          // TODO: fix theWord = '%%' escaping
          dbx.Like("message", theWord),
          dbx.NotLike("message", this.Command).Match(false, true),
        )).
        Row();
      return `Слово ${theWord} было использовано ${count} раз.`;
    }

    const [
      theWord,
      user,
      mention,
    ] = ((): [theWord: string, user: string, mention: string] => {
      if (words.length === 2) {
        return [words[1], msg.User.Name, "Ты"];
      } else {
        const user = words[1];
        return [words[2], user, user];
      }
    })();

    const count: number = db.Select("COUNT(*) AS count").From("messages").Where(dbx.And(
      dbx.NewExp("channel={:channel}", {"channel": msg.Channel}),
      // TODO: fix theWord = '%%' escaping
      dbx.Like("message", theWord),
      dbx.NotLike("message", this.Command).Match(false, true),
      dbx.NewExp("user_name={:user_name}", {"user_name": user}),
    )).Row()

    return `${mention} использовал слово ${theWord} ${count} раз.`;
  },
} as Command;
