import {env} from "process";
import tmi from "tmi.js";
import * as internal from "./internal/internal.ts";
import {LoadFromJSONFile} from "./internal/permissions.ts";

const style = `body {
  background-color: #09090b;
  color: #fff;
}
a {
  color: #ececec;
  text-decoration: underline;
  font-size: 1.8em;
}
a:visited {
  color: #333;
}
table {
  margin: 0 auto;
  text-align: center;
}`;

// TODO: show partially content
type pageEntry = {
  ID:  string,
  URL: string,
};
const idsPage = (x: pageEntry[][]) => `
<head>
  <style>` + style + // TODO: external css file
    `</style>
</head>
<body>
  <table>
  ${x.map(x => `
    <tr>
    ${x.map(({ID, URL}) => `
        <td><a href=${URL}>${ID}</a></td>
    `)}
    </tr>
  `)}
  </table>
</body>
`;
const blabPage = (x: string) => `
<head>
  <style>` + style + // TODO: external css file
    `</style>
</head>
<body>
  <p style="padding: 10% 15%; font-size: 1.8em;">${x}</p>
</body>
`;

// const helixClient = helix.NewClient({
//   ClientID:     env["TWITCH_CLIENT_ID"],
//   ClientSecret: env["TWITCH_CLIENT_SECRET"],
//   RedirectURI:  env["TWITCH_REDIRECT_URI"],
// });

const channels: string[] = app.DB().Select("channel").From("joined_channels").Build().Rows();

// TODO: provide proxy in env
const proxyURL = new URL(env["BALABOBA_PROXY"]!);

// const d = net.Dialer{
//   Timeout: time.Minute,
// };
// const balabobaClient = balaboba.New(balaboba.ClientConfig{
//   Lang: balaboba.Rus,
//   HTTP: &http.Client{
//     Timeout: d.Timeout,
//     Transport: &http.Transport{
//       DialTLSContext:      d.DialContext,
//       TLSHandshakeTimeout: d.Timeout,
//       Proxy:               http.ProxyURL(proxyURL),
//     },
//   },
// })

const permissions = LoadFromJSONFile("permissions.json");

const services: services.Services = {
  // ChatClient:      client,
  // TwitchApiClient: helixClient,
  // Backend:         app,
  // Balaboba:        balabobaClient,
  Permissions:     permissions,
};

const client = new tmi.Client({
	options: {debug: true},
	identity: {
		username: env["TWITCH_USERNAME"],
		password: env["TWITCH_OAUTH_TOKEN"],
	},
	channels: channels,
});
client.connect().catch(console.error);
client.on("whisper", (from, userstate, message, self) => {
	if (self)
    return;
  internal.OnWhisperMessage(services, {user: from, message, channel: from});
});
client.on("message", (channel, tags, message, self) => {
	if (self)
    return;
	// if (message.toLowerCase() === '!hello') {
	// 	client.say(channel, `@${tags.username}, heya!`);
	// }
  internal.OnPrivateMessage(services, {user: tags.username!, message, channel});
});

function chunk<T>(array: readonly T[], size: number): T[][] {
  return Array.from(
    {length: Math.ceil(array.length / size)},
    (_, i) => array.slice(i * size, (i + 1) * size),
  );
}

const server = Bun.serve({
  routes: {
    "/blab": {
      "GET": (req) => {
        const db = app.DB();

        const ids: string[] = db.Select("id").From("blab").All().map((row: {ID: string}) => row.ID);

        const entries = ids.map((ID): pageEntry => ({
          ID:  ID,
          URL: `/blab/${ID}`,
        }))
        const entriesChunks = chunk(entries, 4);
        return new Response(idsPage(entriesChunks), {status: 200});
      },
    },
    "/blab/:id": {
      "GET": (req) => {
        const blabID = req.params.id;
        const db = app.DB();
        const text: string = db.Select("text").From("blab").Where(`id=${blabID}`).Row();
        return new Response(blabPage(text), {status: 200});
      },
    },
  },
});

console.log(`Server running at ${server.url}`);
