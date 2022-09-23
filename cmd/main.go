package main

import (
	"log"
	"net/http"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix"

	abobus "abobus/internal"
)

// app = Flask(__name__)
// with open("model.json", "r") as fd:
//     vocabulary = load(fd)["prior"].keys()

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

// cache = {}
// def levenshteinDistance(s1, s2):
//     import edlib
//     s1, s2 = s1.lower(), s2.lower()
//     if (s1, s2) in cache:
//         return cache[(s1, s2)]
//     cache[(s1, s2)] = edlib.align(s1, s2)["editDistance"]
//     return cache[(s1, s2)]

// def get_begins(message_word, vocabulary):
//     words = {}
//     for prior_ngram in vocabulary:
//         dist = levenshteinDistance(prior_ngram, message_word)
//         if dist < len(message_word) * 1.2:
//             words[prior_ngram] = dist
//     return list(sorted(words.items(), key=lambda kv:kv[1]))[:10]

// @app.route("/g")
// def genn():
//     from random import choices
//     import math
//     from quotes_generator.ngram import NGram
//     model = NGram(3)
//     model.load("model.json")
//     if request.args.get('m'):
//         message = request.args.get('m').replace(' ', '+').replace('_', ' ').split()
//         if len(message) > 1:
//             message = list(zip(message, message[1:]))
//         message = sorted(message, key=len)[-5:]
//         begins = sum([get_begins(w, vocabulary) for w in message], [])
//         print(begins, flush=True)
//         begin = choices([w for w, _ in begins], [math.pow(2, -p) for _, p in begins])
//         return model.generate(begin)
//     return model.generate()

// EMOTES = {
//     "AAUGH": (0, "рыбку"),
//     "VeryLark": (1, "стримера"),
//     "VeryPog": (3, "лысого из C++"),
//     "VerySus": (4, "AMOGUS"),
//     "VeryPag": (5, "анимешку"),
//     "SadgeCry": (6, "FeelsWeakMan"),
//     "LewdIceCream": (7, "бороду"),
// }

// def pack(username, emote_code, count, minute):
//     from base64 import b64encode
//     username = username.encode()
//     # 7 bits - count
//     # 3 bits - emote_code
//     # 6 bits - minute
//     packed = (count << 9) + (emote_code << 6) + minute
//     hhash = username + packed.to_bytes(2, 'big')
//     return b64encode(hhash).decode()

// def unpack(encoded):
//     from base64 import b64decode
//     decoded = b64decode(encoded)
//     username = decoded[:-2].decode()
//     unpacked = int.from_bytes(decoded[-2:], 'big')
//     minute = unpacked & 0x3F
//     unpacked >>= 6
//     emote_code = unpacked & 0x7
//     unpacked >>= 3
//     count = unpacked
//     emote = [(k, v) for k, v in EMOTES.items() if v[0] == emote_code][0]
//     return username, emote, count, minute

// def weight(count):
//     return 100 + 13 * count

// @app.route("/f")
// def feed():
//     from datetime import datetime
//     current_minute = datetime.now().minute
//     emote = request.args.get('e')
//     if emote == '0':
//         return "Ты можешь покормить: " + " ".join(EMOTES.keys())
//     if emote not in EMOTES:
//         return f"{emote} никого не кормит FeelsWeakMan"
//     encoded = request.args.get('c')
//     user = request.args.get('u')
//     if user == "Gekmi":
//         return "Гекми наказан за читерство с талончиками PauseFish"
//     if len(encoded) > 1:
//         try:
//             username, the_emote, count, minute = unpack(encoded)
//         except Exception as e:
//             return "Некорректный талончик"
//         if the_emote[0] != emote:
//             return f"Ты используешь талончик {the_emote[0]} для {emote} . Так нельзя. Получить талончик: !feed {emote}"
//         if username != user:
//             return f"Ты используешь талончик {username}, но тебя зовут {user}. Это называется воровство. Ты приговариваешься к расстрелу AAUGH"
//         if minute <= current_minute < minute + 5:
//             return f"Ты еще не можешь покормить {EMOTES[emote][1]}"
//         count += 1
//         w = weight(count)
//         response = pack(username, the_emote[1][0], count, minute)
//         return f"Ты покормил {the_emote[1][1]} {count} раз. Теперь он весит {w} грамм. Талончик на следующую кормежку: {response}"
//     response = pack(user, EMOTES[emote][0], 1, current_minute)
//     print(unpack(response))
//     w = weight(1)
//     return f"Ты покормил {EMOTES[emote][1]} в первый раз. Теперь он весит {w} грамм. Талончик на следующую кормежку: {response}"

// def read_db():
//     import json
//     with open("db.json", "r", encoding="utf-8") as fd:
//         return json.load(fd)

// def load_db(db):
//     import json
//     with open("db.json", "w", encoding="utf-8") as fd:
//         return json.dump(db, fd)

// last_balaboba = ""

// def balabob(text, skip=0):
//     global last_balaboba
//     if text == "":
//         text = last_balaboba
//     pitsots = 0
//     ln = len(text)
//     tries = 4
//     for _ in range(tries):
//         import requests
//         resp = requests.post('https://pelevin.gpt.dobro.ai/generate/', json={"prompt":text}, verify=False)
//         if resp.status_code == 500 or resp.status_code == 502:
//             pitsots += 1
//         else:
//             print(resp.content.decode('utf-8'))
//             resp = resp.json()
//             text = ' '.join(text.split() + max(resp['replies'], key=len).split())
//     if pitsots == tries:
//         return 'Порфирьевич в ахуе, попробуйте еще раз позже'
//     elif any(x in text.lower() for x in ['пидор', 'негр', 'нигер']):
//         return 'Порфирьевич сказал очень плохое слово, поэтому ловите рыбку AAUGH'
//     else:
//         return text[ln:]

//     from subprocess import check_output
//     TO_SKIP = len("please wait up to 15 seconds Без стиля".split())
//     output = check_output(["./balaboba"] + text.split()).decode("utf-8").split()
//     print(output)
//     if "на острые темы, например, про политику или религию" in ' '.join(output):
//         return "PauseFish"
//     last_balaboba = ' '.join(output[TO_SKIP + skip:])
//     return last_balaboba

// @app.route("/blab/<idd>")
// def long_blab(idd):
//     db = read_db()
//     db[idd] = balabob(db[idd])
//     load_db(db)
//     return f'''<p style="padding: 10% 15%; font-size: 1.8em;">{db[idd]}</p>'''

// @app.route("/b")
// def balaboba():
//     message = request.args.get("m")
//     response = balabob(message, skip=len(message.split()))
//     if len(response) < 300:
//         print("TOO SHORT: ", message, len(response))
//         return response
//     else:
//         db = read_db()
//         idd = max(map(int, db.keys())) + 1
//         db[idd] = message + ' ' + response
//         load_db(db)
//         return f"Читать продолжение в источнике: secure-waters-73337.herokuapp.com/blab/{idd}"

// @app.route("/s")
// def sber():
//     text = request.args.get("m")
//     import requests
//     for _ in range(6):
//         resp = requests.post('https://api-inference.huggingface.co/models/sberbank-ai/rugpt3small_based_on_gpt2', json={"inputs":text})
//         print(resp.content.decode('utf-8'))
//         if resp.status_code == 500 or 'error' in resp:
//             return 'Сбер сейчас в ахуе, попробуйте еще раз позже'
//         resp = resp.json()
//         text = ' '.join(resp[0]['generated_text'].split())
//     return text

func handleTwitchAuth(helixClient *helix.Client) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code, ok := r.URL.Query()["code"]
		if !ok {
			return
		}

		resp, err := helixClient.RequestUserAccessToken(code[0])
		if err != nil {
			log.Println(err.Error())
		}

		// Set the access token on the client
		// TODO: remember token, https://github.com/nicklaw5/helix/blob/main/docs/authentication_docs.md#refresh-user-access-token
		helixClient.SetUserAccessToken(resp.Data.AccessToken)
	})
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Println(err.Error())
	}
}

func run() error {
	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("TWITCH_REDIRECT_URI"),
	})
	if err != nil {
		return err
	}

	go handleTwitchAuth(helixClient)

	url := helixClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       []string{"user:read:email"},
		ForceVerify:  false,
	})
	log.Println(url)

	client := twitch.NewClient("trooba_bot", os.Getenv("TWITCH_OAUTH_TOKEN"))
	client.Join("vs_code")

	srvcs := abobus.Services{
		ChatClient:      client,
		TwitchApiClient: helixClient,
	}

	client.OnPrivateMessage(srvcs.OnPrivateMessage)

	if err := client.Connect(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
