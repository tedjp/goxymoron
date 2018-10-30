all: goxymoron

goxymoron: goxymoron.go
	go build

run: goxymoron
	./$<
