import requests
import json
import re

def normalize_line(line):
    res = line.strip()
    for a, b in zip("acekopxy", "асекорху"):
        res = res.replace(a, b)
    for a, b in zip("ABCEHKMOPTX", "АВСЕНКМОРТХ"):
        res = res.replace(a, b)
    res = re.sub(r"[^абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ]", " ", res)
    for a, b in zip("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ", "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"):
        res = res.replace(a, b)
    return " ".join(res.split())

with open("VK_TOKEN", "r") as f:
    TOKEN = f.readline()

FORBIDDEN = ["http", "like", "лайк", "паблик", "приглашайте", "послушайте", "подпи"]

GROUP_IDS = [27456813, 35927256, 167038479]
for group_id in GROUP_IDS:
    for offset in range(0, 1000, 100):
        result = requests.get(f"https://api.vk.com/method/wall.get?owner_id=-{group_id}&offset={offset}&count=100&filter=owner&access_token={TOKEN}&v=5.102").content
        data = json.loads(result)["response"]["items"]
        for x in data:
            text = normalize_line(x["text"])
            if len(text) > 1 and not any(map(lambda x:x in text, FORBIDDEN)):
                print(text)
