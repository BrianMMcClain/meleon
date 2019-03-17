all: build

build: 
	go build -o meleon -v cmd/meleon/*

clean:
	rm meleon