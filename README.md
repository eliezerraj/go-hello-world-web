# go-hello-world-web

POC just for tests using liveness/readiness probes and Dockerfile using CMD vs Shell (SIGTERM)

## Compile lambda

   Manually compile the function

      GOOD=linux GOARCH=amd64 go build -o ../build/main main.go

      zip -jrm ../build/main.zip ../build/main

https://localhost:PORT/pod-b/a
https://localhost:PORT/pod-b/info