# go-hello-world-web

Workload just for tests:

liveness/readiness probes

Dockerfile using CMD vs Shell (SIGTERM)

Autentication/Authorization (JWT Oauth 2.0 + Scopes)

## Compile

   Manually compile the function

      GOOD=linux GOARCH=amd64 go build -o ../build/main main.go

      zip -jrm ../build/main.zip ../build/main

## Endpoints

GET

      http://localhost:PORT/pod-b/a
      http://localhost:PORT/pod-b/b
      http://localhost:PORT/pod-b/info
      http://localhost:PORT/pod-b/version
      http://localhost:PORT/pod-b/header

POST

      http://localhost:PORT/pod-b/sum/2

GET (authenticatio/authorization - JWT)
      
      http://localhost:PORT/pod-b/methodToken