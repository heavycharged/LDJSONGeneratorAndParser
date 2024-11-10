build: generate
	go build -o bin/parser-modified ./internal/parser/main.go

generate:
	go generate ./...

compare: build 
	hyperfine --warmup 1 --runs 10 --export-markdown results.md --export-json results.json \
	"bin/parser-control < logs.ldjson" \
	"bin/parser-modified < logs.ldjson"