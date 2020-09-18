# fileInfo
## Build
### Native
```
$ go build .
```

### Docker
```
$ make run
$ docker run --rm fileinfo:latest [ARGUMENTS]

# Example
$ docker run --rm fileinfo:latest -v /dev
```

## Run
```
# Show the help menu
./fileInfo -h

# Show info of '/dev'
./fileInfo /dev

# Show info of '/dev' with verbose mode enabled
./fileInfo -v /dev
```