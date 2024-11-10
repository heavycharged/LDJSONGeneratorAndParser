# Generator and Parser for ldjson file

## Install hyperfine
Install [hyperfine](https://github.com/sharkdp/hyperfine) to compare parsing speed.


## For generating ldjson file

```
go build -o generate ./internal/generator/main.go
./generate -n 100000 --from 2024-09-10 --to 2024-10-10
```

## For parsing ldjson file

```
go build -o parse ./internal/parser/main.go
cat logs.ldjson | ./parse
```
