build: bind

bind: *.go go.*
	GOOS=linux GOARCH=amd64 go build -v -o $@ -ldflags="-s -w"

clean:
	rm -rf bind*

.PHONY: build clean
