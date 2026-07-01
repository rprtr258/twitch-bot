import * as message from "./message.ts";
import * as permissions from "./permissions.ts";
import * as services from "./services.ts";

export type Command = {
  // minimal permissions required to use the command
  RequiredPermissions: string[],
  // get command name
  Command: string,
  // get command description
  Description: string,
  // execute command
  Run(_services: services.Services, _perms: permissions.PermissionsList, _msg: message.TwitchMessage): string,
};
