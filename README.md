[![Version](https://img.shields.io/badge/goversion-1.18.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>

# Usho
This project is a url shortener inspired by ursho.

# Build
Mac/Linux
```
go build -o usho ./cmd/http/main.go
```
Windows
```
go build -o usho.exe ./cmd/http/main.go
```
For particular platform:
```
env GOOS=linux GOARCH=amd64 go build -o usho ./cmd/http/main.go
```

# Run
For most cases just running the binary will do just fine.
## With InMem Store:
```
./usho \
-path=./path/to/file
```
## With Sql Database:
```
./usho \
-db
-dsn=dsnOfDatabase
```

# All Flags
```
./usho -help
Usage of usho:
  -db
        indicates wether to use mysql database, if ignored will use inmem store
  -dsn string
        data source name for mysql database
  -path string
        indicates where the file storage should be located (default "./store")
```