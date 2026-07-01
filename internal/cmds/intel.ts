import {type Command} from "../cmds.ts";
import * as message from "../message.ts";
import * as permissions from "../permissions.ts";
import * as services from "../services.ts";

function formatDuration(d: Date): string {
  const now = new Date();

  let parts: string[] = [];
  {
    const years = d.getFullYear() - now.getFullYear();
    if (years > 0)
      parts.push(`${years}y`);
    d.setFullYear(now.getFullYear());
  }

  {
    const months = d.getMonth() - now.getMonth();
    if (months > 0)
      parts.push(`${months}mo`);
    d.setMonth(now.getMonth());
  }

  {
    const days = d.getDate() - now.getDate();
    if (days > 0)
      parts.push(`${days}d`);
    d.setDate(now.getDate());
  }

  const minutes = now.getMinutes() - d.getMinutes();
  if (Math.floor(minutes/60) > 0) {
    parts.push(`${Math.floor(minutes/60)}h`);
  }

  parts.push(`${minutes%60}m`);

  return parts.join("");
}

export default {
  RequiredPermissions: [],
  Command: "?intel",
  Description: "Gather intel on user",
  Run(s: services.Services, perms: permissions.PermissionsList, msg: message.TwitchMessage): string {
    const words = msg.Text.split(" ");
    if (words.length < 2) {
      return "No username provided"
    }

    const login = words[1];
    const resp = s.TwitchApiClient.GetUsers({Logins: [login]});
    if (resp.ErrorMessage !== "") {
      throw new Error(resp.ErrorMessage);
    }

    const users = resp.Data.Users;
    if (users.length === 0) {
      return `${login} not found`;
    }

    const user = users[0];
    const createdAt = user.CreatedAt.Time;
    const delta = formatDuration(createdAt);
    return [
      `created=${createdAt.Format("15:04.2.1.2006")}`,
      `ago=${delta}`,
      `broadcaster_type=${user.BroadcasterType}`,
      `id=${user.ID}`,
      `login=${user.Login}`,
      `display_name=${user.DisplayName}`,
      `view_count=${user.ViewCount}`,
    ].join(" ");
  },
} as Command;
