import {type Command} from "./cmds.ts";
import * as message from "./message.ts";
import * as permissions from "./permissions.ts";
import * as services from "./services.ts";

// TODO: list only available commands
const CommandsCmd: Command = {
  RequiredPermissions(): string[] {
    return ["execute_commands"];
  },
  Command(): string {
    return "?commands"
  },
  Description(): string {
    return "List all commands"
  },
  Run(_s: services.Services, _p: permissions.PermissionsList, _m: message.TwitchMessage): string {
    return allCommands.map(cmd => `${cmd.Command()} - ${cmd.Description()}`).join(", ");
  },
}

// TODO: permissions constants/store in db to dinamically update
// TODO: change to map
const allCommands: Command[] = [
  cmds.IntelCmd,
  cmds.JoinCmd,
  cmds.LeaveCmd,
  cmds.FedCmd,
  cmds.BlabGenCmd,
  cmds.BlabContinueCmd,
  cmds.BlabReadCmd,
  CommandsCmd,
  cmds.PermsCmd,
  cmds.PastaSearchCmd,
];

export function OnPrivateMessage(s: services.Services, msg: twitch.PrivateMessage): void {
  const msgMy: message.TwitchMessage = {
    Text:    msg.Message,
    User:    msg.User,
    At:      msg.Time,
    Channel: msg.Channel,
  };
  s.LogMessage(msgMy);
  const userPermissions = permissions.GetPermissions(s.Permissions, {
    "username": msg.User.Name, // TODO: replace with user id
    "channel":  msg.Channel,
    // TODO: vips, moders
  });

  const firstWord = msg.Message.split(" ")[0];
  const userName = msg.User.Name;

  for (const cmd of allCommands
    .filter(cmd => cmd.Command() === firstWord)
    // TODO: add ban permissions
    .filter(cmd => permissions.Has(userPermissions, ...cmd.RequiredPermissions()))
  ) {
    const whisper = !permissions.Has(userPermissions, "say_response");

    const response = (() => {try {
      const res = cmd.Run(s, userPermissions, msgMy);
      return res === "" ? "Empty response" : res;
    } catch (err) {
      return `Internal error: ${err}`;
    }})().replace(/\n/g, " ");

    // TODO: move responding to commands
    // TODO: rewrite
    // TODO: fix response: Must be less than 500 character(s). somewhere
    // TODO: change to permission check, to send message by parts
    const croppedResponse = cropMessage(response)
    if (whisper) {
      s.ChatClient.Whisper(msg.User.Name, croppedResponse)
    } else {
      s.ChatClient.Reply(msg.Channel, msg.ID, croppedResponse)
    }

    // TODO: fix not logging blab cmds
    s.Insert("chat_commands", {
      "command":  cmd.Command(),
      "args":     msg.Message,
      "response": croppedResponse,
      "user":     userName,
      "channel":  msg.Channel,
    })
  }
}

export function OnWhisperMessage(s: services.Services, msg: twitch.WhisperMessage) {
  const msgMy: message.TwitchMessage = {
    Text:    msg.Message,
    User:    msg.User,
    At:      new Date(),
    Channel: msg.User.Name,
  };
  s.LogMessage(msgMy)
  const userPermissions = s.Permissions.GetPermissions({
    "username": msg.User.Name, // TODO: replace with user id
    "channel":  msgMy.Channel,
  });

  const firstWord = msg.Message.split(" ")[0];
  const userName = msg.User.Name;

  for (const cmd of allCommands
    .filter(cmd => cmd.Command() === firstWord)
    // TODO: add ban permissions
    .filter(cmd => permissions.Has(userPermissions, ...cmd.RequiredPermissions()))
  ) {
    const response = (() => {try {
      const res = cmd.Run(s, userPermissions, msgMy);
      return res === "" ? "Empty response" : res;
    } catch (err) {
      return `Internal error: ${err}`;
    }})().replace(/\n/g, " ");

    // TODO: move responding to commands
    // TODO: rewrite
    // TODO: fix response: Must be less than 500 character(s). somewhere
    // TODO: change to permission check, to send message by parts
    const croppedResponse = cropMessage(response);
    s.ChatClient.Whisper(msg.User.Name, croppedResponse);

    // TODO: fix not logging blab cmds
    s.Insert("chat_commands", {
      "command":  cmd.Command(),
      "args":     msg.Message,
      "response": response,
      "user":     userName,
      "channel":  msgMy.Channel,
    });
  }
}

// TODO: fix cutting
function cropMessage(response: string): string {
  return response.substring(0, message.MaxMessageLength);
}
