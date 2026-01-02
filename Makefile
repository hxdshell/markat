default: buildnrun

build: 
	go build -o markat cmd/main.go

buildnrun: build
	clear && ./markat 