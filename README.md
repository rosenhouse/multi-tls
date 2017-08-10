## tests
```
go get github.com/rosenhouse/multi-tls
cd $GOPATH/src/github.com/rosenhouse/multi-tls
ginkgo integration-test
```

## wat?
```
client --> router ==> backends
        ^          ^
       http       https
```

where:
- backend-choice is governed by request's `Host` header
- backend presents cert with unique common-name, router validates

## deets

### phase 1 (current)
- launch N backends
- backend i listens on port 10000+i
- backend responds to a request by saying 'hello i'm backend listening on port 10000+i'

- launch router, listens on port 2000

- client makes request to router, http header 'Host: backend-10003'
- expects to receive response from backend listening on port 10003


### phase 2
- router should do HTTPS to backend
- backend should present certificate saying common name is "backend-10003"
- router should require that backend present this certificate


### phase 3
- optimize!
