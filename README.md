# Generator and Parser for ldjson file

## For generating ldjson file

```
cd logsGenerator
go run main.go -n 100000 --from 2024-10-01 --to 2024-10-28
```

### or

```
cd logsGenerator
go build -o generator main.go
./generator -n 100000 --from 2024-09-10 --to 2024-10-10
```

## For parsing ldjson file

```
cd logsParser
go run main.go --filePath [pathForYourLDJSONFile]
```

### or

```
cd logsParser
go build -o parser main.go
./parser --filePath [pathForYourLDJSONFile]
```
