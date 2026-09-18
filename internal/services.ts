import * as message from "./message.ts";
import * as permissions from "./permissions.ts";

export type Services = {
  ChatClient:      twitch.Client,
  TwitchApiClient: helix.Client,
  Backend:         pocketbase.PocketBase,
  Balaboba:        balaboba.Client,
  Permissions:     permissions.Permissions,
};

function Insert(s: Services, collectionName: string, data: Record<string, any>): string {
  const collection = s.Backend.Dao().FindCollectionByNameOrId(collectionName);

  const record = models.NewRecord(collection);
  for (const [k, v] of Object.entries(data)) {
    record.Set(k, v);
  }

  const form = forms.NewRecordUpsert(s.Backend.App, record);
  form.Validate()
  form.Submit()

  return record.Id
}

function LogMessage(s: Services, msg: message.TwitchMessage): void {
  s.Insert("messages", {
    "user_id":           msg.User.ID,
    "message":           msg.Text,
    "channel":           msg.Channel,
    "user_name":         msg.User.Name,
    "user_display_name": msg.User.DisplayName,
  })
}
