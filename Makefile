.PHONY:	build clean

all: build

build:
	./build.sh

clean:
	rm -rf ./tmp
	rm -rf ./dist