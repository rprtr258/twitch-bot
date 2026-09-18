export const MaxMessageLength = 490;

export type TwitchMessage = {
  Text:    string,
  User:    twitch.User,
  At:      Date,
  Channel: string,
};
