# Quotes generator

## Usage

### Train model (if you need to) using dataset `dataset.txt` and save into `model.json`:

```bash
./main.py train -d dataset.txt -m model.json
```

### Generate quotes using model `model.json`:

```bash
./main.py generate -m model.json
```

## Datasets

Some default datasets and trained models provided.

### serious.txt, serious.json

Collected from VK groups:

- [Пацанские цитаты©](https://vk.com/public27456813)
- [Пацанские цитаты №1](https://vk.com/krutiecitati)
- [пацанские цитаты и попуги](https://vk.com/tupopopugai) (more from this)

### vanilla.txt, vanilla.json

Collected from:

- [Самые ванильные цитаты (500 цитат)](https://citatnica.ru/citaty/samye-vanilnye-tsitaty-500-tsitat)

