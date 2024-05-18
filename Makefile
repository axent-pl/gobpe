build:
	go build -o bin/load ./cmd/load
	go build -o bin/preprocess ./cmd/preprocess
	go build -o bin/encode ./cmd/encode
	go build -o bin/decode ./cmd/decode
	go build -o bin/train ./cmd/train

# Build commands and execute:
# * bin/load - Load data
# * bin/preprocess - Preprocess data (remove copy notice)
# * bin/train - Train tokenizer
# * bin/encode - Enocde sample string
# * bin/decode - Enocde sample string and decode back to original
demo: build
	cat data/url/mickiewicz.txt | ./bin/load -datadir data/txt
	./bin/preprocess -datadir data/txt
	cat data/txt/*.txt | ./bin/train -params params.json
	echo "Adam Mickiewicz sławnym poetą był" | ./bin/encode -params params.json
	echo "Adam Mickiewicz sławnym poetą był" | ./bin/encode -params params.json | ./bin/decode -params params.json
