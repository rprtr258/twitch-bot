#!/usr/bin/env python3
from argparse import ArgumentParser

from ngram import NGram


WORDS_BAG_SIZE = 3


def train(dataset_fd, model_fd):
    model = NGram(WORDS_BAG_SIZE)
    model.train(dataset_fd)
    model.save(model_fd)

def generate(model_fd):
    import sys; sys.stdout.reconfigure(encoding="utf-8")

    model = NGram(WORDS_BAG_SIZE)
    model.load(model_fd)
    print(model.generate())

if __name__ == "__main__":
    parser = ArgumentParser(
        prog="NGram",
        description="N-Gram model for text generation."
    )
    subparsers = parser.add_subparsers(dest="cmd")
    parser_train = subparsers.add_parser(
        "train",
        help="train model using dataset"
    )
    parser_train.add_argument(
        "-d",
        "--dataset",
        type=str,
        help="dataset to use for training",
        required=True,
    )
    parser_train.add_argument(
        "-m",
        "--model",
        type=str,
        help="file to save model to",
        required=True,
    )
    parser_generate = subparsers.add_parser(
        "generate",
        help="generate text using trained model"
    )
    parser_generate.add_argument(
        "-m",
        "--model",
        type=str,
        help="json model to use for generation",
        required=True,
    )
    args = parser.parse_args()
    if args.cmd == "train":
        train(args.dataset, args.model)
    elif args.cmd == "generate":
        generate(args.model)
