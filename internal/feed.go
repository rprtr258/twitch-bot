package internal

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
