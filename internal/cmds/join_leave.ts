import {type Command} from "../cmds.ts";
import * as message from "../message.ts";
import * as permissions from "../permissions.ts";
import * as services from "../services.ts";

export const JoinCmd: Command = {
  RequiredPermissions: ["global_admin"],
  Command: "?join",
  Description: "Join channel",
  Run(s: services.Services, perms: permissions.PermissionsList, msg: message.TwitchMessage): string {
    const words = msg.Text.split(" ");
    if (words.length != 2) {
      return `Usage: ${this.Command} <channel>`;
    }

    const channel = words[1];

    s.ChatClient.Join(channel)

    s.Insert("joined_channels", {
      "channel": channel,
    })

    return `joined ${channel} successfully`;
  },
};

export const LeaveCmd: Command = {
  RequiredPermissions: ["global_admin"],
  Command: "?leave",
  Description: "Leave channel",
  Run(s: services.Services, perms: permissions.PermissionsList, msg: message.TwitchMessage): string {
    const words = msg.Text.split(" ");
    if (words.length != 2) {
      return `Usage: ${this.Command} <channel>`;
    }

    const channel = words[1];

    s.ChatClient.Depart(channel)

    const db = s.Backend.DB()
    db.Delete("joined_channels", `channel=${channel}`).Execute()

    return `left ${channel} successfully`;
  },
};
