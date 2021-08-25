# Quotes generator
Patsanskie quotes generator with N-Gram model. Dataset was collected from the following groups in vk:

- [Пацанские цитаты©](https://vk.com/public27456813)
- [Пацанские цитаты №1](https://vk.com/krutiecitati)
- [пацанские цитаты и попуги](https://vk.com/tupopopugai) (more from this)

Vanilla quotes were collected from
- [Самые ванильные цитаты (500 цитат)](https://citatnica.ru/citaty/samye-vanilnye-tsitaty-500-tsitat)

## Usage

### Get dataset(if you need to):

Write your vk api token to VK_TOKEN file. Then execute

```bash
python3 get_dataset.py > dataset
```

### Train model(if you need to):

```bash
python3 main.py train
```

### Generate quotes:

```bash
python3 main.py generate
```
