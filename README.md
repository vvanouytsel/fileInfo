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
$ docker run --rm fileinfo:latest -v /tmp
```

## Run
```
# Show the help menu
./fileInfo -h

# Show info of '/tmp/myfile'
./fileInfo /tmp/myfile

# Show info of '/tmp/myfile' with verbose mode enabled
./fileInfo -v
```