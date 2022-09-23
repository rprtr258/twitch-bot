package internal

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
