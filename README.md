# Generator and Parser for ldjson file

## For generating ldjson file

```
go build -o generate ./internal/generator/main.go
./generate -n 100000 --from 2024-09-10 --to 2024-10-10
```

## For parsing ldjson file

```
cd logsParser
go build -o parse ./internal/parser/main.go
cat logs.ldjson | ./parse
```
