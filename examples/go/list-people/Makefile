build: list-people

list-people: *.go go.*
	GOOS=linux GOARCH=amd64 go build -v -o $@ -ldflags="-s -w"

clean:
	rm -rf list-people*

.PHONY: build clean
