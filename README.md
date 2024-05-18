# gobpe
BPE Tokenizer implementation in Go

> Please note that this code has been written solely for learning purposes. 

## Usage

Follow the below example or just run `make demo` to perform all of these steps.

```shell
# build all commands
make build

# download works of Adam Mickiewicz
cat data/url/mickiewicz.txt | ./bin/load -datadir data/txt

# remove copy notice for processing
./bin/preprocess -datadir data/txt

# train tokenizer
cat data/txt/*.txt | ./bin/train -params params.json

# encode text to tokens
echo "SOME TEXT" | ./bin/encode -params params.json

# decode tokens to text
echo "[65,100]" | ./bin/decode -params params.json
```

## Important note

We are using content from https://wolnelektury.pl

## TODO
* Add support for special tokens