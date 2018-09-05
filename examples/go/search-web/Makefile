build: search-web

search-web: *.go go.*
	GOOS=linux GOARCH=amd64 go build -v -o $@ -ldflags="-s -w"

clean:
	rm -rf search-web*

.PHONY: build clean
