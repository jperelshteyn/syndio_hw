# syndio_hw
## About
Local HTTP server for employee data

## Requirements

Linux or MacOS machine

[Go](https://golang.org/doc/install)

## Installation
```bash
go get github.com/jperelshteyn/syndio_hw
```

## Usage
```bash
go run $(go env GOPATH)/src/github.com/jperelshteyn/syndio_hw/main.go [-port] [-seed_db] [-db_path]
```

### Optional Flags
- -port: local port to serve HTTP requests (defaults to `PORT` environment variable)
- -seed_db: cleans and fills DB with seed employee data (defaults to false)
- -db_path: path to SQLite DB file, which gets created if does not exist (defaults to current working directory)

## Access API
```bash
curl localhost:8888/employees
```