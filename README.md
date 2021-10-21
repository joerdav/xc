# xc - eXeCute project tasks from a readme file

## Installation

```
go install github.com/joe-davidson1802/xc/cmd/xc@latest
```

## Tasks

__get__: get dependencies of the project
```
go get ./...
```

### Tests

__test__: test project
```
go test ./...
```

log-message-1: log a message
```
echo "running tests 1!"
```
log-message-2: log a message
```
echo "running tests 2!"
```
log-message-3: log a message
```
echo "running tests 3!"
```

_testshort_: run short tests
!log-message-1, log-message-2
!log-message-3, log-message-1
```
go test ./... -short
```

circular-test-1: test circular deps
!circular-test-2
```
echo 1
```
circular-test-2: test circular deps
!circular-test-1
```
echo 1
```

## More docs

Example docs
