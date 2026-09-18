import {type Command} from "../cmds.ts";
import * as message from "../message.ts";
import * as permissions from "../permissions.ts";
import * as services from "../services.ts";

export default {
  RequiredPermissions: [],
  Command: "?perms",
  Description: "Get permissions",
  Run(s: services.Services, perms: permissions.PermissionsList, msg: message.TwitchMessage): string {
    const words = msg.Text.split(" ");
    let permissionsToShow = perms;
    if (words.length > 1) {
      const claims: permissions.Claims = {};
      for (const word of words.slice(1)) {
        const idx = word.indexOf(':')
        if (idx === -1 || idx === 0 || idx === word.length-1) {
          return `Usage: "${this.Command}" or "${this.Command} key:value key:value"`;
        }
        const key = word.substring(0, idx);
        const value = word.substring(idx+1);
        claims[key] = value
      }
      permissionsToShow = permissions.GetPermissions(s.Permissions, claims);
    }
    return permissionsToShow.entries().toArray().join(",");
  },
} as Command;
